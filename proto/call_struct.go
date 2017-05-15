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
	ExpireDate int `json:"expiryDate"` // day
}

type CheckRobotAuthReq struct {
	AuthCode    string `json:"authCode"`
	MachineCode string `json:"machineCode"`
	Timestamp   int64  `json:"timestamp"`
	SecretKey   string `json:"secretKey"`
}

type CheckRobotAuthRsp struct {
	EndTime   int64  `json:"endTime"`
	SecretKey string `json:"secretKey"`
}

type CheckUpdateRsp struct {
	Version string `json:"version"`
	Url     string `json:"url"`
}
