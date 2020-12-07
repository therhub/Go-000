package model

// 自定义错误的类型
type EasyError struct {
	Err error
	Msg string
}

func (e *EasyError) Error() string {
	return e.Msg + e.Err.Error()
}

// 可以实现，也可以不实现
// errors.As是有判断unwrap
func (e *EasyError) Unwrap() error { return e.Err }
