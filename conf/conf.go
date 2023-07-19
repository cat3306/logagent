package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

type GlobalConf struct {
	Port        string `json:"port"`
	LogFilePath string `json:"log_file_path"`

	MaxSize         int    `json:"max_size"`    // 日志文件的最大大小(MB为单位)
	MaxAge          int    `json:"max_age"`     // 保留旧文件的最大天数
	MaxBackups      int    `json:"max_backups"` // 保留旧文件的最大个数
	Compress        bool   `json:"compress"`    // 是否压缩/归档旧文件
	UseGPool        bool   `json:"use_g_pool"`
	NotificationUrl string `json:"notification_url"`
	NameSpace       string `json:"name_space"`
}

var AppConf *GlobalConf

func Init() {
	var file string
	flag.StringVar(&file, "c", "conf.json", "use -c to bind conf file")
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
