package oerrors

import "errors"

var (
	ErrQueryNotFound = errors.New("缺少 query 字段")
)
