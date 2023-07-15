package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"testing"
	"time"
)

var (
	Chan chan []byte
)

func TestServer_OnOpen(t *testing.T) {
	c, err := net.Dial("udp", ":8899")
	if err != nil {
		t.Log(err)
	}
	Chan = make(chan []byte, 10000)
	logger := initLogger()
	go func() {
		for {
			logger.Sugar().Infof("haha")
			time.Sleep(time.Millisecond * 1000)
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
		EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),

		CallerKey:        "caller",
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
		FunctionKey:      "func",
		NameKey:          "N",
	}

	consoleEncode := zapcore.NewConsoleEncoder(encoder)
	//consoleEncode.AddString("server", "TESTSERVER")
	consoleEncode.AddString("server", "mobile")
	w := zapcore.AddSync(new(LoggerAgent))
	core := zapcore.NewCore(consoleEncode, w, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}

type LoggerAgent struct {
	Prefix string
}

func (l *LoggerAgent) Write(p []byte) (n int, err error) {
	select {
	case Chan <- p:
	default:
		fmt.Printf(string(p), "asdads")
	}
	return len(p), nil
}
