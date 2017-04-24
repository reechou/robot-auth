package controller

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/reechou/holmes"
	"github.com/reechou/robot-auth/models"
	"github.com/reechou/robot-auth/proto"
	"github.com/satori/go.uuid"
)

func (self *Logic) doAuth(rr *HandlerRequest) *proto.Response {
	switch rr.Path {
	case "QhWT4xJ1W7v5PRwV":
		return self.doCreateAuth(rr)
	case "uri":
		return self.doAuthUri(rr)
	case "reset_robot_auth":
		return self.doResetAuth(rr)
	case "get_commission":
		return self.doGetCommission(rr)
	default:
		return self.doCheckAuth(rr)
	}
}

func (self *Logic) doGetCommission(rr *HandlerRequest) *proto.Response {
	return &proto.Response{Code: proto.RESPONSE_OK, Data: self.cfg.CommissionRate}
}

func (self *Logic) doCreateAuth(rr *HandlerRequest) *proto.Response {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}

	req := &proto.CreateRobotAuthReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	var authCodeList []string
	for i := 0; i < req.Num; i++ {
		u := uuid.NewV4()
		authCodeList = append(authCodeList, u.String())
	}

	now := time.Now().Unix()
	//endTime := now + int64(86400*req.ExpiryDate)
	var authList []models.RobotAuth
	for _, v := range authCodeList {
		authList = append(authList, models.RobotAuth{
			AuthCode:  v,
			CreatedAt: now,
			UpdatedAt: now,
			AuthTime:  int64(req.ExpireDate),
		})
	}
	err = models.CreateRobotAuthList(authList)
	if err != nil {
		rsp.Code = proto.RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create auth list of num[%d] error", req.Num)
		return rsp
	}

	rsp.Data = authCodeList

	return rsp
}

func (self *Logic) doAuthUri(rr *HandlerRequest) *proto.Response {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}

	req := &proto.CheckRobotAuthReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}
	
	// check secret key
	secretKey := string(md5Of32(md5Of32([]byte(fmt.Sprintf("%s%d", req.AuthCode, req.Timestamp)))))
	if len(secretKey) < 10 {
		holmes.Error("len(secretKey)[%d] < 10", len(secretKey))
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}
	if secretKey[:10] != req.SecretKey {
		holmes.Error("secretKey[:10][%s] != req.SecretKey[%s]", secretKey[:10], req.SecretKey)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	ra := &models.RobotAuth{
		AuthCode: req.AuthCode,
	}
	has, err := models.GetRobotAuth(ra)
	if err != nil {
		holmes.Error("get robot auth error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}
	if !has {
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	tempUri := uniuri.New()
	ra.TempUri = tempUri
	ra.IfUseUri = 0
	err = models.UpdateRobotAuthTempUri(ra)
	if err != nil {
		holmes.Error("update temp uri error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}
	rsp.Data = tempUri

	return rsp
}

func (self *Logic) doCheckAuth(rr *HandlerRequest) *proto.Response {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}

	req := &proto.CheckRobotAuthReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	ra := &models.RobotAuth{
		AuthCode: req.AuthCode,
	}
	has, err := models.GetRobotAuth(ra)
	if err != nil {
		holmes.Error("get robot auth error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}
	if !has {
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	now := time.Now().Unix()

	md5Key := fmt.Sprintf("%s%s", ra.AuthCode, ra.TempUri)
	allSecretKey := string(md5Of32(md5Of32([]byte(md5Key))))
	if len(allSecretKey) < 8 {
		holmes.Error("md5 error")
		rsp.Code = proto.RESPONSE_SYSTEM
		return rsp
	}
	realSecretKey := allSecretKey[:8]

	if realSecretKey != rr.Path || ra.IfUseUri != 0 {
		holmes.Error("uri check error: %s %s or uri has used[%d]", realSecretKey, rr.Path, ra.IfUseUri)
		rsp.Code = proto.RESPONSE_ERR
		rsp.Msg = "uri error or uri has used."
		return rsp
	}

	if ra.MachineCode == "" {
		ra.MachineCode = req.MachineCode
		ra.IfUseUri = 1
		if ra.IfAuth == 0 {
			ra.IfAuth = 1
			ra.EndTime = now + 86400*ra.AuthTime
			err = models.UpdateRobotAuthMachineFirst(ra)
			if err != nil {
				holmes.Error("update robot first auth machine error: %v", err)
				rsp.Code = proto.RESPONSE_ERR
				return rsp
			}
		} else {
			err = models.UpdateRobotAuthMachineForce(ra)
			if err != nil {
				holmes.Error("update robot auth machine error: %v", err)
				rsp.Code = proto.RESPONSE_ERR
				return rsp
			}
		}
	} else {
		if ra.MachineCode != req.MachineCode {
			holmes.Error("check robot auth db machine_code[%s] != req machine_code[%s]", ra.MachineCode, req.MachineCode)
			rsp.Code = proto.RESPONSE_ERR
			return rsp
		}
		ra.IfUseUri = 1
		models.UpdateRobotAuthTempUriIfUse(ra)
	}
	
	if now > ra.EndTime {
		rsp.Code = proto.RESPONSE_EXPIRED
		return rsp
	}

	rspKey := string(md5Of32(md5Of32(md5Of32([]byte(fmt.Sprintf("%s%d", ra.TempUri, req.Timestamp))))))
	rsp.Data = &proto.CheckRobotAuthRsp{
		EndTime:   ra.EndTime,
		SecretKey: rspKey,
	}

	return rsp
}

func (self *Logic) doResetAuth(rr *HandlerRequest) *proto.Response {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}

	req := &proto.CheckRobotAuthReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return rsp
	}

	ra := &models.RobotAuth{
		AuthCode: req.AuthCode,
	}
	err = models.UpdateRobotAuthMachine(ra)
	if err != nil {
		holmes.Error("reset robot auth error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
	}

	return rsp
}

func md5Of32(src []byte) []byte {
	hash := md5.New()
	hash.Write(src)
	cipherText2 := hash.Sum(nil)
	hexText := make([]byte, 32)
	hex.Encode(hexText, cipherText2)
	return hexText
}
