package re

import (
	"fmt"
	"regexp"
)

func main() {
	str := `共导出 3 条数据 到 competitor_product,总共插入的数据 3 条,重复 0 条`

	// 使用命名分组，显得更清晰
	re := regexp.MustCompile(`共导出 3 条数据 到 competitor_product,总共插入的数据 (?P<count>\d+) 条,重复 0 条`)
	match := re.FindStringSubmatch(str)
	groupNames := re.SubexpNames()

	fmt.Printf("%v, %v, %d, %d\n", match, groupNames, len(match), len(groupNames))

	result := make(map[string]string)
	// 转换为map
	for i, name := range groupNames {
		if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
			result[name] = match[i]
		}
	}
	if _, ok := result["count"]; ok {
		fmt.Println("找到" + result["count"])
	}
}
