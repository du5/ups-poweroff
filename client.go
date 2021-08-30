package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type pingStruct struct {
	OK     bool         `json:"ok"`
	Status statusStruct `json:"status"`
}
type statusStruct struct {
	Power     int64  `json:"power"`
	PsFrom    string `json:"ps_from"`
	Remaining string `json:"remaining"`
	Status    string `json:"status"`
	UpsName   string `json:"ups_name"`
}

func start_client() {
	tc := time.NewTicker(viper.GetDuration("tc_time") * time.Second)
	var ping pingStruct
	var time_out int64
	for {
		<-tc.C
		err := get_with_struct(fmt.Sprintf("http://%s/ping", viper.GetString("listen")), &ping)
		log.Println(fmt.Sprintf("UPS %s, 当前状态 %s, 剩余电量 %d%%", ping.Status.UpsName, ping.Status.Status, ping.Status.Power))
		if err != nil {
			if time_out += 1; time_out > viper.GetInt64("time_out") {
				// 计数器 +1
				_, _ = run_command(shutdown_command())
				os.Exit(0)
			}
			continue
		}
		// 重置计数器
		time_out = 0

		if ping.Status.Status == "charging" && ping.Status.PsFrom == "AC Power" {
			// 正在充电
			continue
		}
		// 正在放电
		if ping.Status.Power <= viper.GetInt64("ups_last") {
			_, _ = run_command(shutdown_command())
			os.Exit(0)
		}

	}
}
