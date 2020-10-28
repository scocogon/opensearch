package clause

import "strings"

type Sort struct {
	src []sourceHandle
}

func NewSort() *Sort {
	return &Sort{}
}

func (s *Sort) Asc(field string) *Sort {
	return s.sc(opAsc, field)
}

// AscSum 多字段求和后升序
func (s *Sort) AscSum(f1 string, fs ...string) *Sort {
	return s.sc(opAsc, f1, fs...)
}

func (s *Sort) Desc(field string) *Sort {
	return s.sc(opDesc, field)
}

// DescSum 多字段求和后降序
func (s *Sort) DescSum(f1 string, fs ...string) *Sort {
	return s.sc(opDesc, f1, fs...)
}

func (s *Sort) sc(op string, field string, fs ...string) *Sort {
	if len(fs) > 0 {
		field = "(" + field + "+" + strings.Join(fs, "+") + ")"
	}

	s.src = append(s.src, &sc{key: field, op: op})
	return s
}

func (s *Sort) Source() string {
	return "sort=" + s.source()
}

func (s *Sort) source() string {
	var ss []string
	for _, v := range s.src {
		ss = append(ss, v.source())
	}

	return strings.Join(ss, ";")
}

type sc struct {
	key string
	op  string
}

func (s *sc) source() string {
	return s.op + s.key
}
