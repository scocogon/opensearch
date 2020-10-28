package osearch

// ParamHandle 未实现的参数，由此接口接入
//   key => source
type ParamHandle interface {
	// 字段名
	Key() string

	Source
}

// Source 已实现的参数
type Source interface {
	// 字段值
	Source() string
}
