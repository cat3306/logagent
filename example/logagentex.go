package main

import (
	"encoding/binary"
	"fmt"
	"go.uber.org/zap"
	"time"

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

func main() {

	TestServer_OnOpen()

}

var (
	Chan         chan []byte
	packetEndian = binary.LittleEndian
)

type LogAgent interface {
	Write(p []byte) (n int, err error)
	LogLevelStrToUint(text []byte) (uint16, error)
	EnCode(payload []byte) []byte
}

func TestServer_OnOpen() {
	c, err := net.Dial("udp", ":8899")
	if err != nil {
		fmt.Println(err)
	}
	Chan = make(chan []byte, 10000)
	logger := initLogger()
	go func() {
		for {
			//	logger.Sugar().Warn("haha")
			logger.Sugar().Warn("haha")
			time.Sleep(time.Millisecond)
		}
	}()
	go handler(c)
	select {}
}
func handler(c net.Conn) {
	for b := range Chan {

		n, err := c.Write(b)
		if err != nil {
			panic(err)
		}
		if n != len(b) {
			fmt.Println(err)
		}
	}
}
func initLogger() *zap.Logger {
	encoder := zapcore.EncoderConfig{
		MessageKey:  "M",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalColorLevelEncoder, // INFO

		TimeKey:    "time",
		EncodeTime: zapcore.TimeEncoderOfLayout(Timeformat),

		CallerKey:        "caller",
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
		FunctionKey:      "func",
		NameKey:          "N",
	}

	consoleEncode := zapcore.NewConsoleEncoder(encoder)
	//consoleEncode.AddString("server", "TESTSERVER")
	w := zapcore.AddSync(new(ZapLoggerAgent).init("AKJZJ"))
	core := zapcore.NewCore(consoleEncode, w, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}

type ZapLoggerAgent struct {
	ServerName string
	color      map[string]zapcore.Level
}

func (l *ZapLoggerAgent) init(server string) *ZapLoggerAgent {
	l.ServerName = server
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
	return l
}
func (l *ZapLoggerAgent) Write(p []byte) (n int, err error) {
	pkg := l.EnCode(p)
	select {
	case Chan <- pkg:
	default:
		fmt.Printf(string(p), "asdads")
	}
	return len(p), nil
}

func (l *ZapLoggerAgent) EnCode(payload []byte) []byte {
	if l.ServerName == "" {
		panic("ServerName invalid")
	}

	hl := uint32(len(l.ServerName)) + logLevelLen
	buf := make([]byte, uint32(len(payload))+hl+headerLen)
	packetEndian.PutUint32(buf, hl)
	packetEndian.PutUint16(buf[headerLen:], l.LogLevelStrToUint(payload))
	copy(buf[headerLen+logLevelLen:], l.ServerName)
	copy(buf[headerLen+hl:], payload)

	return buf
}
func (l *ZapLoggerAgent) LogLevelStrToUint(text []byte) uint16 {
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
