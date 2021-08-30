package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
)

func run_command(stra string, strb []string) (string, error) {
	cmd := exec.Command(stra, strb...)
	o, err := cmd.CombinedOutput()
	var out string
	if err == nil {
		out = string(o)
	}
	return out, err
}

func shutdown_command() (string, []string) {
	switch runtime.GOOS {
	case "macos":
		return "osascript", []string{"-e", `tell app "System Events" to shut down`}
	case "windows":
		return "shutdown", []string{"-s", "-t", "0"}
	case "linux":
		return "poweroff", []string{}
	}
	return "exit", []string{}
}

func get_with_struct(url string, v interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	bytes := []byte(body)
	return json.Unmarshal(bytes, &v)
}
