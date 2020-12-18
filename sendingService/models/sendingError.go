package models

//SendingError is a structure containing sending failure type
type SendingError struct {
	ErrorType int
	Message   *string
}

func (e *SendingError) Error() string {
	switch e.ErrorType {
	case -1:
		return "Failed logical job"
	case -2:
		return "Failed repeatable job"
	default:
		return "Fatal error"
	}
}
