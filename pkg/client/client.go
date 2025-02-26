package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	urlBase           = "https://platform.ringcentral.com/restapi"
	oauthURL          = "/oauth/token"
	getExtensions     = "/v1.0/account/~/extension"
	getAvailableRoles = "/v1.0/account/~/user-role"
	userRoles         = "/v1.0/account/~/extension/%s/assigned-role"
)

type RingCentralClient struct {
	client      *uhttp.BaseHttpClient
	accessToken string
	Config      ClientConfig
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	JWT          string
}

type Option func(c *RingCentralClient)

func WithAccessToken(accessToken string) Option {
	return func(c *RingCentralClient) {
		c.accessToken = accessToken
	}
}

func WithClientID(clientID string) Option {
	return func(c *RingCentralClient) {
		c.Config.ClientID = clientID
	}
}

func WithClientSecret(clientSecret string) Option {
	return func(c *RingCentralClient) {
		c.Config.ClientSecret = clientSecret
	}
}

func WithJWT(jwt string) Option {
	return func(c *RingCentralClient) {
		c.Config.JWT = jwt
	}
}

func (c *RingCentralClient) GetToken() string {
	return c.accessToken
}

func New(ctx context.Context, opts ...Option) (*RingCentralClient, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	cli, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, err
	}

	rcClient := RingCentralClient{
		client: cli,
	}

	for _, o := range opts {
		o(&rcClient)
	}

	if rcClient.Config.ClientID != "" && rcClient.Config.ClientSecret != "" && rcClient.Config.JWT != "" {
		newAccessToken, err := rcClient.requestAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		rcClient.accessToken = newAccessToken
	}

	return &rcClient, nil
}

func (c *RingCentralClient) requestAccessToken(ctx context.Context) (string, error) {
	requestURL, err := url.JoinPath(urlBase, oauthURL)
	if err != nil {
		return "", err
	}

	clientData := c.Config.ClientID + ":" + c.Config.ClientSecret
	encodedClientData := base64.StdEncoding.EncodeToString([]byte(clientData))

	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Add("assertion", c.Config.JWT)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Basic "+encodedClientData)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatalf("Error haciendo la solicitud: %v", err)
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (c *RingCentralClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	body interface{},
	reqOpts ...ReqOpt,
) (http.Header, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, err
	}

	for _, o := range reqOpts {
		o(urlAddress)
	}

	req, err := c.client.NewRequest(
		ctx,
		method,
		urlAddress,
		uhttp.WithAcceptJSONHeader(),
		uhttp.WithContentTypeJSONHeader(),
		uhttp.WithHeader("Authorization", "Bearer "+c.GetToken()),
		uhttp.WithJSONBody(body),
	)
	if err != nil {
		return nil, err
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if res != nil {
		bodyContent, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bodyContent, &res)
		if err != nil {
			return nil, err
		}
	}

	return resp.Header, nil
}

// ListAllUsers returns an array of users of the platform belonging to the company.
// Users withing the platform are named as 'Extension'.
func (c *RingCentralClient) ListAllUsers(ctx context.Context, pageOps PageOptions) ([]Extension, string, error) {
	var response ExtensionResponse

	queryUrl, err := url.JoinPath(urlBase, getExtensions)
	if err != nil {
		return nil, "", err
	}

	nextPage, err := c.getExtensionsListFromAPI(ctx, queryUrl, &response, WithPage(pageOps.Page), WithPageLimit(pageOps.PerPage))
	if err != nil {
		return nil, "", err
	}

	return response.Records, nextPage, nil
}

func (c *RingCentralClient) ListAllAvailableRoles(ctx context.Context, pageOps PageOptions) ([]Role, string, error) {
	var response RoleResponse

	queryUrl, err := url.JoinPath(urlBase, getAvailableRoles)
	if err != nil {
		return nil, "", err
	}

	nextPage, err := c.getRolesListFromAPI(ctx, queryUrl, &response, WithPage(pageOps.Page), WithPageLimit(pageOps.PerPage))
	if err != nil {
		return nil, "", err
	}

	return response.Records, nextPage, nil
}

func (c *RingCentralClient) GetUserAssignedRoles(ctx context.Context, userResource *v2.Resource) ([]UserRole, error) {
	var res UserRoleResponse
	queryUrl, err := url.JoinPath(urlBase, fmt.Sprintf(userRoles, userResource.Id.Resource))
	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(ctx, http.MethodGet, queryUrl, &res, nil)
	if err != nil {
		return nil, err
	}

	return res.Records, nil
}

func (c *RingCentralClient) getExtensionsListFromAPI(
	ctx context.Context,
	urlAddress string,
	res *ExtensionResponse,
	reqOpt ...ReqOpt,
) (string, error) {
	var pageToken string

	_, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, nil, reqOpt...)
	if err != nil {
		return "", err
	}

	if res.Paging.Page < res.Paging.TotalPages {
		pageToken = strconv.Itoa(res.Paging.Page + 1)
	}

	return pageToken, nil
}

func (c *RingCentralClient) getRolesListFromAPI(ctx context.Context, urlAddress string, res *RoleResponse, reqOpt ...ReqOpt) (string, error) {
	var pageToken string

	_, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, nil, reqOpt...)
	if err != nil {
		return "", err
	}

	if res.Paging.Page < res.Paging.TotalPages {
		pageToken = strconv.Itoa(res.Paging.Page + 1)
	}

	return pageToken, nil
}

// IdKeyValue is an auxiliary structure to build the body for the operation of update the roles list of a user.
type IdKeyValue struct {
	Id string `json:"id"`
}

// UpdateUserRoles receives the user resource (the principal of the Grant operation) and request the curren assigned roles for it.
// This function can be called on "revoking mode" or "granting mode". isRevoking sets the behavior.
// While granting: if the role that will be assigned isn't already part of the user roles, it sends the whole role list to the platform.
// While revoking: the list of assigned roles is sent to the platform by previously deleting the desired role.
func (c *RingCentralClient) UpdateUserRoles(ctx context.Context, userResource *v2.Resource, roleID string, isRevoking bool) error {
	// This variable is initialized like this and not with the "var roleIDs []IdKeyValue" semantic since it produces a bug when the array receives no elements.
	roleIDs := []IdKeyValue{}

	// Request the list of the assigned roles of the user to be able to add the new one to that list.
	assignedRoles, err := c.GetUserAssignedRoles(ctx, userResource)
	if err != nil {
		return err
	}

	for _, assignedRole := range assignedRoles {
		arID := assignedRole.Id

		if arID == roleID {
			if isRevoking {
				// While revoking: If the current role ID equals the role that must be removed, it avoids adding it to the roles list.
				continue
			} else {
				// While granting: If the current role ID equals the role should be added, an error is thrown, since the user already has that role.
				return fmt.Errorf("the role with ID: '%s' is already assigned to the user with ID: '%s'", roleID, userResource.Id.Resource)
			}
		}

		roleIDs = append(roleIDs, IdKeyValue{Id: arID})
	}

	if !isRevoking {
		roleIDs = append(roleIDs, IdKeyValue{Id: roleID})
	}

	body := map[string]interface{}{
		"records": roleIDs,
	}
	requestURL, err := url.JoinPath(urlBase, fmt.Sprintf(userRoles, userResource.Id.Resource))
	if err != nil {
		return err
	}

	_, err = c.doRequest(ctx, http.MethodPut, requestURL, nil, body)
	if err != nil {
		return err
	}

	return nil
}
