package client

import (
	"context"
	"encoding/json"
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

func (c *RingCentralClient) getExtensionsListFromAPI(
	ctx context.Context,
	urlAddress string,
	res *ExtensionResponse,
	reqOpt ...ReqOpt,
) (string, error) {
	var pageToken string

	_, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOpt...)
	if err != nil {
		return "", err
	}

	nav := res.Navigation
	if res.Uri == nav.LastPage.Uri {
		pageToken = ""
	}

	return pageToken, nil
}

func (c *RingCentralClient) getRolesListFromAPI(ctx context.Context, urlAddress string, res *RoleResponse, reqOpt ...ReqOpt) (string, error) {
	var pageToken string

	_, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOpt...)
	if err != nil {
		return "", err
	}

	nav := res.Navigation
	if res.Uri == nav.LastPage.Uri {
		pageToken = ""
	}

	return pageToken, nil
}
