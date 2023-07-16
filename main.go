package main

import (
	"cloud/logagent/conf"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
)

const (
	fileMod = 0644
)

type Server struct {
	gnet.BuiltinEventEngine
	eng     gnet.Engine
	fileMap map[string]*lumberjack.Logger
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
	s.writeIO(ctx)
}
func (s *Server) writeIO(ctx *Context) {
	lvl := s.getLogLevel(ctx.LogLevel)
	file := path.Join(conf.AppConf.LogFilePath, ctx.ServerName+"_"+lvl+".log")
	lumberJackLogger, ok := s.fileMap[ctx.ServerName]
	if !ok {
		lumberJackLogger = &lumberjack.Logger{
			Filename:   file,                    // 文件位置
			MaxSize:    conf.AppConf.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
			MaxAge:     conf.AppConf.MaxAge,     // 保留旧文件的最大天数
			MaxBackups: conf.AppConf.MaxBackups, // 保留旧文件的最大个数
			Compress:   conf.AppConf.Compress,   // 是否压缩/归档旧文件
			LocalTime:  true,
		}
		s.fileMap[ctx.ServerName] = lumberJackLogger
	}

	n, err := lumberJackLogger.Write(ctx.Payload)
	if err != nil {
		panic(err)
	}
	if n != len(ctx.Payload) {
		Logger.Sugar().Warn(n, len(ctx.Payload))
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
		fileMap: map[string]*lumberjack.Logger{},
	}
	s.Run()
}
