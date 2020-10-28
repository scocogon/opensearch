package clause

import "testing"

func TestFilter_All(t *testing.T) {
	f := NewFilter()
	f.AddExpr("(hit+sale)*rate>10000").
		AddFloatEQ("geo", 10.24).
		AddIntLT("dis", 100)

	t.Log(f.Source())

	f1 := NewFilter()
	f1.AddFnc("in", "status", 1, 2)
	t.Log(f1.Source())

	f2 := NewFilter()
	f2.AddExpr("fieldlen(title)>1")
	t.Log(f2.Source())

	t.Log(f1.And(f.Or(f2)).Source())
}
