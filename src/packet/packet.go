package packet

type ErrorCode int
type Protocol int

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
	StartMatch
	CancelMatch
)

type RequestPacket struct {
	Code  Protocol               `json:"code"`
	Param map[string]interface{} `json:"param"`
}

func NewRequestPacket() *RequestPacket {
	return &RequestPacket{
		Code:  None,
		Param: map[string]interface{}{},
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
		Param: map[string]interface{}{},
	}
}
