package code

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

/*

##定义配置文件: config.int

mode=dev

MailFrom=guoyumail@qq.com
MailFromName=郭宇
MailPassword=sfsaf
MailHost=qq.com:22
MailTo=

[dev]
requestURL=https://url
responseContent=200

[test]
MailFrom=guoyumail@qq.com
MailFromName=郭宇
MailPassword=sfsaf
MailHost=qq.com:22
MailTo=fdagfdag;gfdgfd
requestURL=https://url
responseContent=200

[prod]
MailFrom=guoyumail@qq.com
MailFromName=郭宇
MailPassword=sfsaf
MailHost=qq.com:22
MailTo=fdagfdag;gfdgfd
requestURL=https://url
responseContent=200

*/

/*
读取配置文件代码:
c := new(Config)
c.Init("./config.int")
mode := c.Read("mode")
v := c.Read("mode.key")
*/

const middle = "."

//Config 配置
type Config struct {
	sync.Mutex
	Mymap  map[string]string
	strcet string
	mode   string
}

//InitConfig 初始化一个配置文件
func (c *Config) InitConfig(path string) error {

	c.Lock()
	defer c.Unlock()

	c.Mymap = make(map[string]string)
	c.strcet = ""
	c.mode = ""

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		//if len(c.strcet) == 0 {
		//	continue
		//}

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := ""
		if c.strcet != "" {
			key += c.strcet + middle
		}
		key += frist

		c.Mymap[key] = strings.TrimSpace(second)
	}

	//读取模式
	c.mode = c.read("mode")

	return nil
}

func (c *Config) Read(key string) string {
	c.Lock()
	defer c.Unlock()
	return c.read(key)
}

//读取模式, 例如: test.db
func (c *Config) read(key string) string {
	if c.mode != "" {
		key = c.mode + "." + key
	}

	v, found := c.Mymap[key]
	if !found {
		split := strings.Split(key, ".")
		if len(split) > 1 {
			v2, found2 := c.Mymap[split[1]]
			if !found2 {
				return ""
			}
			return v2
		}
		return ""
	}
	return v
}

//ReadVar 读取变量,读取引用类型
func (c *Config) ReadVar(key string, v *string) {
	v1 := c.Read(key)
	*v = v1
}
