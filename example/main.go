package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

var logtText = `asdadasd`

func main() {

	Test()
}
func Test() {
	for i := 0; i < 1; i++ {
		agent := ZapLoggerAgent{}
		logger := agent.Init(&LogAgentConf{
			ServerName: fmt.Sprintf("server%d", i),
			AgentAddr:  "logagent:8899",
		}).Conn().Demons().Logger()
		fmt.Println(i)
		go func(l *zap.Logger) {
			for {
				l.Sugar().Debug(logtText)
				l.Sugar().Error(logtText)
				l.Sugar().Info(logtText)
				l.Sugar().Warn(logtText)
				//logger.Sugar().Panic(logtText)
				time.Sleep(time.Millisecond * 1000)
			}
		}(logger)
	}
	select {}
}
