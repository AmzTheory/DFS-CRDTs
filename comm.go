package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"

	set "github.com/emirpasic/gods/sets/linkedhashset"
)

/*
set fixed set of clients
unregistering and registering

incopperate with DFS
test it
*/

type ClientManager struct {
	id         int //client ID
	active     *set.Set
	offline    *set.Set
	broadcast  chan RemoteMsg //remote message
	register   chan *Client   //client
	unregister chan *Client   //client
}

type Client struct {
	id     int
	socket net.Conn
	data   chan []byte
}

//ADD & remove operations of the OR_SET
type RemoteMsg struct {
	SenderID int
	Msg      string
	Op       string
	Params   []interface{} //operation operand for the operation
	//uuid
}

func newClientManager(id int) *ClientManager {
	//register used types in gob for the encoding
	gob.Register(replicationElement{})
	gob.Register(RemoteMsg{})
	gob.Register([]interface{}{})

	fmt.Println("Starting server for " + strconv.Itoa(id))
	listener, error := net.Listen("tcp", ":"+strconv.Itoa(id))

	if error != nil {
		fmt.Println(error)
	}

	manager := ClientManager{
		id:         id,
		active:     set.New(), //empty set (will be added when connected)
		offline:    set.New(),
		broadcast:  make(chan RemoteMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()
	go manager.waitForConns(listener)

	return &manager
}

//wait for new clients
func (manager *ClientManager) waitForConns(listener net.Listener) {
	// 	fmt.Printf("Listener on %d\n",manager.id)

	for true {
		connection, _ := listener.Accept()
		b:=make([]byte,20)
		connection.Read(b)
		n := bytes.IndexByte(b,0)
		id,_:= strconv.Atoi( string(b[:n]))

		client := &Client{id:id, socket: connection, data: make(chan []byte)}
		
		manager.register <- client
		//routines for aparticular client
		// go manager.receive(client)
		go manager.send(client)
		//send unsent operations
	}
}

func (manager *ClientManager) start() {
	for {
		// log.Printf("Still listening %d\n",manager.id)
		select {

		case connection := <-manager.register:
			manager.setClientOnline(connection)
		case connection := <-manager.unregister:
			manager.setClientOffline(connection)
			// log.Println("A connection has terminated!")
		case message := <-manager.broadcast:
			// log.Println("broadcast triggered")
			msg := encodeRemoteMsg(message)
			
			for _, conn := range manager.active.Values() {
				con := conn.(*Client)
				select {
				case con.data <- msg:
				default:
					manager.setClientOffline(con)
				}

			}
		}
	}
}
func (manager *ClientManager) setClientOffline(con *Client) {
	close(con.data)
	manager.offline.Add(con.id)
	
	manager.active.Remove(con)
}
func (manager *ClientManager) setClientOnline(con *Client) {
	manager.active.Add(con)
	manager.offline.Remove(con.id)
}

//recieve remote operation from client manager  (need decode the bytes)
func (client *Client) receive(dfs *Dfs) {
	///keep waiting for remote operations
	for {
		message := make([]byte, 4096)
		_, err := client.socket.Read(message)

		if err != nil {
			client.socket.Close()
			break
		}

		//decode the bytes into RemoeteMessage
		rmsg := decodeRemoteMsg(message)
		// fmt.Println("call by ", dfs)
		// dfs.sendRemoteToRep(rmsg.(RemoteMsg))
		dfs.rem<-rmsg.(RemoteMsg)
	}
}

// client manager will send stream of bytes to a client
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()

	//client will send remote operation as local
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			// log.Println(manager.id," ---------------> ",client.id)
			_,_=client.socket.Write(message)
			//record an unsent operatoins
		}
	}
}

func newClient(id int) *ClientManager {
	// fmt.Println("New client...")
	//attempt to connect to all fixed number of clients
	manager := newClientManager(id)
	return manager
}
func (manager *ClientManager) connectToClients(dfs *Dfs) {
	for _, i := range dfs.clients {
		client:=connectToLocalHost(i,manager.id)
		go client.receive(dfs)
	}
}

//encoding/decoding functions


func connectToLocalHost(port int,myID int) *Client{
	connection, error := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if error != nil {
		fmt.Println(error)
	}
		
	var arr [20]byte
	copy(arr[:],strconv.Itoa(myID))
	connection.Write(arr[:])
	client := &Client{id: myID, socket: connection, data: make(chan []byte)}
	
	return client
}

func (client *Client) recieveInterface(ch chan string){
	for {
		message := make([]byte, 4096)
		_, err := client.socket.Read(message)

		if err != nil {
			client.socket.Close()
			break
		}

		//decode
		msg:=bytestoString(message)
		log.Println(msg)
		ch<-msg
	}	
}
func getClientUIServer(port int,myID int,ch chan string) *Client{
	client:=connectToLocalHost(port,myID)
	go client.recieveInterface(ch)
	return client
}


func bytestoString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	return string(b[:n])
}


func encodeRemoteMsg(rmsg interface{}) []byte {
	var ref bytes.Buffer
	enc := gob.NewEncoder(&ref)

	err := enc.Encode(rmsg)
	logEncDecError(err, "encode remote")
	return ref.Bytes()
}

func decodeRemoteMsg(bits []byte) interface{} {
	var msg RemoteMsg
	buf := bytes.NewBuffer(bits)
	err := gob.NewDecoder(buf).Decode(&msg)
	logEncDecError(err, "decode remote")
	return msg
}

func logEncDecError(err error, str string) {
	if err != nil {
		log.Fatal(str+" error:", err)
	}
}
