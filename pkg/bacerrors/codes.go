package bacerrors

type ErrorCode string

const (
	BadRequestError    ErrorCode = "BadRequest"
	InternalError      ErrorCode = "InternalError"
	NotFoundError      ErrorCode = "NotFound"
	TimeOutError       ErrorCode = "TimeOut"
	UnauthorizedError  ErrorCode = "Unauthorized"
	Forbidden          ErrorCode = "Forbidden"
	ServiceUnavailable ErrorCode = "ServiceUnavailable"
	NotImplemented     ErrorCode = "NotImplemented"
	ResourceExhausted  ErrorCode = "ResourceExhausted"
	ResourceInUse      ErrorCode = "ResourceInUse"
	VersionMismatch    ErrorCode = "VersionMismatch"
	ValidationError    ErrorCode = "ValidationError"
	TooManyRequests    ErrorCode = "TooManyRequests"
	NetworkFailure     ErrorCode = "NetworkFailure"
	ConfigurationError ErrorCode = "ConfigurationError"
	DatastoreFailure   ErrorCode = "DatastoreFailure"
	RequestCancelled   ErrorCode = "RequestCancelled"
	IOError            ErrorCode = "IOError"
	UnknownError       ErrorCode = "UnknownError"
)

func Code(code string) ErrorCode {
	return ErrorCode(code)
}
