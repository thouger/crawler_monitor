package file

import (
	"encoding/csv"
	"fmt"
	"os"
)

func Add(file_name string, data [][]string) {
	f, err := os.OpenFile(file_name, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

func Write(file_name string, data [][]string) {
	f, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// var data = make([][]string, 4)
	// data[0] = []string{"标题", "作者", "时间"}
	// data[1] = []string{"羊皮卷", "鲁迅", "2008"}
	// data[2] = []string{"易筋经", "唐生", "665"}

	f.WriteString("\xEF\xBB\xBF") // 写入一个UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

func Read(file_name string) [][]string {
	f, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	defer f.Close()
	w := csv.NewReader(f)
	data, err := w.ReadAll()
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	// fmt.Println(data)
	return data
}
