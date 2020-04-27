package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"

	set "github.com/emirpasic/gods/sets/linkedhashset"
	// queue "github.com/enriquebris/goconcurrentqueue"
	cmap "github.com/orcaman/concurrent-map"
)


/*
start server
		once c is connected
			set online
			send unsent ops  (map dfs->queue)
			
when sedning brodcast
	if replica offline  (onlineMAP)
		add into unsent op dt(concQueue)
	else
		send
	
socket failure
	failed to send
		set offline
		add the missing ops into Queue

recover lost replica con (client side)
	1-start gr to connect to each replica
		once found , the recive will trigger
			recieve will notice socket fialure
						return to 1



*/
type ClientManager struct {
	id         int //client ID
	active     *set.Set
	offline    *set.Set
	onlineMap  *cmap.ConcurrentMap
	unSentOps  *cmap.ConcurrentMap
	broadcast  chan RemoteMsg //remote message
	register   chan *Client   //client
	unregister chan *Client   //client
	dfs 	   *Dfs
}


func (manager *ClientManager) getOnlineMap(id string) *Client{
	v,_:=manager.onlineMap.Get(id)
	if(v!=nil) {return v.(*Client)}

	return nil//incase the value still nil (not online)
}
func (manager *ClientManager) setOnlineMap(id int,value *Client){
	i:=strconv.Itoa(id)
	manager.onlineMap.Set(i,value)
}
func (manager *ClientManager) getMissingOps(id int) []interface{}{
	i:=strconv.Itoa(id)
	defer manager.unSentOps.Set(i,set.New()) //init new set once returned  
	v,_:=manager.unSentOps.Get(i)
	return v.(*set.Set).Values()
}
func (manager *ClientManager) addMissingOps(id string,op RemoteMsg){
	q,_:=manager.unSentOps.Get(id)
	s:=q.(*set.Set)
	s.Add(op)//added
}


type Client struct {
	id       int
	destPort int
	socket 	 net.Conn
	data   	 chan []byte
}
																																										

type RemoteMsg struct {
	SenderID int
	Msg      string
	Op       string
	P1       interface{}
	P2       interface{}
}


func newClientManager(d *Dfs) *ClientManager {
	//register used types in gob for the encoding
	gob.Register(replicationElement{})
	gob.Register(RemoteMsg{})
	gob.Register([]interface{}{})
	

	// fmt.Println("Starting server for " + strconv.Itoa(id))
	listener, error := net.Listen("tcp", ":"+strconv.Itoa(d.id))

	if error != nil {
		fmt.Println(error)
	}
	//store the initalOnlineMap with nil (yet to be connected)
	onlineMap:=cmap.New()
	unSentOps:=cmap.New()
	for v,_:=range d.clients{
		i:=strconv.Itoa(v)
		onlineMap.Set(i,nil)
		unSentOps.Set(i,set.New())
	}

	manager := ClientManager{
		id:         d.id,
		active:     set.New(), //empty set (will be added when connected)
		offline:    set.New(),
		onlineMap:  &onlineMap,
		unSentOps:  &unSentOps,
		broadcast:  make(chan RemoteMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		dfs:		d,
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
		// fmt.Println("connected with",id)
		manager.register <- client
		go manager.send(client)
		
		//send unSent ops  (locking the queue)
		for _,op := range manager.getMissingOps(id){
			client.data<-encodeMsg(op)
		}
		go manager.receive(client)
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
			
			
			for _, id := range manager.onlineMap.Keys() {
				con:=manager.getOnlineMap(id)
				if(con==nil) { // is offline
					manager.addMissingOps(id,message)
					fmt.Println(id,"op is recorded")
					continue
				}
				msg := encodeMsg(message)
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
	manager.setOnlineMap(con.id,nil)  //set offline
	fmt.Println(con.id,"gone offline")
}
func (manager *ClientManager) setClientOnline(con *Client) {
	manager.setOnlineMap(con.id,con)  //set online
}

//recieve remote operation from client manager  (need decode the bytes)
func (client *Client) receive(dfs *Dfs) {
	///keep waiting for remote operations
	for {
		message := make([]byte, 4096)
		_, err := client.socket.Read(message)

		if err != nil {
			client.socket.Close()
			go tryConnect(client.destPort,client.id,dfs)  //socket failure -- try to receieve again
			break
		}

		//decode the bytes into RemoeteMessage
		rmsg:=decodeMsg(message)
		// fmt.Println("call by ", dfs)
		// dfs.sendRemoteToRep(rmsg.(RemoteMsg))
		dfs.rem<-rmsg.(RemoteMsg)
	}
}
func (manager *ClientManager) receive(client *Client){
	message := make([]byte, 4096)
	_, err := client.socket.Read(message)
	if(err!=nil){
		log.Println(err)
	}

	//indicate soc closed
	manager.unregister<-client
}

// client manager will send stream of bytes to a client
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()

	//client will send remote operation as local
	for {
		select {
		case message, _ := <-client.data:
			_,err:=client.socket.Write(message)
			if err!=nil{
				log.Println(err)
			}
		}
	}
}

func newClient(d *Dfs) *ClientManager {
	// fmt.Println("New client...")
	//attempt to connect to all fixed number of clients
	manager := newClientManager(d)
	return manager
}
func (manager *ClientManager) connectToClients(dfs *Dfs) {
	for i,_ := range dfs.clients {
		go tryConnect(i,manager.id,dfs)
	}
}
func tryConnect(port int,myID int,dfs *Dfs){
	client:=connectToLocalHost(port,myID)//recieve client
	dfs.clients[port]=client
	go client.receive(dfs)
} 

//encoding/decoding functions

//loop that tries to connect to port and return once able to do so
func connectToLocalHost(port int,myID int) *Client{
	var connection net.Conn
	var err error
	for true{
		connection, err = net.Dial("tcp", "localhost:"+strconv.Itoa(port))
		if err != nil {
			// log.Println(err)
			continue//try another time
		}else{break} //connected
	}
			
	var arr [20]byte
	copy(arr[:],strconv.Itoa(myID))
	connection.Write(arr[:])
	client := &Client{id: myID,destPort:port, socket: connection, data: make(chan []byte)}
	
	return client
}

func (client *Client) recieveInterface(ch chan string,d *Dfs){
	for {
		message := make([]byte, 4096)
		_, err := client.socket.Read(message)

		if err != nil {
			client.socket.Close()
			break
		}

		//decode
		msg:=bytestoString(message)
		if(msg=="state"){
			client.data<-encodeMsg(d.getCurrentState())
			continue
		}
		ch<-msg
	}	
}
func (client *Client) sendStateToServer(){
	defer client.socket.Close()

	//send State to the Dfs starter 
	for {
		select {
		case message, _ := <-client.data:
			_,err:=client.socket.Write(message)
			if err!=nil{
				log.Println(err)
			}
		}
	}
}
func getClientUIServer(d *Dfs,port int,myID int,ch chan string) *Client{
	client:=connectToLocalHost(port,myID)
	go client.recieveInterface(ch,d)
	go client.sendStateToServer()
	return client
}


func bytestoString(b []byte) string {
	n := bytes.IndexByte(b, 0)
	return string(b[:n])
}


func encodeMsg(rmsg interface{}) []byte {
	var ref bytes.Buffer
	enc := gob.NewEncoder(&ref)

	err := enc.Encode(rmsg)
	logEncDecError(err, "encode remote")
	return ref.Bytes()
}

func decodeMsg(bits []byte) interface{} {
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
