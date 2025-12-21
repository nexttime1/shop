package free_port

import (
	"net"
	"option_service/global"
)

func GetFreePort() (int, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", global.Config.LocalInfo.Addr+":0")
	if err != nil {
		return 0, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}
