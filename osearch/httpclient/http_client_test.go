package httpclient

import (
	"testing"
)

const (
	accessKeyID     = ""
	accessKeySecret = ""
	host            = "http://intranet.opensearch-cn-qingdao.aliyuncs.com"
	appName         = "qa_qtt_content"
)

// go test -timeout 30s -tags qa github.com/scocogon/opensearch/osearch/httpclient -run "^(TestHttpClient_Authorization)$" -v -count=1
func TestHttpClient_Authorization(t *testing.T) {

	cfg := NewConfig(accessKeyID, accessKeySecret, host)
	c := New(cfg)

	params := map[string]string{
		"query":        `config=format:json,start:0,hit:20,rerank_size:200&&query=create_time:[1579449600000,1582127999000]&&filter=content_type=3 AND in(status,"1|2") AND create_time >=1579449600000 AND create_time<=1582127999000&&sort=-pv`,
		"fetch_fields": "id;source_id;source_name",
	}

	resp, err := c.Request(appName, params, map[string]string{})
	t.Log(err, resp)
}
