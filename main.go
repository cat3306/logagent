package main

import (
	"cloud/logagent/conf"
	"cloud/logagent/util"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"os"
	"path"
)

const (
	fileMod = 0644
)

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
	ctx := Decode(buf)
	fmt.Println(s.getLogLevel(ctx.LogLevel))
	s.writeIO(ctx)
}
func (s *Server) writeIO(ctx *Context) {
	lvl := s.getLogLevel(ctx.LogLevel)
	file := path.Join(conf.AppConf.LogFilePath, ctx.ServerName+"_"+lvl)
	f := s.fileMap[ctx.ServerName]
	var err error
	if f == nil {
		f, err = os.OpenFile(file+".log", os.O_RDWR|os.O_APPEND|os.O_CREATE, fileMod)
		if err != nil {
			Logger.Sugar().Infof("os.OpenFile err:%s", err.Error())
			return
		}
		s.fileMap[ctx.ServerName] = f
	}

	n, err := fmt.Fprint(f, util.BytesToString(ctx.Payload))
	if err != nil {
		Logger.Sugar().Infof("f.Write err:%s", err.Error())
		return
	}
	fmt.Println(ctx.ServerName, ctx.LogLevel)
	if n != len(ctx.Payload) {
		//Logger.Sugar().Warn(n, len(text))
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
func (s *Server) getLogLevel(l uint16) string {

	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case DPanicLevel:
		return "DPANIC"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return ""
	}

}
func main() {
	InitLog()
	conf.Init()
	s := Server{
		fileMap: map[string]*os.File{},
	}
	s.Run()
}
