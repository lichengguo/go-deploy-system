package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"go-deploy-system-server/utils"
	"math"
	"os"
	"time"
)

func Log() gin.HandlerFunc {
	filePath := fmt.Sprintf("%s/%s/%s", utils.ExecBaseDir, utils.LogPath, utils.LogFileName) // 日志目录文件

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("打开日志文件错误", err)
		panic(err)
	}

	logger := logrus.New()             // 构造函数，构造一个*Logger
	logger.Out = f                     // 将日志输出到文件
	logger.SetLevel(logrus.DebugLevel) // 设置日志记录的级别

	// 日志分割
	// 构造函数
	logWriter, _ := retalog.New(
		filePath, // 参数1 日志文件名称
		retalog.WithMaxAge(time.Duration(utils.LogSaveTime)*24*time.Hour), // 参数2 文件最大保存时间
		retalog.WithRotationSize(int64(utils.LogSplitSize*1024*1024)),     // 参数3 日志按照大小切割 单位:MB
	)
	// 所有级别的日志写入到logWriter
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(Hook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0))) // 函数耗时
		hostName, err := os.Hostname()                                                               // 主机名
		if err != nil {
			hostName = "unknown"
		}
		statusCode := c.Writer.Status()    // 请求状态码
		clientIp := c.ClientIP()           // 客户端IP
		userAgent := c.Request.UserAgent() // 浏览器
		dataSize := c.Writer.Size()        // 请求大小
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   // 请求方法
		path := c.Request.RequestURI // 请求路径

		// logrus的标准写法
		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		// gin框架内部错误
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		// 状态码错误
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
