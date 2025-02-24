package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"io"
	"net/http"
	"net/url"
)

const (
	urlBase       = "https://platform.ringcentral.com/restapi/v1.0"
	getExtensions = "/account/~/extension"
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

func (c *RingCentralClient) getExtensionsListFromAPI(
	ctx context.Context,
	urlAddress string,
	res *ExtensionResponse,
	reqOpt ...ReqOpt,
) (string, annotations.Annotations, error) {
	_, annotation, err := c.doRequest(ctx, http.MethodGet, urlAddress, &res, reqOpt...)

	if err != nil {
		return "", nil, err
	}

	var pageToken string
	nav := res.Navigation
	logger := ctxzap.Extract(ctx)
	logger.Info(fmt.Sprintf("Paging Token: %v", nav))
	logger.Info(fmt.Sprintf("URI: %v", res.Uri))
	logger.Info(fmt.Sprintf("FIRST PAGE URI: %v", nav.FirstPage.Uri))
	logger.Info(fmt.Sprintf("LAST PAGE URI: %v", nav.LastPage.Uri))

	if res.Uri == nav.LastPage.Uri {
		pageToken = ""
	}

	return pageToken, annotation, nil
}

func (c *RingCentralClient) doRequest(
	ctx context.Context,
	method string,
	endpointUrl string,
	res interface{},
	reqOpts ...ReqOpt,
) (http.Header, annotations.Annotations, error) {
	var (
		resp *http.Response
		err  error
	)

	urlAddress, err := url.Parse(endpointUrl)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	if res != nil {
		bodyContent, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, nil, err
		}
		err = json.Unmarshal(bodyContent, &res)
		if err != nil {
			return nil, nil, err
		}
	}
	annotation := annotations.Annotations{}
	return resp.Header, annotation, nil
}

// ListAllUsers returns an array of users of the platform belonging to the company.
// Users withing the platform are named as 'Extension'.
func (c *RingCentralClient) ListAllUsers(ctx context.Context, pageOps PageOptions) ([]Extension, string, error) {
	var response ExtensionResponse
	queryUrl, err := url.JoinPath(urlBase, getExtensions)

	if err != nil {
		return nil, "", err
	}

	nextPage, _, err := c.getExtensionsListFromAPI(ctx, queryUrl, &response, WithPage(pageOps.Page), WithPageLimit(pageOps.PerPage))
	if err != nil {
		return nil, "", err
	}

	return response.Records, nextPage, nil
}
