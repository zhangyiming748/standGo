package command

import (
	"bufio"
	"io"
	"os/exec"
)

// ExecCommandWithOutput 执行命令并分别处理标准输出和标准错误
func ExecCommandWithOutput(cmd *exec.Cmd) error {
	// 获取标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// 获取标准错误管道
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return err
	}

	// 处理标准输出
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					return
				}
				break
			}
			// 这里可以根据需要处理输出，比如打印到控制台或写入日志
			print(line)
		}
	}()

	// 处理标准错误
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					return
				}
				break
			}
			// 这里可以根据需要处理错误输出
			print("Error: " + line)
		}
	}()

	// 等待命令执行完成
	return cmd.Wait()
}
