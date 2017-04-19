package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type CreateRobotAuthReq struct {
	Num int `json:"num"`
}

type CheckRobotAuthReq struct {
	AuthCode string `json:"authCode"`
	MachineCode string `json:"machineCode"`
}
