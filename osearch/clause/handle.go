package clause

const (
	opAnd    = " AND "
	opOr     = " OR "
	opAndNot = " ANDNOT "
	opRank   = " RANK "

	// 排序
	opAsc  = "+"
	opDesc = "-"
)

type sourceHandle interface {
	source() string
}
