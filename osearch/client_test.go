package osearch

import (
	"testing"

	"github.com/scocogon/opensearch/osearch/clause"
)

const (
	accessKeyID     = ""
	accessKeySecret = ""
	host            = "http://intranet.opensearch-cn-qingdao.aliyuncs.com"
	appName         = ""
)

// go test -timeout 30s -tags qa github.com/scocogon/opensearch/osearch -run "^(TestClient_All)$" -v -count=1
func TestClient_All(t *testing.T) {
	cli := New(accessKeyID, accessKeySecret, host)

	query := clause.NewQuery()
	query.AddRangeII("create_time", 1579449600000, 1582127999000)

	cfg := clause.NewConfig().
		SetHit(20).
		SetRerankSize(100)

	filter := clause.NewFilter()
	filter.AddIntEQ("content_type", 3).
		AddFnc("in", "status", 1, 2)

	srt := clause.NewSort()
	srt.Desc("pv")

	// 添加 query 子句
	// query=query=..&&filter=..&&sort=..&config=..
	cli.AddQuerys(query, filter, srt, cfg).

		// 添加 fetch_fields
		// id;source;source_type;comment_pv;bitrate
		AddFetchFields("id;source").
		AddFetchFields("source_type").
		AddFetchFields("comment_pv", "bitrate").

		// 添加其他参数
		// nofield=
		AddOtherParams(map[string]string{
			"nofield": "",
		})

	resp, err := cli.Send(appName, nil)
	t.Log(err, resp)
}
