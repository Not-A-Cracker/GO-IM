package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorPurple = "\033[35m"
)

func printHelp() {
	fmt.Println(colorCyan + "===== 命令列表 =====" + colorReset)
	fmt.Println("  直接输入文字        → 发送广播消息")
	fmt.Println("  /online             → 查看在线用户")
	fmt.Println("  /rename <新名字>    → 修改用户名")
	fmt.Println("  /to <用户名> <消息> → 私聊指定用户")
	fmt.Println("  /help               → 显示帮助")
	fmt.Println("  /quit               → 退出")
	fmt.Println(colorCyan + "====================" + colorReset)
}

func main() {
	addr := "127.0.0.1:9999"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(colorRed+"连接服务器失败:", err, colorReset)
		return
	}
	defer conn.Close()

	fmt.Println(colorGreen+"已连接到服务器 "+addr+colorReset, "")
	printHelp()

	// 读取服务器消息
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(colorRed + "\n与服务器断开连接" + colorReset)
				os.Exit(0)
			}
			line = strings.TrimRight(line, "\n\r")

			switch {
			case strings.Contains(line, "Private Message"):
				fmt.Println(colorPurple + line + colorReset)
			case strings.Contains(line, "Error"):
				fmt.Println(colorRed + line + colorReset)
			case strings.Contains(line, "now online!") || strings.Contains(line, "now offline!"):
				fmt.Println(colorYellow + line + colorReset)
			default:
				fmt.Println(line)
			}
		}
	}()

	// 读取终端输入
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch input {
		case "/quit":
			fmt.Println(colorYellow + "已退出" + colorReset)
			return
		case "/help":
			printHelp()
			continue
		}

		_, err := conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println(colorRed+"发送失败:", err, colorReset)
			return
		}
	}
}
