package pkg

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log *logrus.Logger
)

func New(apps string, debug bool) *logrus.Logger {
	l := logrus.New()

	l.Out = &lumberjack.Logger{
		Filename:   UseHomeDir() + "/.ambpi/" + apps + ".log",
		MaxSize:    viper.GetInt("log.maxsize"),
		MaxBackups: viper.GetInt("log.maxbackups"),
		MaxAge:     1500,
		Compress:   true,
		LocalTime:  true}

	if debug {
		l.SetLevel(logrus.DebugLevel)
	}

	return l
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return log.WithFields(fields)
}

func Error(args interface{}) {
	log.Error(args)
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Infoln(args ...interface{}) {
	log.Infoln(args)
}

func Println(args ...interface{}) {
	log.Println(args)
}

func UseHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
