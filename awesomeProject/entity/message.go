package entity

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/fatih/set"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FromId   int64
	TargetId int64
	Type     int
	Media    int
	Content  string
	Pic      string
	Url      string
	Desc     string
	Amount   int
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node)

// 读写锁
var rwLocker sync.RWMutex

// 需要：发送者ID，接收者ID，消息类型(文字、音频、图片等)，消息的内容(说了什么话、发了什么图片)，发送类型
// (私聊、群聊、广播)
func Chat(writer http.ResponseWriter, request *http.Request) {
	//校验token 待补充
	query := request.URL.Query()
	Id, _ := strconv.Atoi(query.Get("fromId"))
	fromId := int64(Id)
	//targetId := query.Get("targetId")
	//token := query.Get("token")
	//content := query.Get("content")
	//messageType := query.Get("type")
	conn, err := (&websocket.Upgrader{
		//token检验
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}
	//获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//用户关系
	rwLocker.Lock()
	clientMap[fromId] = node
	rwLocker.Unlock()
	//完成发送逻辑
	sendMsg(fromId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			if err := node.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				panic(err)
			}

		}
	}
}

var udpSendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

// udp数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case data := <-udpSendChan:
			if _, err := conn.Write(data); err != nil {
				panic(err)
			}
		}
	}
}

func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		buf := [512]byte{}
		if _, err := conn.Read(buf[0:]); err != nil {
			panic(err)
		}
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		panic(err)
	}
	switch msg.Type {
	case 1: //私聊
		sendMsg(msg.TargetId, data)
	case 2: //群聊
	case 3: //广播
	default:

	}
}
func sendMsg(targetId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
func init() {

}
