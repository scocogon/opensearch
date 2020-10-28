package clause

import "testing"

func TestQuery_All(t *testing.T) {
	q := NewQuery()
	q.AddString("title", "中关村").
		AddRangeEE("created_at", nil, 100).
		AddRangeII("create_time", 40, 1000)

	q2 := NewQuery()
	q2.AddRangeEI("created_at", 1000, 3000)

	q3 := NewQuery()
	q3.AddString("area", "中国")

	q2 = q2.Or(q3)
	src := q2.Source()
	if src != `query=created_at:(1000,3000] OR area:"中国"` {
		t.Fatal(src)
	}

	q = q.And(q2.Or(q3))
	src = q.Source()

	if src != `query=(title:"中关村" AND created_at:(,100) AND create_time:[40,1000]) AND (created_at:(1000,3000] OR area:"中国" OR area:"中国")` {
		t.Fatal(src)
	}

	t.Log("DONE")
}
