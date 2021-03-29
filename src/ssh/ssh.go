package ssh

import (
	"deploy-tool-go/src/log"
	"fmt"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// 日志变量
var sugar *zap.SugaredLogger

var config *ssh.ClientConfig

//var host, user string

func init() {
	sugar = log.Sugar
}

// 返回ssh client
func GetConnect(_host string, _user string, _authType string, _passwd string, _port int) (*ssh.Client, error) {

	/* 	host = _host
	   	port = _port
	   	user = _user */

	config = &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            _user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if _authType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(_passwd)}
	} else {
		//config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", _host, _port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		sugar.Infof("create ssh client fail: %v", err)
		return nil, err
	}
	//defer sshClient.Close()

	//创建ssh-session
	/* 	session, err := sshClient.NewSession()
	   	if err != nil {
	   		sugar.Infof("create ssh session fail %v", err)
	   		return nil, err
	   	}
	   	defer session.Close() */
	//执行远程命令
	/* 	combo, err := session.CombinedOutput("whoami; cd /; ls -al")
	   	if err != nil {
	   		sugar.Infof("execute cmd fail: %v", err)
	   		return nil, err
	   	}
	   	fmt.Println(string(combo)) */
	//sugar.Infof("cmd output: %v", string(combo))
	return sshClient, nil
}

// GetScpClient 从sshClient获取scpClient
func GetScpClient(sshClient *ssh.Client) (scp.Client, error) {
	return scp.NewClientBySSH(sshClient)
}
