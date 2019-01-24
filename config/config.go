package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/golang/glog"
)

type DBConfig struct {
	DBAddr   string `json:DBAddr`
	DBPort   string `json:DBPort`
	DBUser   string `json:DBUser`
	DBPasswd string `json:DBPasswd`
	Database string `json:database`
}
type TasksConfig struct {
	AbnormalGas bool
}
type Config struct {
	Debug             bool   `json:"debug"`
	HttpServiceAddr   string `json:httpServiceAddr`
	ConcurrenceNum    int    `json:concurrenceNum`
	DefaultConnSize   int32  `json:defaultConnSize`
	MaxConnSize       int32  `json:maxConnSize`
	MinIdleConnSize   int32  `json:minIdleConnSize`
	NewComerThreshold int    `json:newComerThreshold`
	//matrix
	EngineSwitch bool

	//weedfs
	StorageAddr    string
	RankerSwitch   bool
	KafkaAddress   []string
	KafkaTopics    []string
	KafkaGroupID   string
	BatchSize      int
	BatchIntv      int
	LoginTimeOut   int64
	TimeZone       int
	TaskConfig     string
	ConsumerTag    string
	WorkID         int64
	DbBi           *DBConfig
	DbData         *DBConfig
	DbDeep         *DBConfig
	RedisAddr      string
	AuthTimeout    string
	AuthMaxRefresh string
	Tasks          *TasksConfig
}

var (
	ThorConfig *Config
	once       sync.Once
)

func GetConfig() *Config {
	if ThorConfig != nil {
		return ThorConfig
	}
	once.Do(func() {
		ThorConfig = parseConfig()
	})
	return ThorConfig
}

func ParseConfig(config_file *string) {
	configContent, err := ioutil.ReadFile(*config_file)
	glog.Infoln(string(configContent))
	if err != nil {
		glog.Fatalln(err)
	}
	if configContent == nil {
		glog.Fatalln("Error: empty config file")
	}

	config := new(Config)
	err = json.Unmarshal(configContent, config)
	if err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln(config)
	ThorConfig = config
	return
}
func parseConfig() *Config {
	//flag.Parse()
	config_file := flag.String("f", "config.json", "config file root")
	fmt.Print(*config_file)
	flag.Parse()
	configContent, err := ioutil.ReadFile(*config_file)
	if err != nil {
		glog.Fatalln(err)
	}
	if configContent == nil {
		glog.Fatalln("Error: empty config file")
	}

	config := new(Config)
	err = json.Unmarshal(configContent, config)
	if err != nil {
		glog.Fatalln(err)
	}
	body, _ := json.MarshalIndent(config, "", "    ")
	glog.Infoln(string(body))
	return config
}
