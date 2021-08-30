package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func os_x_ups_info() (string, []string) {
	if runtime.GOOS != `darwin` {
		log.Panicln(errors.New("此服务只能在 macOS 上运行。"))
	}
	return "pmset", []string{"-g", "batt"}
}

func ups_or_ac(i string) string {
	myregexp := regexp.MustCompile(`^Now\sdrawing\sfrom\s'([\s\S]*)'`)
	params := myregexp.FindStringSubmatch(i)
	return params[1]
}

// ups 品牌名
func ups_name(i string) string {
	myregexp := regexp.MustCompile(`\s-([\s\S]*)\s\(`)
	params := myregexp.FindStringSubmatch(i)
	return params[1]
}

// ups 电量状态
//
// int64: 剩余电量
//
// string:
//     charging 充电
//     discharging 放电
func ups_status(i string) (int64, string, string) {
	myregexp := regexp.MustCompile(`[\s]*([0-9]*)\%;\s([A-Z?a-z]*)[\s?;][\s?\S]`)
	params := myregexp.FindStringSubmatch(i)
	var last_time string
	ps, _ := strconv.ParseInt(params[1], 10, 64)
	if params[2] == "discharging" {
		myregexp = regexp.MustCompile(`discharging;\s([0-9]*:[0-9]*)`)
		last_time = myregexp.FindStringSubmatch(i)[1]
	}
	return ps, params[2], last_time
}
func start_service() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		out, _ := run_command(os_x_ups_info())
		pw, st, lt := ups_status(out)
		c.JSON(200, gin.H{
			"ok": true,
			"status": gin.H{
				"ps_from":   ups_or_ac(out),
				"ups_name":  ups_name(out),
				"power":     pw,
				"status":    st,
				"remaining": lt,
			},
		})
	})
	log.Println(fmt.Sprintf("run at http://%s/ping .", viper.GetString("listen")))
	_ = r.Run(viper.GetString("listen"))
}
