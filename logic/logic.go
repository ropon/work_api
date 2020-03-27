package logic

import (
	"encoding/json"
	"fmt"
	"time"
	"workApi/dblayer"
	"workApi/models"
)

func parseToken(token string) (user models.User, err error) {
	var jsonStr string
	jsonStr, err = decrypt(token)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		return
	}
	return
}

func CheckAuth(token string) (resMap map[string]interface{}, err error) {
	//s, _ := encrypt(`{"phoneNum": "134XXXXXXXX", "machineId": "SDDF2323SDFSFSDF"}`)
	//fmt.Println(s)
	var user models.User
	user, err = parseToken(token)
	if err != nil {
		err = fmt.Errorf("解码失败")
		return
	}
	var redisCache interface{}
	redisKey := user.PhoneNum + user.MachineId
	redisCache, err = dblayer.RedisGet(redisKey)
	//不存在redis缓存
	if err != nil || redisCache == nil {
		var resUser models.User
		err = models.GetUser(&resUser, user.PhoneNum, user.MachineId)
		if err != nil {
			err = fmt.Errorf("未授权")
			return
		}
		if resUser.ExpiryTime.Sub(time.Now()).Seconds() <= 0 {
			err = fmt.Errorf("授权已过期")
			return
		}
		resMap = map[string]interface{}{
			"token":    token,
			"uid":      resUser.ID,
			"realName": resUser.RealName,
			"nickName": resUser.NickName,
		}
		//序列化
		resByte, _ := json.Marshal(resMap)
		_ = dblayer.RedisSet(redisKey, string(resByte), 120)
	} else {
		//fmt.Println("redis cache")
		//反序列化
		_ = json.Unmarshal([]byte(redisCache.(string)), &resMap)
		//序列化反序列化后变成float64
		resMap["uid"] = uint(resMap["uid"].(float64))
	}
	return
}

func NewTaskLog(taskLog *models.TaskLog) (err error) {
	if (taskLog.XiuJin == 0 && taskLog.XueBei == 0) || taskLog.Note == "" {
		err = fmt.Errorf("参数不完整")
		return
	}
	err = models.CreateTaskLog(taskLog)
	if err != nil {
		err = fmt.Errorf("新增记录失败")
		return
	}
	return
}
