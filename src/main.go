package main

import (
	"deploy-tool-go/src/log"
	ssh "deploy-tool-go/src/ssh"
	"deploy-tool-go/src/un_zip_file"
	"deploy-tool-go/src/zip_file"
	"flag"
	"os"

	"go.uber.org/zap"
)

// 命令行参数
var act string
var src string
var dst string
var user string
var host string
var passwd string
var cmd string
var port int

// 日志变量
var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	// 初始化命令行参数
	flag.StringVar(&act, "act", "", "what action to do: zip, unzip, scp, sh")
	flag.StringVar(&src, "src", "", "zip/unzip/scp src")
	flag.StringVar(&dst, "dst", "", "zip/unzip/scp dst")
	flag.StringVar(&user, "user", "", "ssh user")
	flag.StringVar(&host, "host", "", "ssh host")
	flag.StringVar(&passwd, "passwd", "", "ssh passwd")
	flag.StringVar(&cmd, "cmd", "", "when use sh action, what cmd to run")
	flag.IntVar(&port, "port", 22, "ssh port")
	logger = log.Logger
	sugar = log.Sugar
}

func main() {
	defer func() {
		err := logger.Sync() // flushes buffer, if any
		if err != nil {
			sugar.Infow("logger.Sync() fail", err)
		}
	}()

	flag.Parse()
	sugar.Infof(
		"cmd params: act: %s, src: %s, dst: %s, user: %s, host: %s, passwd: %s, port: %d, cmd: %s",
		act, src, dst, user, host, passwd, port, cmd)

	switch act {
	case "zip":
		logger.Info("current act: zip")
		if src == "" && dst == "" {
			logger.Error("src or dst is empty!")
		} else {
			if err := zip_file.Zip(dst, src); err != nil {
				logger.Error("zip error", zap.Error(err))
			}
		}
	case "unzip":
		logger.Info("current act: unzip")
		if src == "" || dst == "" {
			logger.Error("src or dst is empty!")
		} else {
			if err := un_zip_file.UnZip(dst, src); err != nil {
				logger.Error("unzip error", zap.Error(err))
			}
		}
	case "scp":
		logger.Info("current act: scp")
		if src == "" || dst == "" ||
			user == "" || host == "" || passwd == "" {
			logger.Error("cmd params is empty!")
		} else {
			// ssh
			sshClient, err := ssh.GetConnect(
				host, user, "password", passwd, port,
			)
			if err != nil {
				sugar.Errorf("ssh.GetConnect fail: %v", err)
				return
			} else {
				logger.Info("sshClient connection established!")
				defer sshClient.Close()
			}
			// scp
			scpClient, err := ssh.GetScpClient(sshClient)
			if err != nil {
				sugar.Errorf("ssh.GetScpClient fail %v", err)
				return
			} else {
				logger.Info("scpClient connection established!")
				defer scpClient.Close()
			}
			// Open a file
			f, err := os.Open(src)
			if err != nil {
				sugar.Errorf("os.Open fail %v", err)
				return
			} else {
				defer f.Close()
			}

			// scp transit file
			// 创建目录
			/* 			err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
			   			if err != nil {
			   				sugar.Errorf("os.MkdirAll fail %v", err)
			   				return
			   			} */
			err = scpClient.CopyFile(f, dst, "0644")
			if err != nil {
				sugar.Errorf("scpClient.CopyFile fail %v", err)
			} else {
				sugar.Infof("copy success!")
			}

		}
	case "sh":
		logger.Info("current act: sh")
		if cmd == "" ||
			user == "" || host == "" || passwd == "" {
			logger.Error("cmd params is empty!")
		} else {
			// ssh
			sshClient, err := ssh.GetConnect(
				host, user, "password", passwd, port,
			)
			if err != nil {
				sugar.Errorf("ssh.GetConnect fail: %v", err)
				return
			} else {
				logger.Info("sshClient connection established!")
				defer sshClient.Close()
			}
			// ssh session
			session, err := sshClient.NewSession()
			if err != nil {
				sugar.Errorf("create ssh session fail %v", err)
				return
			} else {
				defer session.Close()
			}
			//执行远程命令
			combo, err := session.CombinedOutput(cmd)
			if err != nil {
				sugar.Errorf("execute cmd fail: %v", err)
				return
			}
			sugar.Infof(string(combo))
		}
	default:
		sugar.Infof("%s,No such action", act)
	}

	sugar.Infof("main end")
	os.Exit(0)
}
