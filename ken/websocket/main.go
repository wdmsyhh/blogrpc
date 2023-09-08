package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// strconv: str conversion
// Atoi(s string) (int, error) (表示 ascii to integer)是把字符串转换成整型数的一个函数
// Itoa(i int) string

var (
	counter = 0
)

func main() {
	http.HandleFunc("/websocket", handleWebSocket)
	log.Println("WebSocket服务器启动成功，监听端口：9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("WebSocket服务器启动失败：", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级WebSocket连接失败：", err)
		return
	}

	defer conn.Close()

	for {
		counter++
		message := []byte(strconv.Itoa(counter))
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("发送消息失败：", err)
			break
		}
		// 控制递增速度
		time.Sleep(time.Second)
	}
}
