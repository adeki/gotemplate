package errcode

type Code uint32

type CodeStruct struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	InvalidArguments Code = 1
	NotFound         Code = 2
	UnAuthorized     Code = 3
	Internal         Code = 4
)

func (c Code) String() string {
	switch c {
	case InvalidArguments:
		return "InvalidArguments"
	case NotFound:
		return "NotFound"
	case UnAuthorized:
		return "UnAuthorized"
	case Internal:
		return "Internal"
	default:
		return "UnknownCode"
	}
}

func (c Code) Message() string {
	switch c {
	case InvalidArguments:
		return "invalid arguments"
	case NotFound:
		return "not found"
	case UnAuthorized:
		return "unauthorized"
	case Internal:
		return "internal error"
	default:
		return "unknown error"
	}
}

func (c Code) Struct() CodeStruct {
	return CodeStruct{
		Code:    c.String(),
		Message: c.Message(),
	}
}
