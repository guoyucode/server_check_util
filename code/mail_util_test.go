package code

import "testing"

func TestSendMail(t *testing.T) {
	type args struct {
		subject string
		body    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "guoyu",
			args: args{
				subject: "邮件主题",
				body:    "邮件主体",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMail(tt.args.subject, tt.args.body)
		})
	}
}
