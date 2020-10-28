package httpclient

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/scocogon/opensearch/osearch/response"
)

const (
	UrIPrefix = "/v3/openapi/apps/"
	OsPrefix  = "OPENSEARCH"
)

type Client struct {
	cfg *Config

	cli *http.Client
	m   sync.Mutex
}

// New new http client
func New(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Request(appName string, params, headers map[string]string) (*response.Response, error) {
	byts, err := c.SendRequest(appName, params, headers)
	if err != nil {
		return nil, err
	}

	r := &response.Response{}
	err = json.Unmarshal(byts, r)
	return r, err
}

func (c *Client) TryRun(appName string, params, headers map[string]string) {
	url := c.cfg.host + c.buildQuery(appName, params, headers)

	fmt.Printf("curl ")
	for k, v := range headers {
		fmt.Printf("-H \"%s: %s\" ", k, v)
	}

	fmt.Println(url)
	return
}

func (c *Client) SendRequest(appName string, params, headers map[string]string) ([]byte, error) {
	url := c.cfg.host + c.buildQuery(appName, params, headers)

	if c.cli == nil {
		c.m.Lock()

		if c.cli == nil {
			c.cli = &http.Client{
				Timeout: 3 * time.Second,
			}
		}

		c.m.Unlock()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// 标记这个操作失败
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, e := c.cli.Do(req)
	if e != nil {
		fmt.Fprintf(os.Stderr, "[HTTP GET] addr: %s, err: %s\n", url, e.Error())
		return nil, e
	}

	defer resp.Body.Close()

	// fmt.Fprintf(os.Stderr, "[HTTP GET] addr: %s, resp: %+v\n", addr, resp)
	buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))

	_, e = buf.ReadFrom(resp.Body)

	return buf.Bytes(), e
}

func (c *Client) SetHttpClient(cli *http.Client) {
	c.m.Lock()
	c.cli = cli
	c.m.Unlock()
}

func (c *Client) buildQuery(appName string, params, headers map[string]string) string {
	uri := UrIPrefix
	if len(appName) != 0 {
		uri += appName
	}
	uri += "/search"

	c.buildRequestHeader(uri, params, headers)

	return c.canonicalizedResource(uri, params)
}

func (c *Client) buildRequestHeader(uri string, params, headers map[string]string) {
	for k, v := range headers {
		headers[k] = v
	}

	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}

	if _, ok := headers["Date"]; !ok {
		headers["Date"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	}

	if _, ok := headers["X-Opensearch-Nonce"]; !ok {
		headers["X-Opensearch-Nonce"] = generateNonce()
	}

	if _, ok := headers["Authorization"]; !ok {
		headers["Authorization"] = c.buildAuthorization(uri, params, headers)
	}
}

func (c *Client) buildAuthorization(uri string, params, headers map[string]string) string {
	resource := c.canonicalizedResource(uri, params)
	signature, _ := c.Signature(resource, headers)

	return fmt.Sprintf("%s %s:%s", OsPrefix, c.cfg.accessKeyID, signature)
}

func (c *Client) Signature(resource string, reqHeaders map[string]string) (signature string, err error) {
	contentMD5 := ""
	contentType := ""
	date := time.Now().UTC().Format(http.TimeFormat)

	if v, exist := reqHeaders["Content-MD5"]; exist {
		contentMD5 = v
	}

	if v, exist := reqHeaders["Content-Type"]; exist {
		contentType = v
	}

	if v, exist := reqHeaders["Date"]; exist {
		date = v
	}

	header := c.canonicalizedHeaders(reqHeaders)

	stringToSign := `GET` + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		header + "\n" +
		resource

	sha1Hash := hmac.New(sha1.New, []byte(c.cfg.accessKeySecret))
	if _, e := sha1Hash.Write([]byte(stringToSign)); e != nil {
		return "", e
	}

	signature = base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))
	return
}

func (c *Client) canonicalizedHeaders(reqHeaders map[string]string) string {
	var headers sort.StringSlice
	for k, v := range reqHeaders {
		if strings.HasPrefix(k, "X-Opensearch-") && len(v) > 0 {
			headers = append(headers, strings.ToLower(k)+":"+strings.TrimSpace(v))
		}
	}

	if headers.Len() == 0 {
		return ""
	}

	sort.Sort(headers)
	return strings.Join(headers, "\n")
}

func (c *Client) canonicalizedResource(uri string, params map[string]string) string {
	canonicalized := strings.Replace(quote(uri), "%2F", "/", -1)

	var param sort.StringSlice
	for k, v := range params {
		if len(v) > 0 {
			param = append(param, quote(k)+"="+quote(v))
		}
	}

	if param.Len() == 0 {
		return ""
	}
	sort.Sort(param)

	return canonicalized + "?" + strings.Join(param, "&")
}

func generateNonce() string {
	v := 1000 + rand.Int63n(8999) + time.Now().UnixNano()/1000
	return strconv.FormatInt(v, 10)
}

func quote(s string) string {
	s = url.QueryEscape(s)
	s = strings.Replace(s, "+", "%20", -1)

	return s
}
