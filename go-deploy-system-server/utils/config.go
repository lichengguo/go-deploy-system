package utils

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

// 加载ini配置文件

// 声明配置文件参数变量
var (
	AppMode      string // 启动模式 debug 开发模式；release 生产模式
	HttpPort     string // 监听地址和端口
	JwtKey       string // JWT加盐字符串
	PwdKey       string // 登录密码加盐字符串
	ServerGitKey string // 服务器、Git密码模式下的key
	KeyFilePath  string // 秘钥文件存储目录
	CodePath     string // 代码存放目录 data/git/

	DbHost     string // 数据库连接地址
	DbPort     string // 数据库连接端口
	DbUser     string // 数据库连接账号
	DbPassWord string // 数据库连接密码
	DbName     string // 数据库名称

	LogPath      string // 日志文件存储目录
	LogFileName  string // 日志文件名称
	LogSaveTime  int    // 日志最大保存时间
	LogSplitSize int    // 日志切割

	GitCodePath string // 代码存放的绝对路径 .../data/git/
	LockDirPath string // 锁文件目录
	ExecBaseDir string // 执行程序文件所在的绝对目录
	err         error
)

var (
	configFile string // 配置文件
	help       bool   // 帮助
)

// init 只要包被导入就执行
func init() {
	// 当前程序执行的绝对目录
	ExecBaseDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("获取当前程序执行目录失败")
		panic(err)
	}

	// 配置文件
	configPathDefault := fmt.Sprintf("%s/%s", ExecBaseDir, "config/config.ini")       // 默认配置文件位置
	flag.StringVar(&configFile, "config", configPathDefault, "set config file path.") // 自定义配置文件位置
	flag.BoolVar(&help, "help", false, "Show help information.")
	flag.Parse()

	if help {
		fmt.Printf("Usage:\n\t")
		fmt.Printf("-config config.ini\n")
		os.Exit(1)
	}

	file, err := ini.Load(configFile) // 读取配置文件内容
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件.")
		panic(err)
	}

	loadConfig(file)
	initDir()
}

func loadConfig(file *ini.File) {
	// server
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
	PwdKey = file.Section("server").Key("PwdKey").MustString("89js82js72")
	ServerGitKey = file.Section("server").Key("ServerGitKey").MustString("89js82js72")
	KeyFilePath = file.Section("server").Key("KeyFilePath").MustString("upload/key")
	CodePath = file.Section("server").Key("CodePath").MustString("data/git")

	// database
	DbHost = file.Section("database").Key("DbHost").MustString("127.0.0.1")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("root")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("root123")
	DbName = file.Section("database").Key("DbName").MustString("yiihua_go")

	// log
	LogPath = file.Section("log").Key("LogPath").MustString("log")
	LogFileName = file.Section("log").Key("LogFileName").MustString("ops")
	LogSaveTime = file.Section("log").Key("LogSaveTime").MustInt(10)
	LogSplitSize = file.Section("log").Key("LogSplitSize").MustInt(100)
}

// 创建相关目录
func initDir() {
	// 日志文件目录
	logDirPath := fmt.Sprintf("%s/%s", ExecBaseDir, LogPath)
	_, err = os.Stat(logDirPath) // 判断目录是否存在
	if err != nil {
		_ = os.MkdirAll(logDirPath, 0755) // 不存在则创建目录
	}

	// 秘钥文件目录
	keyDirPath := fmt.Sprintf("%s/%s", ExecBaseDir, KeyFilePath)
	_, err = os.Stat(keyDirPath)
	if err != nil {
		_ = os.MkdirAll(keyDirPath, 0755)
	}

	// 代码文件存储的绝对目录
	GitCodePath = fmt.Sprintf("%s/%s", ExecBaseDir, CodePath)
	_, err = os.Stat(GitCodePath)
	if err != nil {
		_ = os.MkdirAll(GitCodePath, 0755)
	}

	// 锁文件目录
	LockDirPath = fmt.Sprintf("%s/%s/%s", ExecBaseDir, CodePath, "lock")
	_, err = os.Stat(LockDirPath)
	if err != nil {
		_ = os.MkdirAll(LockDirPath, 0755)
	}
}
