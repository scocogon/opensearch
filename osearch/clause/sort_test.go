package clause

import "testing"

func TestSort_All(t *testing.T) {
	s := NewSort()
	t.Log(s.Source())

	s.Asc("pv")
	t.Log(s.Source())

	s.Desc("sum")
	t.Log(s.Source())

	s.AscSum("duration", "value")
	t.Log(s.Source())
}
