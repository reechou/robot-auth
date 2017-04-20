package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type RobotAuth struct {
	ID          int64  `xorm:"pk autoincr" json:"id"`
	AuthCode    string `xorm:"not null default '' varchar(128) unique" json:"authCode"`
	MachineCode string `xorm:"not null default '' varchar(1024)" json:"machineCode"`
	TempUri     string `xorm:"not null default '' varchar(64)" json:"tempUri"`
	IfUseUri    int64  `xorm:"not null default 0 int" json:"ifUseUri"`
	IfAuth      int64  `xorm:"not null default 0 int" json:"ifAuth"`
	AuthTime    int64  `xorm:"not null default 0 int" json:"authTime"`
	EndTime     int64  `xorm:"not null default 0 int" json:"endTime"`
	CreatedAt   int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt   int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateRobotAuth(info *RobotAuth) error {
	if info.AuthCode == "" {
		return fmt.Errorf("robot autn[%s] cannot be nil.", info.AuthCode)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create robot auth error: %v", err)
		return err
	}
	holmes.Info("create robot auth[%v] success.", info)

	return nil
}

func CreateRobotAuthList(list []RobotAuth) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create robot auth list error: %v", err)
		return err
	}
	return nil
}

func GetRobotAuth(info *RobotAuth) (bool, error) {
	has, err := x.Where("auth_code = ?", info.AuthCode).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find robot auth from auth[%s]", info.AuthCode)
		return false, nil
	}
	return true, nil
}

func UpdateRobotAuthTempUri(info *RobotAuth) error {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("temp_uri", "if_use_uri", "updated_at").Update(info)
	if affected == 0 {
		return fmt.Errorf("auth update uri error")
	}
	return err
}

func UpdateRobotAuthTempUriIfUse(info *RobotAuth) error {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("if_use_uri", "updated_at").Update(info)
	if affected == 0 {
		return fmt.Errorf("auth update uri if use error")
	}
	return err
}

func UpdateRobotAuthMachine(info *RobotAuth) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("machine_code", "updated_at").Update(info)
	return err
}

func UpdateRobotAuthMachineForce(info *RobotAuth) error {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("machine_code", "if_use_uri", "updated_at").Where("machine_code = ''").Update(info)
	if affected == 0 {
		return fmt.Errorf("auth[%s] has bind machine", info.AuthCode)
	}
	return err
}

func UpdateRobotAuthMachineFirst(info *RobotAuth) error {
	info.UpdatedAt = time.Now().Unix()
	affected, err := x.ID(info.ID).Cols("machine_code", "if_use_uri", "if_auth", "end_time", "updated_at").Where("machine_code = ''").And("if_auth = 0").Update(info)
	if affected == 0 {
		return fmt.Errorf("auth[%s] has bind machine or has authed", info.AuthCode)
	}
	return err
}
