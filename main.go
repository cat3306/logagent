package main

import (
	"cloud/logagent/conf"
	"cloud/logagent/internal"
	"fmt"
	"path"
	"sync"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	fileMod = 0644
)

type Server struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	fileMap      map[string]*lumberjack.Logger
	locker       sync.RWMutex
	gPool        *goroutine.Pool
	errMsgChan   chan *errMsg
	errMsgMap    map[string][]string
	capMsg       int
	m            map[string]int
	mainLoopChan chan *Context
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

	select {
	case s.mainLoopChan <- ctx:
	default:
		s.writeIO(ctx)
	}
}
func (s *Server) getLogger(fileName string) (*lumberjack.Logger, bool) {
	s.locker.RLock()
	defer s.locker.RUnlock()
	lg, ok := s.fileMap[fileName]
	return lg, ok
}
func (s *Server) setLogger(fileName string, lg *lumberjack.Logger) {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.fileMap[fileName] = lg
}
func (s *Server) writeIO(ctx *Context) {
	lvl := s.getLogLevel(ctx.LogLevel)
	serverName := ctx.ServerName + "_" + lvl + ".log"
	file := path.Join(conf.AppConf.LogFilePath, serverName)
	lumberJackLogger, ok := s.getLogger(serverName)
	if !ok {
		lumberJackLogger = &lumberjack.Logger{
			Filename:   file,                    // 文件位置
			MaxSize:    conf.AppConf.MaxSize,    // 进行切割之前,日志文件的最大大小(MB为单位)
			MaxAge:     conf.AppConf.MaxAge,     // 保留旧文件的最大天数
			MaxBackups: conf.AppConf.MaxBackups, // 保留旧文件的最大个数
			Compress:   conf.AppConf.Compress,   // 是否压缩/归档旧文件
			LocalTime:  true,
		}
		s.setLogger(serverName, lumberJackLogger)
	}
	payload := ctx.Payload[:ctx.PayloadLen]
	f := func() {
		n, err := lumberJackLogger.Write(payload)
		if err != nil {
			Logger.Sugar().Errorf("err:%s", err.Error())
			return
		}
		if n != len(payload) {
			Logger.Sugar().Warn(n, len(payload))
		}
	}
	// if ctx.LogLevel == ErrorLevel {
	// 	s.errMsgChan <- &errMsg{
	// 		serverName: ctx.ServerName,
	// 		text:       string(payload),
	// 	}
	// }
	//不发钉钉消息了
	f()
	internal.BUFFERPOOL.Put(ctx.Payload)
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
func (s *Server) mainLoop() {
	for ctx := range s.mainLoopChan {
		s.writeIO(ctx)
	}
}
func main() {
	InitLog()
	conf.Init()
	s := Server{
		fileMap:      map[string]*lumberjack.Logger{},
		gPool:        goroutine.Default(),
		errMsgChan:   make(chan *errMsg, 1024),
		errMsgMap:    make(map[string][]string),
		capMsg:       64,
		m:            map[string]int{},
		mainLoopChan: make(chan *Context, 1024),
	}
	go s.handleErrMsg()
	go s.mainLoop()
	s.Run()
}
