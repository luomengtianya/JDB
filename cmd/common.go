// cmd通用的一些操作
package cmd

import (
	"fmt"
	"jdb/app"
	"jdb/base"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var ComCmd = &Common{}

type Common struct {
}

// 添加链接数据库的配置
func (common *Common) addClient(flag *pflag.FlagSet) {
	flag.StringP("host", "", "127.0.0.1", "数据库host")
	flag.StringP("port", "", "3306", "数据库端口")
	flag.StringP("u", "u", "root", "用户名")
	flag.StringP("p", "p", "", "密码")
	flag.StringP("o", "o", "", "输出的文件路径")
	flag.StringP("f", "f", "", "HOST,PORT,USERNAME,PASS,DBNAME的properties数据/暂未实现")
}

// 验证数据库相关参数
func (common *Common) checkAndInit(cmd *cobra.Command, args []string) {
	out, _ := cmd.Flags().GetString("o")
	app.Instance().JGetConf.Out = out

	common.checkDataBaseAndInit(cmd, args)
}

// 验证数据库相关参数
func (common *Common) checkDataBaseAndInit(cmd *cobra.Command, args []string) {

	// 默认数据
	conf := base.MySQLConf{
		Host: "127.0.0.1",
		Port: "3306",
		DB:   "information_schema",
		User: "root",
	}

	var host, port, userName, passWord string

	// todo 解析文件配置先忽略
	file, _ := cmd.Flags().GetString("f")
	if file != "" {

	}

	// 优先使用命令行参数
	if v, _ := cmd.Flags().GetString("host"); port != "" {
		host = v
	}
	if v, _ := cmd.Flags().GetString("port"); port != "" {
		port = v
	}
	if v, _ := cmd.Flags().GetString("u"); port != "" {
		userName = v
	}
	if v, _ := cmd.Flags().GetString("p"); port != "" {
		passWord = v
	}

	fmt.Println(fmt.Sprintf("获取参数 == host:%s, port: %s, scheme: %+v, userName: %s, passWord: %s, file: %s",
		host, port, args, userName, passWord, file))

	// todo 参数校验以后再做，先默认参数正常

	if host != "" {
		conf.Host = host
	}

	if port != "" {
		conf.Port = port
	}

	if userName != "" {
		conf.User = userName
	}

	if passWord != "" {
		conf.Password = passWord
		fmt.Println("直接输入密码有危险，请使用stdin的方式输入[初始命令不输入密码]")
	} else {
		fmt.Println("请输入密码:")
		pass := common.stdInPassword()

		conf.Password = pass
	}

	app.Instance().MySQL = conf

	if len(args) == 0 {
		fmt.Println("必须要填写需要操作的库名")
		os.Exit(-1)
	}

	app.Instance().JGetConf.Scheme = args

	if err := app.Instance().InitGormDB(); err != nil {
		fmt.Println(fmt.Sprintf("连接数据库异常: %s", err.Error()))
		//退出当前go程
		// runtime.Goexit()
		os.Exit(-1)
	}
}

// stdInPassword 命令行输入密码
func (common Common) stdInPassword() string {
	tmp, _ := terminal.ReadPassword(int(os.Stdin.Fd())) // os.Stdin.Fd() 或者 syscall.Stdin -- idea不支持展示

	return string(tmp)
}
