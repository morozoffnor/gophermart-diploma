package accrual

type ErrorTooManyRequests struct {
	Timeout int
}

func (e *ErrorTooManyRequests) Error() string {
	return "too many requests"
}

type ErrorNotRegistered struct{}

func (e *ErrorNotRegistered) Error() string {
	return "order not registered in accrual system"
}

type ErrorInternalError struct{}

func (e *ErrorInternalError) Error() string {
	return "internal error"
}
