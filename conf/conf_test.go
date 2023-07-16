package conf

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGen(t *testing.T) {
	c := GlobalConf{
		Port:        "8899",
		LogFilePath: "/Users/joker/code/go/src/cloud/logagent/",
		MaxSize:     1,
		MaxAge:      7,
		MaxBackups:  3,
		Compress:    true,
	}
	raw, err := json.Marshal(c)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Log(ioutil.WriteFile("conf.json", raw, 0644))
}
