package responses

const (
	SUCCESS           = 0
	UnknownError      = 10000
	ParameterError    = 10001
	OperationFailed   = 10002
	MissingParameters = 10003
	InvalidOperation  = 10004
	UserNotAuthorize  = 10005
	FrequencyTooFast  = 10006
	DataAlreadyExists = 10007
	DataDoesNotExist  = 10008
	UntrustedSource   = 10009

	UserAlreadyExists = 20000
	LoginExpires   = 21000
)

func GetMsg(code int) string {
	switch code {
	case UserAlreadyExists:
		return "User already exists"
	case LoginExpires:
		return "The login expires"
	case SUCCESS:
		return "SUCCESS"
	}
	return ""
}
