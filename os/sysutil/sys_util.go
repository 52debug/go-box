package sysutil

import (
	"hash/fnv"
	"net"
	"strconv"
	"strings"
)

const (
	basePort    = 20000
	portRange   = 45536 // 20000~65535
	defaultPort = 29876
)

// ExistsInstance 是否已经存在实例
func ExistsInstance(port int) bool {
	if port < 1 || port > 65535 {
		port = defaultPort // 默认端口
	}

	addr := "127.0.0.1:" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}

	// 成功绑定 → 我们是第一个实例
	// 保持 listener 不关闭，直到程序退出
	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close() // 不需要处理连接
		}
	}()

	return false
}

// ExistsInstanceByName 根据应用名称生成端口，减少冲突可能
func ExistsInstanceByName(appName string) bool {
	if appName == "" {
		return ExistsInstance(defaultPort)
	}

	// 小写 + 去掉空白
	appName = strings.ToLower(strings.TrimSpace(appName))

	// 使用 fnv hash（简单、质量较好、不易溢出）
	h := fnv.New32a()
	_, _ = h.Write([]byte(appName))
	hashValue := int(h.Sum32())

	// 映射到端口范围
	port := basePort + (hashValue % portRange)
	return ExistsInstance(port)
}
