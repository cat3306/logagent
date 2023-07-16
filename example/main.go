package main

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

var logtText = `{
    "node_fake_ip": "10.0.0.1",
    "router_fake_ip": "10.0.0.2",
    "port": 5000,
    "mac_addr": "qweqeqwzdad"
}`

func main() {

	Test()
}
func Test() {
	for i := 0; i < 100; i++ {
		agent := ZapLoggerAgent{}
		logger := agent.Init(&LogAgentConf{
			ServerName: fmt.Sprintf("server%d", i),
			AgentAddr:  "192.168.1.7:8899",
		}).Conn().Demons().Logger()
		go func(l *zap.Logger) {
			for {
				l.Sugar().Debug(logtText)
				l.Sugar().Errorf(logtText)
				l.Sugar().Info(logtText)
				l.Sugar().Warn(logtText)
				//logger.Sugar().Panic(logtText)
				time.Sleep(time.Millisecond * 10)
			}
		}(logger)
	}
	select {}
}
