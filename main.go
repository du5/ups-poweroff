package main

import (
	"fmt"
	"os/exec"
	"strconv"

	"context"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-ping/ping"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type CPing struct {
	Host    string `toml:"host"`
	Retry   int    `toml:"retry"`
	Timeout string `toml:"timeout"`
}

type Config struct {
	Hostname string `toml:"hostname"`
	System   string `toml:"system"`
	Redis    string `toml:"redis"`
	Passwd   string `toml:"passwd"`
	Port     int    `toml:"port"`
	Ping     CPing  `toml:"ping"`
	Max      int    `toml:"max"`
}

func (c Config) getCommand() (string, []string) {
	if c.System == "macos" {
		return "osascript", []string{"-e", `tell app "System Events" to shut down`}
	} else if c.System == "windows" {
		return "shutdown", []string{"-s", "-t", "0"}
	} else if c.System == "linux" {
		return "poweroff", []string{}
	} else {
		return "exit", []string{}
	}
}
func (c Config) setCountPP(rdb *redis.Client, count int) {
	_ = rdb.Set(ctx, fmt.Sprintf("%s.loss", c.Hostname), strconv.Itoa(count), 0)
}

func (c Config) getCountPP(rdb *redis.Client) int {
	val, err := rdb.Get(ctx, fmt.Sprintf("%s.loss", c.Hostname)).Result()
	if err == redis.Nil {
		val = "0"
	}
	int, _ := strconv.Atoi(val)
	return int
}

func main() {
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis, conf.Port),
		Password: conf.Passwd,
		DB:       0,
	})

	pinger, _ := ping.NewPinger(conf.Ping.Host)
	pinger.Count = conf.Ping.Retry
	timeOut, _ := time.ParseDuration(conf.Ping.Timeout)
	pinger.Timeout = timeOut
	err := pinger.Run()
	var stats *ping.Statistics
	if err == nil {
		stats = pinger.Statistics()
	}
	isLoss := false
	println(stats.PacketsRecv)
	println(stats.PacketsSent)
	if err != nil || stats.PacketsSent > 0 {
		if err != nil || stats.PacketsRecv == 0 {
			conf.setCountPP(rdb, conf.getCountPP(rdb)+1)
			isLoss = true
		}
	}
	if !isLoss {
		conf.setCountPP(rdb, 0)
	}
	if conf.getCountPP(rdb) > conf.Max {
		stra, strb := conf.getCommand()
		cmd := exec.Command(stra, strb...)
		err = cmd.Run()
	}
}
