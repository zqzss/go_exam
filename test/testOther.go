package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type student struct {
	name string
	age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handle(w http.ResponseWriter,r *http.Request){
	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	for  {
		messageType,p,err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		p = append([]byte("收到消息: "),p...)
		err = conn.WriteMessage(messageType,p)
		if err != nil {
			log.Fatal(err)
		}

	}

}
func main() {
	//url := "https://orbits.seewintech.com/metro/login"
	//contentType := "application/json"
	//data := `{"username": "xm_yunwei", "password": "admin123"}`
	//resp,err := http.Post(url,contentType,strings.NewReader(data))
	//if err!= nil{
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//b,err := io.ReadAll(resp.Body)
	//fmt.Println(string(b))


	http.HandleFunc("/ws",handle)
	http.ListenAndServe(":8080",nil)
}
