package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

type Params struct {
	File_name string
}

func Init_log(params Params) *logrus.Logger {
	log := logrus.New()
	//日志写入文件
	logrus.SetLevel(logrus.WarnLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	if params.File_name != "" {
		file, err := os.OpenFile(fmt.Sprintf("log/%s.log", params.File_name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Error("Failed to log to File_name, using default stderr")
		}
		log.SetOutput(log.Writer())
	}
	return log
}
