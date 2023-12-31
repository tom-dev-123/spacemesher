package main

import (
	"log"
	"net"
	"os/exec"
)


func CheckBin(path string){
	_, err := exec.LookPath(path)
	if err != nil {
		log.Fatal(err)
	}
}


func GetLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}