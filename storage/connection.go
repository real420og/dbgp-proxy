package storage

import "net"

type IdeConnection struct {
	IdeHost string `json:"ideHost"`
	RemoteAddr string `json:"remoteAddr"`
	LocalAddr    string `json:"localAddr"`
	IdeKey     string `json:"ideKey"`
	Port     string `json:"port"`
}

func (thus *IdeConnection) FullAddress() string {
	return net.JoinHostPort(thus.IdeHost, thus.Port)
}
