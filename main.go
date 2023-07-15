package main

import (
	"cloud/logagent/conf"
	"cloud/logagent/util"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"os"
	"path"
)

const (
	fileMod = 0644
)

type ServerInfo struct {
	Server string `json:"server"`
}
type Server struct {
	gnet.BuiltinEventEngine
	eng     gnet.Engine
	fileMap map[string]*os.File
}

func (s *Server) OnBoot(e gnet.Engine) (action gnet.Action) {
	s.eng = e
	Logger.Sugar().Infof("logagent Server is listening on:%s", conf.AppConf.Port)
	return
}

func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	buf, err := c.Next(-1)
	if err != nil {
		Logger.Sugar().Errorf("OnTraffic err:%s", err.Error())
		return gnet.None
	}

	s.handlerLog(buf)
	return gnet.None
}
func (s *Server) handlerLog(buf []byte) {
	//logSrc := util.BytesToString(buf)
	//fmt.Printf(logSrc)
	split := func(s []byte) (ok bool, s1 []byte, s2 []byte) {
		max := 64
		//fmt.Println(string(s[len(s)-2]))
		if s[len(buf)-2] == '}' {
			i := len(buf) - 1
			cnt := 0
			for ; i >= 0; i-- {
				if cnt >= max {
					return
				}
				if buf[i] == '{' {
					ok = true
					s1 = s[:i]
					s2 = s[i : len(buf)-1]
					fmt.Println(cnt)
					return
				}
				cnt++
			}
			return
		}
		return
	}

	ok, log, server := split(buf)
	if !ok {
		return
	}
	info := &ServerInfo{}
	err := json.Unmarshal(server, info)
	if err != nil {
		return
	}
	fmt.Println(info.Server)
	fmt.Println(string(log))
	s.writeIO(info.Server, log)
}
func (s *Server) writeIO(fileName string, text []byte) {
	file := path.Join(conf.AppConf.LogFilePath, fileName)
	f := s.fileMap[fileName]
	var err error
	if f == nil {
		f, err = os.OpenFile(file+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, fileMod)
		if err != nil {
			Logger.Sugar().Infof("os.OpenFile err:%s", err.Error())
			return
		}
		s.fileMap[fileName] = f
	}

	n, err := fmt.Fprintln(f, util.BytesToString(text))
	if err != nil {
		Logger.Sugar().Infof("f.Write err:%s", err.Error())
		return
	}
	if n != len(text) {
		Logger.Sugar().Warn(n, len(text))
	}
}
func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {

	return gnet.None
}

func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	Logger.Sugar().Infof("server:%s connect", c.RemoteAddr().String())
	return
}
func (s *Server) Run() {
	addr := fmt.Sprintf("udp4://:%s", conf.AppConf.Port)
	f := func() {
		err := gnet.Run(s, addr,
			gnet.WithMulticore(true),
			//gnet.WithSocketSendBuffer(conf.GameConfig.ConnWriteBuffer),
			//gnet.WithSocketRecvBuffer(conf.GameConfig.ConnWriteBuffer),
			//gnet.WithTCPKeepAlive()
		)
		panic(err)
	}
	f()
}
func main() {
	InitLog()
	conf.Init()
	s := Server{
		fileMap: map[string]*os.File{},
	}
	s.Run()
}
