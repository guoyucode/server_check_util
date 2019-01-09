package code

import (
	"fmt"
	"testing"
)

func Test_read(t *testing.T) {
	myConfig := new(Config)
	myConfig.InitConfig("./config_test.text")
	fmt.Println(myConfig.Read("default.path"))
	fmt.Println(myConfig.Read("v1"))
	fmt.Println(myConfig.Read("default.v2"))
	fmt.Printf("%v", myConfig.Mymap)
}

func TestConfig_ReadVar(t *testing.T) {
	type args struct {
		key string
		v   *string
	}
	tests := []struct {
		name string
		c    *Config
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ReadVar(tt.args.key, tt.args.v)
		})
	}
}
