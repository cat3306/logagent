package main

import (
	"bytes"
	"cloud/logagent/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type errMsg struct {
	serverName string
	text       string
}

func (s *Server) handleErrMsg() {
	for item := range s.errMsgChan {
		list, ok := s.errMsgMap[item.serverName]
		if !ok {
			list = make([]string, 0, s.capMsg)
		}
		list = append(list, item.text)
		if len(list) >= s.capMsg {
			s.sendMsg(list, item.serverName)
			list = list[0:0]
		}
		s.errMsgMap[item.serverName] = list
	}
}

type dingDingReq struct {
	MsgType  string `json:"msgtype"`
	MarkDown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}

func (s *Server) sendMsg(msgs []string, sName string) {
	reqJson := dingDingReq{
		MsgType: "markdown",
		MarkDown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: "游戏版本更新通知",
		},
	}
	reqJson.MarkDown.Text = fmt.Sprintf("## 环境:%s,服务:%s\n", conf.AppConf.NameSpace, sName) + strings.Join(msgs, "\n")
	reqRaw, err := json.Marshal(reqJson)
	if err != nil {
		Logger.Sugar().Errorf("sendMsg Marshal err:%s", err.Error())
		return
	}
	b := bytes.NewBuffer(reqRaw)
	req, err := http.NewRequest(http.MethodPost, conf.AppConf.NotificationUrl, b)
	if err != nil {
		Logger.Sugar().Errorf("sendMsg err:%s", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		Logger.Sugar().Errorf("client.Do err:%s", err.Error())
		return
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		Logger.Sugar().Errorf("StatusCode:%d,Status:%s", rsp.StatusCode, rsp.Status)
		return
	}
	raw, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		Logger.Sugar().Errorf("ioutil.ReadAll:%s", err.Error())
		return
	}
	Logger.Sugar().Infof("sendMsg:%s", string(raw))
}
