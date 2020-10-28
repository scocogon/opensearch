### 说明

+ 对 opensearch v3 api 简易封装
+ [官网说明](https://help.aliyun.com/document_detail/57155.html?spm=a2c4g.11174283.6.638.b1db5a19qzARRK)

### 安装

```
go get -u git.qutoutiao.net/x/qopensearch
```

### 示例

```golang
const (
	accessKeyID     = ""
	accessKeySecret = ""
	host            = "http://intranet.opensearch-cn-qingdao.aliyuncs.com"
	appName         = "qa_qtt_content"
)

func main() {
	cli := New(accessKeyID, accessKeySecret, host)

    query := clause.NewQuery()
    query.AddRangeII("create_time", 1579449600000, 1582127999000)

    cfg := clause.NewConfig().
        // hit=20
        SetHit(20).
        // rerank_size=100
        SetRerankSize(100)

    filter := clause.NewFilter()
    // content_type=3
    filter.AddIntEQ("content_type", 3).
        // in(status, "1|2")
        AddFnc("in", "status", 1, 2)

    srt := clause.NewSort()
    // -pv
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
    fmt.Printf("%s: %+v\n", err, resp)
}
```