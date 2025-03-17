package conf

import (
	"fmt"
	"github.com/kitex-contrib/config-consul/consul"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
	"os"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Env      string
	Kitex    Kitex    `yaml:"kitex"`
	MySQL    MySQL    `yaml:"mysql"`
	Redis    Redis    `yaml:"redis"`
	Registry Registry `yaml:"registry"`
}

type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kitex struct {
	Service       string `yaml:"service"`
	Address       string `yaml:"address"`
	LogLevel      string `yaml:"log_level"`
	LogFileName   string `yaml:"log_file_name"`
	LogMaxSize    int    `yaml:"log_max_size"`
	LogMaxBackups int    `yaml:"log_max_backups"`
	LogMaxAge     int    `yaml:"log_max_age"`
	MetricsPort   string `yaml:"metrics_port"`
}

type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	//prefix := "conf"
	//confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	//content, err := ioutil.ReadFile(confFileRelPath)
	//if err != nil {
	//	panic(err)
	//}
	//conf = new(Config)
	//err = yaml.Unmarshal(content, conf)
	//if err != nil {
	//	klog.Error("parse yaml error - %v", err)
	//	panic(err)
	//}
	//if err := validator.Validate(conf); err != nil {
	//	klog.Error("validate config error - %v", err)
	//	panic(err)
	//}
	//conf.Env = GetEnv()
	//pretty.Printf("%+v\n", conf)
	client, err := consul.NewClient(consul.Options{
		Addr: "localhost:8500",
	})
	if err != nil {
		panic(err)
	}
	//egisterConfigCallback 方法用来注册一个配置文件的回调。
	//当 Consul 中指定的配置文件发生变化时，回调函数就会被执行。

	//onsul.AllocateUniqueID()：这是一个生成唯一 ID 的方法，Consul 通过该 ID 来标识此回调请求。
	//在 Consul 中，多个客户端可以注册回调函数，AllocateUniqueID() 会确保每个回调有一个唯一的标识
	client.RegisterConfigCallback("product/test.yaml", consul.AllocateUniqueID(), func(s string, parser consul.ConfigParser) {
		//s这是一个字符串，表示从 Consul 获取到的配置内容（product/test.yaml 的内容）。
		err = yaml.Unmarshal([]byte(s), &conf)
		if err != nil {
			fmt.Println(err)
		}
		pretty.Printf("%+v\n", conf)
	})

}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}

func LogLevel() klog.Level {
	level := GetConf().Kitex.LogLevel
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
