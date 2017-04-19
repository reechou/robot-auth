package controller

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/robot-auth/proto"
)

const (
	AuthPrefix = "/auth"
)

type HandlerRequest struct {
	Method string
	Path   string
	Val    []byte
}

type AuthHandler struct {
	l *Logic
}

func (self *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr, err := parseRequest(r)
	if err != nil {
		holmes.Error("parse request error: %v", err)
		writeRsp(w, &proto.Response{Code: proto.RESPONSE_ERR})
		return
	}

	rsp := self.l.doAuth(rr)
	writeRsp(w, rsp)
}

func parseRequest(r *http.Request) (*HandlerRequest, error) {
	req := &HandlerRequest{}
	req.Path = r.URL.Path[len(AuthPrefix)+1:]

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, errors.New("parse request read error")
	}
	r.Body.Close()

	req.Method = r.Method
	req.Val = result

	return req, nil
}

func writeRsp(w http.ResponseWriter, rsp *proto.Response) {
	w.Header().Set("Content-Type", "application/json")

	if rsp != nil {
		WriteJSON(w, http.StatusOK, rsp)
	}
}
