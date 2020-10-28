package clause

import (
	"fmt"
	"strings"
)

// Query query clause
type Query struct {
	src []sourceHandle
	op  string
}

// NewQuery new query
func NewQuery() *Query { return &Query{} }

func (q *Query) AddString(key string, value string) *Query {
	q.src = append(q.src, &qc{key: key, val: value})
	return q
}

func (q *Query) AddRangeII(key string, lval, rval interface{}) *Query {
	return q.addRange(optII, key, lval, rval)
}

func (q *Query) AddRangeIE(key string, lval, rval interface{}) *Query {
	return q.addRange(optIE, key, lval, rval)
}

func (q *Query) AddRangeEE(key string, lval, rval interface{}) *Query {
	return q.addRange(optEE, key, lval, rval)
}

func (q *Query) AddRangeEI(key string, lval, rval interface{}) *Query {
	return q.addRange(optEI, key, lval, rval)
}

func (q *Query) Or(r *Query) *Query     { return mergeQuery(opOr, q, r) }
func (q *Query) And(r *Query) *Query    { return mergeQuery(opAnd, q, r) }
func (q *Query) AndNot(r *Query) *Query { return mergeQuery(opAndNot, q, r) }
func (q *Query) Rank(r *Query) *Query   { return mergeQuery(opRank, q, r) }

func mergeQuery(op string, l, r *Query) *Query {
	return &Query{
		src: []sourceHandle{l, r},
		op:  op,
	}
}

func (q *Query) addRange(opt int, key string, lval, rval interface{}) *Query {
	q.src = append(q.src, &qc2{key: key, val: []interface{}{lval, rval}, opt: opt})
	return q
}

func (q *Query) Source() string {
	return "query=" + q.source()
}

func (q *Query) source() string {
	switch q.op {
	case opRank:
		return q.src[0].source() + q.op + q.src[1].source()

	case opAnd, opOr, opAndNot:
		return "(" + q.src[0].source() + ")" + q.op + "(" + q.src[1].source() + ")"

	default:
		var ss []string
		for _, v := range q.src {
			ss = append(ss, v.source())
		}

		return strings.Join(ss, opAnd)
	}
}

const (
	// 数组边界
	optEE = iota // 左右都不包含
	optIE        // 左包含，右不包含
	optEI        // 左不包含，右包含
	optII        // 左右都包含
)

type qc struct {
	key string
	val string // string
}

func (q *qc) source() string {
	return q.key + `:"` + q.val + `"`
}

type qc2 struct {
	key string
	val []interface{}

	opt int
}

func (q *qc2) source() string {
	var l, r string
	switch q.opt {
	case optII:
		l, r = "[", "]"
	case optIE:
		l, r = "[", ")"
	case optEE:
		l, r = "(", ")"
	case optEI:
		l, r = "(", "]"
	}

	switch {
	case q.val[0] == nil:
		return fmt.Sprintf(`%s:%s,%v%s`, q.key, l, q.val[1], r)
	case q.val[1] == nil:
		return fmt.Sprintf(`%s:%s%v,%s`, q.key, l, q.val[0], r)
	}

	return fmt.Sprintf(`%s:%s%v,%v%s`, q.key, l, q.val[0], q.val[1], r)
}
