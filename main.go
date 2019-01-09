package main

import (
	"errors"
	"fmt"
	"gitee.com/guoyucode/logs"
	"io/ioutil"
	"mailutil/code"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var config *code.Config

//LoadConfigFile 加载配置文件
func LoadConfigFile()  {
	if e := config.InitConfig("./config.ini"); e != nil {
		logs.Emergency("读取配置文件出错:" + e.Error())
	}
}


func main() {

	logs.WriteToFile()
	logs.Async(100)
	logs.SetLevel(logs.LevelDebug)

	logs.Info("mailutil开始启动")

	config = new(code.Config)
	LoadConfigFile()

	var level, fromAddress, password, subject, host, sendTos, testURL, macthContent, sleepSecond string
	readVal := func (){
		config.ReadVar("level", &level)
		config.ReadVar("fromAddress", &fromAddress)
		config.ReadVar("password", &password)
		config.ReadVar("subject", &subject)
		config.ReadVar("host", &host)
		config.ReadVar("sendTos", &sendTos)
		config.ReadVar("testURL", &testURL)
		config.ReadVar("macthContent", &macthContent)
		config.ReadVar("sleepSecond", &sleepSecond)
	}
	readVal()

	//循环加载配置文件
	go func() {
		for ;true ;  {
			time.Sleep(time.Second * 12)
			LoadConfigFile()
			logs.Debug("重新加载配置文件")
		}
	}()

	for ;true ;  {

		time.Sleep(time.Second * 10)
		readVal()
		if e := validate(level, subject, password, fromAddress, host, sendTos); e != nil {
			logs.Emergency("配置文件中参数有误:" + e.Error())
		}

		var configStr = `{"level":LEVEL,"subject":"SUBJECT","fromAddress":"FROMADDRESS","username":"USERNAME","password":"PASSWORD","host":"HOST","sendTos":[SENDTOS]}`
		configStr = strings.Replace(configStr, "LEVEL", level+"", -1)
		configStr = strings.Replace(configStr, "SUBJECT", subject, -1)
		configStr = strings.Replace(configStr, "FROMADDRESS", fromAddress, -1)
		configStr = strings.Replace(configStr, "USERNAME", fromAddress, -1)
		configStr = strings.Replace(configStr, "PASSWORD", password, -1)
		configStr = strings.Replace(configStr, "HOST", host, -1)
		s := sendTos
		s = "\"" + s + "\""
		s = strings.Replace(s, ",", "\",\"", -1)
		configStr = strings.Replace(configStr, "SENDTOS", s, -1)
		if err := logs.SetLogger("smtp", configStr); err != nil{
			fmt.Println(err.Error())
		}

		logs.Info("测试开始", testURL)
		client := http.Client{
			Timeout: time.Second * 15,
		}
		resp, err := client.Get(testURL)
		if err != nil{
			logs.Alert("请求出错了", err)
			continue
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Alert("请求时读取数据出错了", err)
			continue
		}

		responseStr := string(buf)
		//logs.Debug(responseStr)
		b := strings.Contains(responseStr, macthContent)
		if !b {
			logs.Alert("测试错误", "请求返回值错误", responseStr, macthContent)
			continue
		}else{
			logs.Info("测试正常", testURL)
		}

		i, _ := strconv.Atoi(sleepSecond)
		time.Sleep(time.Duration(i) * time.Second)
	}

}


func validate(level, subject, password, fromAddress, host, sendTos string) error {
	msg := ""
	if level == "" {
		msg += " level参数不能为空"
	}
	if subject == "" {
		msg += " subject参数不能为空"
	}
	if password == "" {
		msg += " password参数不能为空"
	}
	if fromAddress == "" {
		msg += " MailFromName"
	}
	if host == "" {
		msg += " host参数不能为空"
	}
	if sendTos == "" {
		msg += " sendTos参数不能为空"
	}
	if msg != "" {
		return errors.New(msg)
	}

	return nil
}
