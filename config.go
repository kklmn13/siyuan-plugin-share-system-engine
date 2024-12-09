package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
	"gopkg.in/yaml.v3"
)

type appConfig struct {
	Basic `yaml:"basic"`
	Web   `yaml:"web"`
	SQL   `yaml:"sql"`
	db    *sql.DB
}
type Basic struct {
	ListenPort    string `yaml:"listen"`
	SavePath      string `yaml:"savePath"`
	ShareBaseLink string `yaml:"shareBaseLink"`
}
type Web struct {
	FileMaxMB int64  `yaml:"fileMaxMB"`
	SSLEnable bool   `yaml:"sslEnable"`
	SSLCERT   string `yaml:"sslCERT"`
	SSLKEY    string `yaml:"sslKEY"`
}
type SQL struct {
	APIEnable   bool   `yaml:"apiEnable"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	SYSFilename string `yaml:"sysFilename"`
}

type flagStruct struct {
	config_file string
	db_file     string
	version     bool
	loglevel    int
}

func (app *appConfig) is_empty() {
	if app.ListenPort == "" {
		app.ListenPort = "0.0.0.0:25934"
		log.Infof("监听端口使用强制参数:%s", app.ListenPort)
	}
	if app.SavePath == "" {
		app.SavePath = "/data"
		log.Infof("存储路径使用强制参数:%s", app.SavePath)
	}
	if app.ShareBaseLink == "" {
		app.ShareBaseLink = "http://127.0.0.1:25934"
		log.Infof("分享地址使用强制参数:%s", app.Basic.ShareBaseLink)
	}
	if app.FileMaxMB < 0 {
		app.FileMaxMB = 100
		log.Infof("最大文件使用强制参数:%d", app.FileMaxMB)
	}
	if app.SQL.SYSFilename == "" {
		app.SQL.SYSFilename = "info.db"
		log.Infof("数据库存储文件使用强制参数:%s", app.SQL.SYSFilename)

	}
}
func init_flag() flagStruct {
	var f flagStruct
	flag.StringVar(&f.config_file, "c", "config.yaml", "打开配置文件")
	flag.StringVar(&f.db_file, "d", "", "清空数据并导入SQL数据文件")
	flag.IntVar(&f.loglevel, "l", log.LEVELINFOINT, fmt.Sprintf("日志等级,%d-%d", log.LEVELFATALINT, log.LEVELDEBUG3INT))
	flag.BoolVar(&f.version, "v", false, "查看版本号")

	flag.Parse()
	return f
}
func init_config(flag flagStruct) {
	if flag.version {
		fmt.Println(version)
		os.Exit(0)
	}
	log.SetLevelInt(flag.loglevel)
	_, g := log.GetLevel()
	fmt.Printf("日志等级:%s\n", g)

	log.Infof("读取配置文件")
	l := []string{flag.config_file, "/data/config.yaml", "/config.yaml"}
	for _, f := range l {
		if !tools.FileExist(f) {
			log.Warnf("配置文件不存在 路径:%s", f)
			continue
		}
		log.Infof("配置文件: %s", flag.config_file)
		yamlFile, err := os.ReadFile(flag.config_file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &app)
		if err != nil {
			panic(err)
		}

		log.Infof("共享链接: %s", app.Basic.ShareBaseLink)
		log.Infof("资源文件保存位置: %s", app.Basic.SavePath)
		return
	}
	log.Info("使用默认配置")

}
func init_env() {
	if v := os.Getenv("SPSS_LISTEN"); v != "" {
		log.Infof("SPSS_LISTEN=%s", v)
		app.ListenPort = v
	}

	if v := os.Getenv("SPSS_SAVE_PATH"); v != "" {
		log.Infof("SPSS_SAVE_PATH=%s", v)
		app.Basic.SavePath = v
	}

	if v := os.Getenv("SPSS_SHARE_LINK"); v != "" {
		log.Infof("SPSS_SHARE_LINK=%s", v)
		app.Basic.ShareBaseLink = v
	}
	if v := os.Getenv("SPSS_WEB_FILE_MAX"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			app.Web.FileMaxMB = 100
			log.Infof("最大文件使用强制参数:%d", app.FileMaxMB)
		} else {

			app.Web.FileMaxMB = int64(i)
		}
		log.Infof("SPSS_WEB_FILE_MAX=%s", v)
	}
	if v := os.Getenv("SPSS_WEB_SSL"); v == "true" || v == "TRUE" {
		app.Web.SSLEnable = true
		if v := os.Getenv("SPSS_WEB_SSL_CERT"); v != "" {
			log.Infof("SPSS_WEB_SSL_CERT=%s", v)
			app.Web.SSLCERT = v
		}
		if v := os.Getenv("SPSS_WEB_SSL_KEY"); v != "" {
			log.Infof("SPSS_WEB_SSL_KEY=%s", v)
			app.Web.SSLKEY = v
		}
	}
	if v := os.Getenv("SPSS_DB_API"); v == "true" || v == "TRUE" {
		log.Infof("SPSS_DB_API=%s", v)
		app.SQL.APIEnable = true
	}

	if v := os.Getenv("SPSS_DB_AUTH"); v != "" {
		if !strings.Contains(v, ":") {
			log.Fatal("SPSS_DB_LOGIN变量格式错误")
		}
		l := strings.Split(v, ":")
		app.SQL.APIEnable = true
		app.SQL.Username = l[0]
		app.SQL.Password = l[1]
		log.Infof("SPSS_DB_AUTH=%s", v)

	}
	if v := os.Getenv("SPSS_DB_SAVE"); v != "" {
		app.SQL.SYSFilename = v
	}
	app.is_empty()
}
