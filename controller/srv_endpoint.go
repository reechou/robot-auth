package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-auth/models"
	"github.com/reechou/robot-auth/proto"
	"github.com/satori/go.uuid"
)

func (self *Logic) CreateRobotAuth(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()

	if r.Method != "POST" {
		return
	}

	req := &proto.CreateRobotAuthReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateRobotAuth json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	var authCodeList []string
	for i := 0; i < req.Num; i++ {
		u := uuid.NewV4()
		authCodeList = append(authCodeList, u.String())
	}

	now := time.Now().Unix()
	endTime := now + int64(86400*req.ExpireDate)
	var authList []models.RobotAuth
	for _, v := range authCodeList {
		authList = append(authList, models.RobotAuth{
			AuthCode:  v,
			CreatedAt: now,
			UpdatedAt: now,
			EndTime:   endTime,
		})
	}
	err := models.CreateRobotAuthList(authList)
	if err != nil {
		rsp.Code = proto.RESPONSE_ERR
		rsp.Msg = fmt.Sprintf("create auth list of num[%d] error", req.Num)
		return
	}

	rsp.Data = authCodeList
}

func (self *Logic) CheckRobotAuth(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()

	if r.Method != "POST" {
		return
	}

	req := &proto.CheckRobotAuthReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CheckRobotAuth json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	ra := &models.RobotAuth{
		AuthCode: req.AuthCode,
	}
	has, err := models.GetRobotAuth(ra)
	if err != nil {
		holmes.Error("get robot auth error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if ra.MachineCode == "" {
		ra.MachineCode = req.MachineCode
		err = models.UpdateRobotAuthMachineForce(ra)
		if err != nil {
			holmes.Error("update robot auth machine error: %v", err)
			rsp.Code = proto.RESPONSE_ERR
		}
		return
	}
	if ra.MachineCode != req.MachineCode {
		holmes.Error("check robot auth db auth_code[%s] != req auth_code[%s]", ra.MachineCode, req.MachineCode)
		rsp.Code = proto.RESPONSE_ERR
	}
}

func (self *Logic) ResetRobotAuth(w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		WriteJSON(w, http.StatusOK, rsp)
	}()

	if r.Method != "POST" {
		return
	}

	req := &proto.CheckRobotAuthReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("ResetRobotAuth json decode error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	ra := &models.RobotAuth{
		AuthCode: req.AuthCode,
	}
	err := models.UpdateRobotAuthMachine(ra)
	if err != nil {
		holmes.Error("reset robot auth error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
	}
}
