package errors

import "fmt"

const (
	DefaultClientErrorStatus = 400
	DefaultClientErrorCode   = "SDK.ClientError"

	UnsupportedCredentialErrorCode    = "SDK.UnsupportedCredential"
	UnsupportedCredentialErrorMessage = "Specified credential (type = %s) is not supported, please check"

	CanNotResolveEndpointErrorCode    = "SDK.CanNotResolveEndpoint"
	CanNotResolveEndpointErrorMessage = "Can not resolve endpoint(param = %s), please check your accessKey with secret, and read the user guide\n %s"

	UnsupportedParamPositionErrorCode    = "SDK.UnsupportedParamPosition"
	UnsupportedParamPositionErrorMessage = "Specified param position (%s) is not supported, please upgrade sdk and retry"

	AsyncFunctionNotEnabledCode    = "SDK.AsyncFunctionNotEnabled"
	AsyncFunctionNotEnabledMessage = "Async function is not enabled in client, please invoke 'client.EnableAsync' function"

	UnknownRequestTypeErrorCode    = "SDK.UnknownRequestType"
	UnknownRequestTypeErrorMessage = "Unknown Request Type: %s"

	MissingParamErrorCode = "SDK.MissingParam"
	InvalidParamErrorCode = "SDK.InvalidParam"

	JsonUnmarshalErrorCode    = "SDK.JsonUnmarshalError"
	JsonUnmarshalErrorMessage = "Failed to unmarshal response, but you can get the data via response.GetHttpStatusCode() and response.GetHttpContentString()"

	TimeoutErrorCode    = "SDK.TimeoutError"
	TimeoutErrorMessage = "The request timed out %s times(%s for retry), perhaps we should have the threshold raised a little?"

	EmptyCredentialsCode    = "SDK.EmptyCredentials"
	EmptyCredentialsMessage = "invalid credentials: key & email must not be empty"
)

type ClientError struct {
	errorCode   string
	message     string
	originError error
}

func NewClientError(errorCode, message string, originErr error) Error {
	return &ClientError{
		errorCode:   errorCode,
		message:     message,
		originError: originErr,
	}
}

func (err *ClientError) Error() string {
	clientErrMsg := fmt.Sprintf("[%s] %s", err.ErrorCode(), err.message)
	if err.originError != nil {
		return clientErrMsg + "\ncaused by:\n" + err.originError.Error()
	}
	return clientErrMsg
}

func (err *ClientError) OriginError() error {
	return err.originError
}

func (*ClientError) HttpStatus() int {
	return DefaultClientErrorStatus
}

func (err *ClientError) ErrorCode() string {
	if err.errorCode == "" {
		return DefaultClientErrorCode
	} else {
		return err.errorCode
	}
}

func (err *ClientError) Message() string {
	return err.message
}

func (err *ClientError) String() string {
	return err.Error()
}
