package logger

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/daodao97/egin/utils/config"
)

func NewLogger(channel string) Logger {
	return &logger{
		channel: channel,
		logrus:  loggerFactory(),
	}
}

type Logger interface {
	Info(message interface{}, content ...interface{})
	Error(message interface{}, content ...interface{})
}

type logger struct {
	channel string
	logrus  *logrus.Logger
}

func (l logger) Info(message interface{}, content ...interface{}) {
	l.logrus.WithFields(logrus.Fields{
		"content": content,
		"channel": l.channel,
	}).Info(message)
}

func (l logger) Error(message interface{}, content ...interface{}) {
	l.logrus.WithFields(logrus.Fields{
		"content": content,
		"channel": l.channel,
	}).Error(message)
}

func loggerFactory() *logrus.Logger {
	conf := config.Config.Logger
	switch conf.Type {
	case "stdout":
		return stdoutLogger()
	case "file":
		return fileLogger()
	default:
		return stdoutLogger()
	}
}

func stdoutLogger() *logrus.Logger {
	loggerConf := config.Config.Logger

	logger := logrus.New()

	logWriter := io.Writer(os.Stdout)
	logger.Out = logWriter

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.Level(loggerConf.Level))

	return logger
}

func fileLogger() *logrus.Logger {
	loggerConf := config.Config.Logger

	fileName := loggerConf.FileName

	logger := logrus.New()

	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	logger.Out = logWriter

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.Level(loggerConf.Level))

	return logger
}

// TODO
func esLogger() {}

func mongoLogger() {}
