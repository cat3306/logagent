package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type GlobalConf struct {
	Port        string `json:"port"`
	LogFilePath string `json:"log_file_path"`
}

var AppConf *GlobalConf

func Init() {
	var file string
	flag.StringVar(&file, "c", "", "use -c to bind conf file")
	flag.Parse()
	appConf := new(GlobalConf)
	err := LoadJsonConfigLocal(file, appConf)
	if err != nil {
		panic(err)
	}
	AppConf = appConf
}
func LoadJsonConfigLocal(file string, v interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
