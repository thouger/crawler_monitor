package main

// todo: 更新爬虫状态后需要立刻查询数据库判断该爬虫是否已经爬取完毕
// todo: 检测爬虫程序第二天还在不在，及时kill掉
import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"runtime"
	mycli "spider/daemon/cli"
	_log "spider/utils/log"
)

// 定义全局变量
var log = _log.Init_log(_log.Params{})

func init() {

	//改变一下当前目录为上一级目录
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b) + "/../"
	err := os.Chdir(basepath)
	if err != nil {
		fmt.Errorf("进入上一级目录失败：%s", err)
	}

	pwd, _ := os.Getwd()
	log.Debugf("当前目录是：", pwd)

	log.Debugf("守护进程init函数执行完毕")
}

func main() {
	//run("猪*通")
	//test()
// 	mycli.Monitor(log)
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run",
			Subcommands: []*cli.Command{
				{
					Name:  "id",
					Usage: "任务id",
					Action: func(c *cli.Context) error {
						id := c.Args().Get(0)
						log.Debugf("开始执行run函数，id: %s", id)
						mycli.Run("id", id)
						return nil
					},
				},
				{
					Name:  "name",
					Usage: "任务名字",
					Action: func(c *cli.Context) error {
						name := c.Args().Get(0)
						log.Debugf("开始执行run函数，name: %s", name)
						mycli.Run("name", name)
						return nil
					},
				},
			},
		},
		{
			Name:    "monitor",
			Aliases: []string{"m"},
			Usage:   "monitor",
			Action: func(context *cli.Context) error {
				log.Debugf("开始执行monitor函数")
				mycli.Monitor(log)
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test",

			Action: func(context *cli.Context) error {
				log.Debugf("开始执行test函数")
				mycli.Test()
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("cli启动失败:%s", err)
	}
}
