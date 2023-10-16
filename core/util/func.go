package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	COOKIE_KEY = "r5jO-PYHr50S6EU88Rt9v70FiEwxXvAC"

	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Mili = "2006-01-02T15:04:05.999Z07:00"

	MAX_INT64 int64 = 9223372036854775807
	MAX_INT   int   = 4294967295
)

func StrInArray(search string, items *[]string) bool {
	if items == nil {
		return false
	}
	contains := false
	for _, item := range *items {
		if item == search {
			contains = true
			break
		}
	}
	return contains
}

func GetIp() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func IsRunningInContainer() bool {
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//docker：表示进程运行在 Docker 容器中。
		//kubepods：表示进程运行在 Kubernetes 容器中。
		//containerd：表示进程运行在 containerd 容器中。
		//rkt：表示进程运行在 rkt 容器中。
		if strings.Contains(line, "docker") || strings.Contains(line, "kubepods") || strings.Contains(line, "containerd") || strings.Contains(line, "rkt") {
			return true
		}
	}
	return false
}
