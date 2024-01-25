package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	port := get_port()
	// exec command lsof with result string
	output := run_lsof(port)
	// get pid by parse lsof result
	pid := get_pid(output)
	// killing process by pid
	kill_pid(pid)
}

func kill_pid(pid string) {
	cmd := exec.Command("kill", "-9", pid)
	if _, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintln("Failed killed process PID:", pid))
	}
	fmt.Println("killed process PID:", pid)
}

func run_lsof(port string) string {
	cmd := exec.Command("/usr/sbin/lsof", "-t", "-i", fmt.Sprint(":", port))
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintln("Not found PID:", port, ".", err)
		panic(msg)
	}
	return string(output)
}

func get_port() string {
	var port *string
	// 定义命令行参数
	port = flag.String("p", "", "socket port number")
	// 解析输入的命令行参数
	flag.Parse()
	if len(*port) == 0 {
		// 没有通过 -p 传入参数，尝试从 cli 参数获取
		if len(os.Args) > 1 {
			return os.Args[1]
		} else {
			panic("Please input port number by argument p")
		}
	}
	return *port
}

func get_pid(pidStr string) string {
	// pid, err := strconv.ParseUint(strings.TrimSpace(pidStr), 10, 32)
	// if err != nil {
	// 	panic(fmt.Sprintln("Error parse string to uint of pid", err))
	// }
	return strings.TrimSpace(pidStr)
}

// 不需要了 lsof -t 直接返回 PID
func parse_port(output string) string {
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	// 使用 bufio.Scanner 逐行读取输出
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过标题行，COMMAND开头的
		if strings.Contains(line, "COMMAND") {
			continue
		}

		// 分割每一列，获取 PID
		fields := strings.Fields(line)
		fmt.Println("fields:", fields)
		if len(fields) > 2 {
			pidStr := fields[1]
			fmt.Println("get parse PID:", pidStr)
			return get_pid(pidStr)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}
	panic("Failed Parse port")
}
