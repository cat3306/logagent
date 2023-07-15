package conf

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"testing"
)

func TestGen(t *testing.T) {
	c := GlobalConf{
		Port: "8869",
	}
	raw, err := json.Marshal(c)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Log(ioutil.WriteFile("conf.json", raw, fs.ModePerm))
}
