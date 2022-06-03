package packet

type ErrorCode int
type Protocol int

const (
	START_BY_WEB = 2
)

const (
	Unknown ErrorCode = iota - 1
	Success
	NotFoundParam
	NotFoundRoom
)

const (
	None Protocol = iota
	Login
	Logout
	CreateRoom
	LookUpRoom
	Match
	CancelMatch
	ClientCount
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

func ContainsParamReq(res *ResponsePacket, req *RequestPacket, paramKey string) bool {
	var param = req.Param[paramKey]
	if param == nil {
		res.Error = NotFoundParam
		return false
	}
	return true
}

func ContainsParamRes(res *ResponsePacket, paramKey string) bool {
	var param = res.Param[paramKey]
	if param == nil {
		res.Error = NotFoundParam
		return false
	}
	return true
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
