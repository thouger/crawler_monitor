package monitor

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os/exec"
	"spider/daemon/child"
	"spider/utils/alert"
	viper "spider/utils/viper"
	"syscall"
)

func Git_pull(log *logrus.Logger) {
	cmd := exec.Command("zsh", "-c", "git pull")
	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Start()
	stdout := outbuf.String()
	if err != nil {
		log.Printf("git pull failed:%s", err)
	}
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exit_status := status.ExitStatus()
				content := fmt.Sprintf("同步代码最终失败，Exit code is %d\n", exit_status)
				log.Error(content)
				alert.Alert(content)
				return
			}
		}
	}
	log.Debugf("拉取了一次代码,%s", stdout)
}

// 添加爬虫定时任务
func Tasks(log *logrus.Logger, c *cron.Cron) {

	//添加爬虫定时任务逻辑部分
	run_tasks := map[string]map[string]interface{}{}
	add := func(settings map[string]interface{}) {
		for k := range settings {
			v := settings[k]
			if rec, ok := v.(map[string]interface{}); ok {
				cmd := rec["cmd"].(string)
				_cron := rec["cron"].(string)
				name := rec["name"].(string)
				//回调函数，如果是配置文件被修改，则删除当前任务id，重新添加
				if _, ok := run_tasks[name]; ok {
					c.Remove(run_tasks[name]["id"].(cron.EntryID))
					delete(run_tasks, name)
					log.WithFields(logrus.Fields{
						"爬虫名字是": name,
						"爬虫命令是": cmd,
						"爬虫调度是": _cron,
					}).Info("爬虫任务已经存在但是配置不对了")
				} else {
					log.WithFields(logrus.Fields{
						"爬虫名字是": name,
						"爬虫命令是": cmd,
						"爬虫调度是": _cron,
					}).Info("首次加入爬虫任务")
				}

				id, _ := c.AddFunc(_cron, func() {
					child.Start_process(name, cmd)
				})
				run_tasks[name] = map[string]interface{}{"cron": _cron, "id": id, "cmd": cmd}
			}
		}
	}

	// 遍历所有爬虫任务
	runtime_viper, settings := viper.Reader()
	add(settings)

	//当配置文件改变后，重新加载配置
	runtime_viper.WatchConfig()
	runtime_viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		_, settings = viper.Reader()
		add(settings)
	})
}
