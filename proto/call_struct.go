package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
	RESPONSE_EXPIRED
	RESPONSE_SYSTEM
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type CreateRobotAuthReq struct {
	Num        int `json:"num"`
	ExpiryDate int `json:"expiryDate"` // day
}

type CheckRobotAuthReq struct {
	AuthCode    string `json:"authCode"`
	MachineCode string `json:"machineCode"`
	Timestamp   int64  `json:"timestamp"`
}

type CheckRobotAuthRsp struct {
	EndTime   int64  `json:"endTime"`
	SecretKey string `json:"secretKey"`
}
