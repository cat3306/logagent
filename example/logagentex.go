package main

import (
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"strings"
)

//copy go.uber.org/zap/internal/color
// Foreground colors.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color represents a text color.
type Color uint8

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

//copy go.uber.org/zap/internal/color
const (
	headerLen   = uint32(4)
	logLevelLen = uint32(2)
	Timeformat  = "2006-01-02 15:04:05.000"
)

type LogAgent interface {
	Write(p []byte) (n int, err error)
	LogLevelStrToUint(text []byte) (uint16, error)
	EnCode(payload []byte) []byte
}
type LogAgentConf struct {
	ServerName  string
	AgentAddr   string
	ChanBuffer  int
	EncoderConf *zapcore.EncoderConfig
}
type ZapLoggerAgent struct {
	conf        *LogAgentConf
	color       map[string]zapcore.Level
	logger      *zap.Logger
	bufferChan  chan []byte
	c           net.Conn
	encoderConf *zapcore.EncoderConfig
	offlineMode bool
}

func (l *ZapLoggerAgent) Demons() *ZapLoggerAgent {
	go func() {
		for b := range l.bufferChan {
			_, err := l.c.Write(b)
			if err != nil {
				fmt.Printf(string(b))
				continue
			}
		}
	}()
	return l
}
func (l *ZapLoggerAgent) Conn() *ZapLoggerAgent {
	c, err := net.Dial("udp", l.conf.AgentAddr)
	if err != nil {
		fmt.Println(err)
		l.offlineMode = true
		return l
	}
	l.c = c
	return l
}
func (l *ZapLoggerAgent) initLogger() *ZapLoggerAgent {
	if l.conf.EncoderConf == nil {
		l.conf.EncoderConf = &zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder, // INFO

			TimeKey:    "time",
			EncodeTime: zapcore.TimeEncoderOfLayout(Timeformat),

			CallerKey:        "caller",
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
			FunctionKey:      "func",
		}
	}

	consoleEncode := zapcore.NewConsoleEncoder(*l.conf.EncoderConf)
	w := zapcore.AddSync(l)
	core := zapcore.NewCore(consoleEncode, w, zapcore.DebugLevel)
	l.logger = zap.New(core, zap.AddCaller())
	return l
}
func (l *ZapLoggerAgent) Init(config *LogAgentConf) *ZapLoggerAgent {
	if config == nil {
		panic("config nil")
	}
	if config.ServerName == "" {
		panic("ServerName invalid")
	}
	if config.AgentAddr == "" {
		panic("AgentAddr invalid")
	}
	l.conf = config
	l.color = map[string]zapcore.Level{}
	debug := zapcore.DebugLevel.CapitalString()
	info := zapcore.InfoLevel.CapitalString()
	warn := zapcore.WarnLevel.CapitalString()
	er := zapcore.ErrorLevel.CapitalString()
	dPanic := zapcore.DPanicLevel.CapitalString()
	pan := zapcore.PanicLevel.CapitalString()
	fatal := zapcore.FatalLevel.CapitalString()
	l.color[Magenta.Add(debug)] = zapcore.DebugLevel
	l.color[Blue.Add(info)] = zapcore.InfoLevel
	l.color[Yellow.Add(warn)] = zapcore.WarnLevel
	l.color[Red.Add(er)] = zapcore.ErrorLevel
	l.color[Red.Add(dPanic)] = zapcore.DPanicLevel
	l.color[Red.Add(pan)] = zapcore.PanicLevel
	l.color[Red.Add(fatal)] = zapcore.FatalLevel
	l.initLogger()
	if l.conf.ChanBuffer == 0 {
		l.conf.ChanBuffer = 1024
	}
	l.bufferChan = make(chan []byte, l.conf.ChanBuffer)
	return l
}
func (l *ZapLoggerAgent) Write(p []byte) (n int, err error) {
	if l.offlineMode {
		fmt.Printf(string(p))
		return len(p), nil
	}
	pkg := l.EnCode(p)
	select {

	case l.bufferChan <- pkg:
	default:
		fmt.Printf(string(p))
	}
	return len(p), nil
}

func (l *ZapLoggerAgent) EnCode(payload []byte) []byte {
	if l.conf.ServerName == "" {
		panic("ServerName invalid")
	}

	hl := uint32(len(l.conf.ServerName)) + logLevelLen
	buf := make([]byte, uint32(len(payload))+hl+headerLen)
	binary.LittleEndian.PutUint32(buf, hl)
	binary.LittleEndian.PutUint16(buf[headerLen:], l.logLevelStrToUint(payload))
	copy(buf[headerLen+logLevelLen:], l.conf.ServerName)
	copy(buf[headerLen+hl:], payload)

	return buf
}
func (l *ZapLoggerAgent) logLevelStrToUint(text []byte) uint16 {
	timeLen := len(Timeformat)
	s := strings.Builder{}
	//timeLen+16 防止遍历整个text
	for i := timeLen + 1; i < timeLen+16; i++ {
		if text[i] == ' ' {
			break
		}
		s.WriteByte(text[i])
	}
	le := l.color[s.String()]
	le++ //zapcore.DebugLevel value -1 but to convert uint16,so +1
	return uint16(le)
}

func (l *ZapLoggerAgent) Logger() *zap.Logger {
	return l.logger
}
