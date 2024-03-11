package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

//消息
type Message struct {
	gorm.Model
	FormId uint
	//发送者
	TargetId uint
	//接受者
	Type string //消息类型 群聊 私聊 广播
	Media int //消息类型 文字 图片 音频
	Content string //消息内容
	Pic string
	Url string
	Desc string
	Amount int //其他数字统计
}
func (table *Message) TableName() string {
	return "message"
}

type Node struct{
	Conn *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系
var clientMap map[int64]*Node = make(map[int64]*Node,0)
//读写锁
var rwLocker = sync.RWMutex{}

// 需要 ：发送者ID ，接受者ID ，消息类型，发送的内容，发送类型
func chat(writer http.ResponseWriter,request * http.Request){
	//1. 获取参数 并 检验 token 等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId,_ := strconv.ParseInt(Id,10,64)
	//msgType := query.Get("type")
	//targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalida := true //checkToke() 待.........
	conn,err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		//token 校验
		return true
	}}).Upgrade(writer,request,nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2.获取conn
	node := &Node{
		Conn: conn,
		DataQueue: make(chan []byte,50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3. 用户关系
	//4. userid 跟 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//5.完成发送逻辑
	go sendProc(node)
	//6.完成接收逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

func sendProc(node *Node){
	for {
		select{
		case data := <- node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage,data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node){
	for  {
		_,data,err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(data)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<< ",data)
	}
}

var udpsendChan chan []byte = make(chan []byte,1024)
func broadMsg(data []byte){
	udpsendChan <- data
}

func init(){
	go udpSendProc()
	go udpReceProc()
}

//完成udp数据发送协程
func udpSendProc(){
	con,err := net.DialUDP("udp",nil,&net.UDPAddr{
		IP:  net.IPv4(192,168,0,255),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <- udpsendChan:
			_,err := con.Write(data)
			if err != nil{
				fmt.Println(err)
				return
			}
		}
	}
}

//完成udp数据接收协程
