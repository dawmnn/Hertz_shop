package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kr/pretty"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
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
}

type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	//Once.Do 方法用于确保某个函数 f() 只执行一次。
	//错误的实现会导致并发调用时出现不正确的返回顺序（第二个调用可能不等待第一个调用完成）。
	//解决方案是通过“慢路径”（doSlow）来确保正确的同步，使得 f() 完成后才返回。
	return conf
}

func initConf() {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	//这行代码通过 filepath.Join 来连接路径，构建配置文件的相对路径。首先，
	//调用 GetEnv() 获取当前的环境（后面会讲解该函数），然后将其与 "conf" 和 "conf.yaml"
	//拼接成完整的配置文件路径。
	//假设 GetEnv() 返回 "test"，那么 confFileRelPath 就是 "conf/test/conf.yaml"。
	content, err := ioutil.ReadFile(confFileRelPath)
	//使用 ioutil.ReadFile 来读取指定路径 confFileRelPath 中的配置文件内容，
	//并将文件内容存储在 content 变量中。
	//如果读取文件失败，err 将会包含错误信息。
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	//使用 yaml.Unmarshal 将从配置文件读取到的字节数据（content）解析为 Config 结构体的实例。
	//yaml.Unmarshal 是 Go 语言的一个函数，用于将 YAML 格式的数据反序列化到 Go 结构体中。
	if err != nil {
		klog.Error("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		klog.Error("validate config error - %v", err)
		panic(err)
	}
	//在成功解析配置文件之后，使用 validator.Validate 来验证配置的内容是否符合要求。这是一个常见的做法，
	//可以确保配置文件中的数据完整且符合预期。如果验证失败（err != nil），
	//会记录错误日志并通过 panic 抛出异常。
	conf.Env = GetEnv()
	//通过调用 GetEnv() 函数获取当前环境的值（如 "test" 或 "prod"），
	//并将其赋值给配置对象中的 Env 字段。这样，配置对象就包含了当前的运行环境。
	pretty.Printf("%+v\n", conf)
	//使用 pretty.Printf（这是一个用于格式化输出的第三方包）打印 conf 配置对象的内容，
	//格式化为易读的结构化信息。%+v 会打印结构体的字段名和对应的值，输出时会更具可读性。
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	//通过 os.Getenv("GO_ENV") 获取环境变量 GO_ENV 的值。
	//如果该环境变量存在，返回其值；如果不存在，返回一个空字符串。
	if len(e) == 0 {
		return "test"
	}
	//如果 GO_ENV 环境变量没有设置（即长度为 0），则返回默认值 "test"，表示默认环境为 test。
	//如果环境变量存在且有值，则返回其值，通常是 "prod" 或 "dev" 等环境名称。
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
