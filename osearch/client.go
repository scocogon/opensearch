package osearch

import (
	"strings"

	"github.com/scocogon/opensearch/oerrors"
	"github.com/scocogon/opensearch/osearch/httpclient"
	"github.com/scocogon/opensearch/osearch/response"
)

type clientHandle interface {
	TryRun(appName string, params, headers map[string]string)

	Request(appName string, params, headers map[string]string) (*response.Response, error)
	SendRequest(appName string, params, headers map[string]string) ([]byte, error)
}

// Client client
type Client struct {
	cli clientHandle

	query   []Source
	ffields []string

	params map[string]string
}

// New new client
func New(accessKeyID, accessKeySecret, host string) *Client {
	var cli clientHandle

	switch {
	case strings.HasPrefix(host, "http://"):
		cli = httpclient.New(httpclient.NewConfig(accessKeyID, accessKeySecret, host))

	default:
		panic("不支持 host")
	}

	return &Client{
		cli: cli,
	}
}

func (c *Client) TryRun(appName string, headers map[string]string) {
	if len(c.query) == 0 {
		return
	}

	var param = map[string]string{}
	var s []string

	for _, v := range c.query {
		s = append(s, v.Source())
	}

	param["query"] = strings.Join(s, "&&")

	// fetch_fields 可选
	if len(c.ffields) > 0 {
		param["fetch_fields"] = strings.Join(c.ffields, ";")
	}

	// 其他参数
	if len(c.params) > 0 {
		for k, v := range c.params {
			param[k] = v
		}
	}

	if headers == nil {
		headers = map[string]string{}
	}

	c.cli.TryRun(appName, param, headers)
}

// Send send request
func (c *Client) Send(appName string, headers map[string]string) (*response.Response, error) {
	var param = map[string]string{}
	var s []string

	for _, v := range c.query {
		if v != nil {
			s = append(s, v.Source())
		}
	}

	if len(s) == 0 {
		return nil, oerrors.ErrQueryNotFound
	}

	param["query"] = strings.Join(s, "&&")

	// fetch_fields 可选
	if len(c.ffields) > 0 {
		param["fetch_fields"] = strings.Join(c.ffields, ";")
	}

	// 其他参数
	if len(c.params) > 0 {
		for k, v := range c.params {
			param[k] = v
		}
	}

	if headers == nil {
		headers = map[string]string{}
	}

	return c.cli.Request(appName, param, headers)
}

// AddQuerys add query
/**
 * add config, query, sort, filter, ...
 */
func (c *Client) AddQuerys(src ...Source) *Client {
	c.query = append(c.query, src...)
	return c
}

// AddFetchFields add fetch_fields
func (c *Client) AddFetchFields(fields ...string) *Client {
	c.ffields = append(c.ffields, fields...)
	return c
}

// AddOtherParams add other params
/**
 * url: https://help.aliyun.com/document_detail/57155.html?spm=a2c4g.11174283.6.638.b1db5a19qzARRK
 * 除 query, fetchfiles 字段外的参数
 */
func (c *Client) AddOtherParams(params map[string]string) *Client {
	if c.params == nil {
		c.params = params
	} else {
		for key, value := range params {
			c.params[key] = value
		}
	}

	return c
}

// AddOtherParams2 add other params
/**
 * url: https://help.aliyun.com/document_detail/57155.html?spm=a2c4g.11174283.6.638.b1db5a19qzARRK
 * 除 query, fetchfiles 字段外的参数
 */
func (c *Client) AddOtherParams2(src ...ParamHandle) *Client {
	if c.params == nil {
		c.params = map[string]string{}
	}

	for _, v := range src {
		c.params[v.Key()] = v.Source()
	}

	return c
}
