package main

import (
		"fmt"
		set "github.com/emirpasic/gods/sets/linkedhashset"			
		"strconv"
	)


type ClientManager struct {
	active     *set.Set
	offline    *set.Set
    broadcast  chan RemoteMsg
    register   chan *Client
    unregister chan *Client
}

type Client struct {
	id 	   int
    socket net.Conn
    data   chan []RemoteMsg
}
type RemoteMsg struct{
	clientId	int
	msg			string
}


func newClientManager() *ClientManager {
	
	d := Dfs{id:id,hier: newhierLayer(), rep: newReplicationLayer(id)}
	return &d
}

func startServer(){
	fmt.Println("Starting server...")
    listener, error := net.Listen("tcp", ":1234")
	
	if error != nil {
        fmt.Println(error)
	}
	    
	//we've got 3 replicas for now
	
	manager := ClientManager{
		active:		set.New(),//empty set (will be added when connected)
		offline:	set.New(),
        broadcast:  make(chan RemoteMsg),
        register:   make(chan *Client),
        unregister: make(chan *Client),
	}
	
    go manager.start()
    go manager.waitForConns()
}

func (manager *ClientManager) waitForConns(){
	for {
        connection, _ := listener.Accept()
        if error != nil {
            fmt.Println(error)
        }
        client := &Client{socket: connection, data: make(chan string)}
        manager.register <- client
        go manager.receive(client)
        go manager.send(client)
    }
}

func (manager *ClientManager) start() {
    for {
        select {
		case connection := <-manager.register:
			 manager.setClientOnline(&connection)
             fmt.Println("Added new connection!")
        case connection := <-manager.unregister:
			manager.setClientOffline(&connection)	
			fmt.Println("A connection has terminated!")
        case message := <-manager.broadcast:
            for connection := range manager.active {
				con:=connection.(*Client)
				if con.id!=message.clientId{
					select {
                		case connection.data <- message:
					default:
						manager.setClientOffline(&con)
					}		
				}
				
            }
        }
    }
}
func (manager *ClientManager) setClientOffline(con *Client){
	close(con.data)
	manager.offline.Add(con.id)
	manager.active.Remove(&con)
}
func (manager *ClientManager) setClientOnline(con *Client){
	manager.active.Add(&con)
	manager.remove.Remove(con.id)
}

func (client *Client) receive() {
	///keep waiting for remote operations
	for {
        length, err := client.socket.Read(message)
        if err != nil {
            client.socket.Close()
            break
		}
		
        sendtorep()
    }
}

func (ClientManager *manager) receive(client *Client){
	//recieve the id
	message := make(string)
	length, _ := client.socket.Read(message)
	client.id=strconv.Atoi(message)//set the id

	
	//wait for local operations to be sent	
	for {
        length, err := client.socket.Read(message)
        if err != nil {
            manager.unregister <- client
            client.socket.Close()
            break
        }
        if length > 0 {
            fmt.Println("RECEIVED: " + string(message))
            manager.broadcast <- message
        }
    }
}
func (ClientManager *manager) send(client *Client){
	 defer client.socket.Close()
	 //send remote operations as you recieve msg in the channel
	 for {
        select {
        case message, ok := <-client.data:
            if !ok {
                return
            }
            client.socket.Write(message)
        }
    }
}
