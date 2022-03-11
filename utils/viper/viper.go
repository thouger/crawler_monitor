package viperr

//package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func Print(Settings map[string]interface{}) {
	for k := range Settings {
		v := Settings[k]
		fmt.Println(v)
		if rec, ok := v.(map[string]interface{}); ok {
			for key, val := range rec {
				log.Printf(" [========>] %s = %s", key, val)
			}
		} else if rec, ok := v.(map[int]interface{}); ok {
			for key, val := range rec {
				log.Printf(" [========>] %s = %s", key, val)
			}
		} else {
			fmt.Printf("record not a map[string]interface{}: %runtime_viper\n", v)
		}
	}
}
func Reader() (*viper.Viper, map[string]interface{}) {
	Settings := runtime_viper.AllSettings()

	//Print(Settings)

	return runtime_viper, Settings
}

func Writer() {

	//runtime_viper.Set("name", "爬虫任务表")

	runtime_viper.Set("叮*买*.id", "1")
	runtime_viper.Set("叮*买*.cron", "0 5 * * *")
	runtime_viper.Set("叮*买*.cmd", "cd ./spider/spiders&&python3 dingdong_app.py")
	runtime_viper.Set("叮*买*.name", "叮*买*")
	runtime_viper.Set("叮*买*.status", "start")

	if err := runtime_viper.WriteConfig(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("写完")
}

var runtime_viper = viper.New()

func init() {

	runtime_viper.SetConfigName("./utils/config/spider_tasks.toml")
	runtime_viper.SetConfigType("toml")
	runtime_viper.AddConfigPath(".")

	if err := runtime_viper.ReadInConfig(); err != nil {
		fmt.Printf("读取爬虫任务文件失败: %s", err)
		return
	}
}

func main() {
	Writer()
}
