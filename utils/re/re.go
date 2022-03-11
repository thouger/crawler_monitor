package re

import (
	"regexp"
)

func FindGroup(str string, reg string) map[string]string {

	// 使用命名分组，显得更清晰
	re := regexp.MustCompile(reg)
	match := re.FindStringSubmatch(str)
	groupNames := re.SubexpNames()

	if match == nil {
		return nil
	}

	//fmt.Printf("%v, %v, %d, %d\n", match, groupNames, len(match), len(groupNames))

	result := make(map[string]string)
	// 转换为map
	for i, name := range groupNames {
		if i != 0 && name != "" { // 第一个分组为空（也就是整个匹配）
			result[name] = match[i]
		}
	}
	return result
	//if _, ok := result["name"]; ok {
	//	fmt.Println(result["name"])
	//}
}
