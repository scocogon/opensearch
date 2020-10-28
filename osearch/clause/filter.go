package clause

import (
	"fmt"
	"strings"
)

const (
	ExprEQ = "="
	ExprLT = "<"
	ExprLE = "<="
	ExprGT = ">"
	ExprGE = ">="
	ExprNE = "!="
)

type Filter struct {
	src []sourceHandle
	op  string
}

func NewFilter() *Filter { return &Filter{} }

func (f *Filter) AddStringEQ(key string, val string) *Filter {
	return f.addExpr1(ExprEQ, key, val)
}

func (f *Filter) AddStringNE(key string, val string) *Filter {
	return f.addExpr1(ExprNE, key, val)
}

func (f *Filter) AddIntEQ(key string, val int) *Filter {
	return f.addExpr2(ExprEQ, key, val)
}

func (f *Filter) AddIntNE(key string, val int) *Filter {
	return f.addExpr2(ExprNE, key, val)
}

func (f *Filter) AddIntLT(key string, val int) *Filter {
	return f.addExpr2(ExprLT, key, val)
}

func (f *Filter) AddIntLE(key string, val int) *Filter {
	return f.addExpr2(ExprLE, key, val)
}

func (f *Filter) AddIntGT(key string, val int) *Filter {
	return f.addExpr2(ExprGT, key, val)
}

func (f *Filter) AddIntGE(key string, val int) *Filter {
	return f.addExpr2(ExprGE, key, val)
}

func (f *Filter) AddFloatEQ(key string, val float64) *Filter {
	return f.addExpr2(ExprEQ, key, val)
}

func (f *Filter) AddFloatNE(key string, val float64) *Filter {
	return f.addExpr2(ExprNE, key, val)
}

func (f *Filter) AddFloatLT(key string, val float64) *Filter {
	return f.addExpr2(ExprLT, key, val)
}

func (f *Filter) AddFloatLE(key string, val float64) *Filter {
	return f.addExpr2(ExprLE, key, val)
}

func (f *Filter) AddFloatGT(key string, val float64) *Filter {
	return f.addExpr2(ExprGT, key, val)
}

func (f *Filter) AddFloatGE(key string, val float64) *Filter {
	return f.addExpr2(ExprGE, key, val)
}

func (f *Filter) AddExpr(expr string) *Filter {
	f.src = append(f.src, &fc2{expr: strings.TrimSpace(expr)})
	return f
}

func (f *Filter) AddFnc(fnc string, args ...interface{}) *Filter {
	src := generateByFunc(fnc, args...)
	if src != nil {
		f.src = append(f.src, src)
	}

	return f
}

func (f *Filter) Or(r *Filter) *Filter  { return mergeFilter(opOr, f, r) }
func (f *Filter) And(r *Filter) *Filter { return mergeFilter(opAnd, f, r) }

func (f *Filter) addExpr1(op string, key string, val string) *Filter {
	f.src = append(f.src, &fc{key: key, val: val, op: op})
	return f
}

func mergeFilter(op string, l, r *Filter) *Filter {
	return &Filter{
		src: []sourceHandle{l, r},
		op:  op,
	}
}

func (f *Filter) addExpr2(op string, key string, val interface{}) *Filter {
	f.src = append(f.src, &fc1{key: key, val: val, op: op})
	return f
}

func (f *Filter) Source() string {
	return "filter=" + f.source()
}

func (f *Filter) source() string {
	switch f.op {
	case opOr:
		return f.src[0].source() + f.op + f.src[1].source()

	case opAnd:
		return "(" + f.src[0].source() + ")" + f.op + "(" + f.src[1].source() + ")"

	default:
		var ss []string
		for _, v := range f.src {
			ss = append(ss, v.source())
		}

		return strings.Join(ss, opAnd)
	}
}

// key=val (>,>=,<=,<,!=)
type fc struct {
	key string
	val string // string

	op string
}

func (f *fc) source() string {
	return f.key + f.op + `"` + f.val + `"`
}

type fc1 struct {
	key string
	val interface{} // float, int

	op string
}

func (f *fc1) source() string {
	return f.key + f.op + fmt.Sprintf("%+v", f.val)
}

type fc2 struct {
	expr string
}

func (f *fc2) source() string {
	return f.expr
}

// 函数
type fc3 struct {
	fnc string

	fc2
}

func generateByFunc(name string, args ...interface{}) sourceHandle {
	switch name {
	case "in", "notin":
		if len(args) < 2 {
			return nil
		}

		field, ok := args[0].(string)
		if !ok {
			return nil
		}

		var s []string
		for _, v := range args[1:] {
			switch v := v.(type) {
			case string:
				s = append(s, v)
			default:
				s = append(s, fmt.Sprintf("%v", v))
			}
		}

		return &fc2{expr: name + "(" + field + `,"` + strings.Join(s, `|`) + `")`}

	default:
		panic("PANIC: filter不支持的函数调用")
	}
}
