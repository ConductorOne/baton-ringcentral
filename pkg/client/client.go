package client

import (
	"context"
	"encoding/json"
	"fmt"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"io"
	"net/http"
	"net/url"
)

const (
	urlBase           = "https://platform.ringcentral.com/restapi/v1.0"
	getExtensions     = "/account/~/extension"
	getAvailableRoles = "/account/~/user-role"
	userRoles         = "/account/~/extension/%s/assigned-role"
)

type RingCentralClient struct {
	client      *uhttp.BaseHttpClient
	accessToken string
}

type Option func(c *RingCentralClient)

func (c *RingCentralClient) GetToken() string {
	return c.accessToken
}

func WithAccessToken(accessToken string) Option {
	return func(c *RingCentralClient) {
		c.accessToken = accessToken
	}
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

	return &rcClient, nil
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

	// TODO implement proper pagination for this API
	nav := res.Navigation
	if res.Uri == nav.LastPage.Uri {
		pageToken = ""
	}

	return pageToken, nil
}

func (c *RingCentralClient) getRolesListFromAPI(ctx context.Context, urlAddress string, res *RoleResponse, reqOpt ...ReqOpt) (string, error) {
	var pageToken string

	_, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, nil, reqOpt...)
	if err != nil {
		return "", err
	}

	// TODO implement proper pagination for this API
	nav := res.Navigation
	if res.Uri == nav.LastPage.Uri {
		pageToken = ""
	}

	return pageToken, nil
}

// IdKeyValue is an auxiliary structure to build the body for the operation of update the roles list of a user.
type IdKeyValue struct {
	Id string `json:"id"`
}

// UpdateUserRole receives the user resource (the principal of the Grant operation) and request the curren assigned roles for it.
// Then, if the role that will be assigned isn't already part of the user roles, it sends the whole role list to the platform.
func (c *RingCentralClient) UpdateUserRole(ctx context.Context, userResource *v2.Resource, roleID string) error {
	var roleIDs []IdKeyValue

	// Request the list of the assigned roles of the user to be able to add the new one to that list.
	assignedRoles, err := c.GetUserAssignedRoles(ctx, userResource)
	if err != nil {
		return err
	}

	for _, assignedRole := range assignedRoles {
		arID := assignedRole.Id

		if arID == roleID {
			return fmt.Errorf("the role with ID: '%s' is already assigned to the user with ID: '%s'", roleID, userResource.Id.Resource)
		}
		roleIDs = append(roleIDs, IdKeyValue{Id: arID})
	}
	roleIDs = append(roleIDs, IdKeyValue{Id: roleID})

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
