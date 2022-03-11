package child

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"spider/utils/alert"
	dbs "spider/utils/db"
	log "spider/utils/log"
	"spider/utils/re"
	"strconv"
	"strings"
	"syscall"
)

func Start_process(name string, _cmd string) {
	params := log.Params{}
	process_log := log.Init_log(params)
	// 开始执行爬虫
	process_log.Infof("开始执行%s爬虫任务", name)
	cmd := exec.Command("zsh", "-c", "source ./venv/bin/activate&&"+_cmd)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	//_ = stdoutIn
	//_ = stderrIn
	err2 := cmd.Start()
	if err2 != nil {
		process_log.WithFields(logrus.Fields{
			"命令是": cmd,
			"错误是": err2,
		}).Info("运行命令失败")
	}

	// 初始化爬虫运行状态
	sql := fmt.Sprintf("insert into crawler_done (name,status,start_time) values ('%s','%s',NOW())", name, "start")
	lastrowid := dbs.Execute(sql)
	count := 0

	// stdout
	scanner1 := bufio.NewScanner(stdoutIn)
	outStdout := ""
	go func() {
		for scanner1.Scan() {
			outStdout = scanner1.Text()
			result := re.FindGroup(outStdout, "实际上插入的数据 (?P<count>\\d+) 条")
			fmt.Printf(">>>>>>>>>>>%s", result)
			if result != nil {
				if _, ok := result["count"]; ok {
					intVar, _ := strconv.Atoi(result["count"])
					fmt.Printf(">>>>>>>>>>>%s", intVar)
				}
			}
			process_log.WithFields(logrus.Fields{}).Error(outStdout)
		}
	}()
	// stderr,爬虫运行出错会发送倒数第二条stderr
	last_err := ""
	last_two_err := ""
	scanner2 := bufio.NewScanner(stderrIn)
	go func() {
		for scanner2.Scan() {
			last_two_err = last_err + "\n"
			last_err = scanner2.Text()
			last_two_err += last_err
			//process_log.WithFields(logrus.Fields{}).Error(">>>>>>>>>>" + last_err)
		}
	}()

	// 等待爬虫运行完成
	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exit_status := status.ExitStatus()
				process_log.Errorf("Exit code is %d\n", exit_status)
				process_log.Errorf("%s\n", last_two_err)
				sql = fmt.Sprintf("update crawler_done set status='%s',end_time=now(),exit_status='%d' where id=%d", "error", exit_status, lastrowid)
				dbs.Execute(sql)
				content := fmt.Sprintf("%s爬虫任务执行失败,"+
					"失败的原因：\n"+
					strings.Repeat("*", 10)+"\n"+
					"%s\n"+
					strings.Repeat("*", 10)+"\n"+
					"请检查日志", name, last_two_err)
				alert.Alert(content)
				return
			}
		}
	}

	// 更新爬虫运行状态
	sql = fmt.Sprintf("update crawler_done set status='%s',exit_status='%s',end_time=now(),count='%d' where id='%d'", "success", "0", count, lastrowid)
	dbs.Execute(sql)
	process_log.Infof("爬虫任务%s已经成功结束", name)
}
