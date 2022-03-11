package cli

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"spider/daemon/child"
	dbs "spider/utils/db"
	viper "spider/utils/viper"
	"time"

	monitor "spider/daemon/monitor"
)

func Run(name_name string, name string) {
	_, settings := viper.Reader()
	for k := range settings {
		v := settings[k]
		if rec, ok := v.(map[string]interface{}); ok {
			if rec[name_name].(string) == name {
				cmd := rec["cmd"].(string)
				dbs.Connect()
				child.Start_process(name, cmd)
				dbs.Close()
			}
		}
		//path, _ := os.Getwd()
		//
		//tasks := file.Read(path + "/utils/config/spider_task")
		//for _, char := range tasks[1:] {
		//	spider_name := char[2]
		//	if spider_name == name {
		//		dbs.Connect()
		//		cmd := char[1]
		//		child.Start_process(spider_name, cmd)
		//		dbs.Close()
		//		return
		//	}
		//}
		//fmt.Printf("没找到%s程序", name)
	}
}

func Monitor(log *logrus.Logger) {
	c := cron.New()

	//启动监控任务并且添加任务队列
	monitor.Tasks(log, c)

	//1.每一小时自动拉取最新的github仓库
	_, err := c.AddFunc("@hourly", func() {
		monitor.Git_pull(log)
	})
	if err != nil {
		log.Fatalf("拉取git代码失败,错误是：%s", err)
	}

	//c.AddFunc("@every 10s", func() {
	//
	//})
	c.Start()
	dbs.Connect()
	for {
		time.Sleep(time.Second)
	}
}

func Test() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	fmt.Println(basepath)
}
