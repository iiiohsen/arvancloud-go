package errors

type Error interface {
	error
	HttpStatus() int
	ErrorCode() string
	Message() string
	OriginError() error
}
