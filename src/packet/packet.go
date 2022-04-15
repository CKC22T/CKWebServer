package packet

type ErrorCode int
type Protocol int

const (
	START_BY_WEB = 2
)

const (
	Unknown ErrorCode = iota - 1
	Success
)

const (
	None Protocol = iota
	Login
	Logout
	CreateRoom
	LookUpRoom
	Match
	CancelMatch
)

type RequestPacket struct {
	Code  Protocol               `json:"code"`
	Param map[string]interface{} `json:"param"`
}

func NewRequestPacket() *RequestPacket {
	return &RequestPacket{
		Code:  None,
		Param: make(map[string]interface{}),
	}
}

type ResponsePacket struct {
	Error ErrorCode              `json:"error"`
	Code  Protocol               `json:"code"`
	Param map[string]interface{} `json:"param"`
}

func NewResponsePacket() *ResponsePacket {
	return &ResponsePacket{
		Error: Success,
		Code:  None,
		Param: make(map[string]interface{}),
	}
}
