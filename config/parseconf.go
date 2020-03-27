package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

func parseConf(fileName string, result interface{}) (err error) {
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)
	if t.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		err = fmt.Errorf("result必须是一个指针且是结构体指针")
		return
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		err = fmt.Errorf("打开配置文件%s失败", fileName)
		return
	}
	//splitFlag := checkOs()
	lineSlice := strings.Split(string(data), "\n")
	resMap := map[string]map[string]string{
		"global": {},
	}
	var section string
	for index, line := range lineSlice {
		//去除首尾空格
		line = strings.TrimSpace(line)
		//排查空行或注释行
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			if line[0] != '[' || line[len(line)-1] != ']' || len(section) == 0 {
				err = fmt.Errorf("第%d行有语法错误", index+1)
				return
			}
			resMap[section] = map[string]string{}
		} else {
			//排查不含=的行
			eqIndex := strings.Index(line, "=")
			if eqIndex == -1 {
				err = fmt.Errorf("第%d行有语法错误", index+1)
				return
			}
			//排查没有key的行 比如 =dsfafads
			//Trim 去除收尾指定字符串
			//TrimSpace 去除收尾空格
			if len(section) == 0 {
				section = "global"
			}
			key := strings.ToLower(strings.TrimSpace(line[:eqIndex]))
			//去除首尾双引号"
			val := strings.Trim(strings.TrimSpace(line[eqIndex+1:]), "\"")
			if key == "" {
				err = fmt.Errorf("第%d行有语法错误", index+1)
				return
			}
			resMap[section][key] = val
		}
	}
	//结构体反射 reflect.TypeOf(result) 下的方法
	//Elem()方法获取指针对应的值
	tElem := t.Elem()
	vElem := v.Elem()
	var structName string
	for k2, v2 := range resMap {
		for i := 0; i < tElem.NumField(); i++ {
			// 取到结构体字段信息
			field := tElem.Field(i)
			// 通过字段名取到tag
			tagName := field.Tag.Get("conf")
			//找到节
			if k2 == tagName && field.Type.Kind() == reflect.Struct {
				structName = field.Name
				break
			}
		}
		//通过structName 找到值信息
		sVal := vElem.FieldByName(structName)
		sType := sVal.Type() // 找到类型信息
		for i := 0; i < sType.NumField(); i++ {
			// 取到结构体字段信息
			field2 := sType.Field(i)
			// 通过字段名取到tag
			tagName2 := field2.Tag.Get("conf")
			switch field2.Type.Kind() {
			//字符串
			case reflect.String:
				//再找到值信息设置值
				sVal.Field(i).SetString(v2[tagName2])
			//int64
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				//将任意类型转换为Int64 转换失败就是对应类型零值
				int64Val, _ := strconv.ParseInt(v2[tagName2], 10, 64)
				sVal.Field(i).SetInt(int64Val)
			//float64
			case reflect.Float32, reflect.Float64:
				//将任意类型转换为Float64 转换失败就是对应类型零值
				float64Val, _ := strconv.ParseFloat(v2[tagName2], 64)
				sVal.Field(i).SetFloat(float64Val)
			//bool
			case reflect.Bool:
				boolVal, _ := strconv.ParseBool(v2[tagName2])
				sVal.Field(i).SetBool(boolVal)
			}
		}
	}
	return
}
