package main

import (
	"time"
)

// "fmt"

/*
	instance of
		UserInterface
		Hier
		replication
	Model communication
		between layers

	View of the DFS

	Comunication between replicas
		assign operations
*/

type Dfs struct {
	id      int
	ui      *UserInterface
	hier    *hierLayer
	rep     *replicationLayer
	manager *ClientManager
	clients map[int]*Client
	servID  int
	dfsServ *Client 
	rem     chan RemoteMsg
}

//messages type
type UiToHier struct {
	path     string
	name     string
	fileType string
	op       string
}

type HierToRep struct {
	path     string
	fileType string
	op       string
}

var on bool

func newDfs(id int, clients map[int]*Client,servID int) *Dfs {
	d := Dfs{id: id, clients: clients, hier: newhierLayer(), rep: newReplicationLayer(id,true),rem:make(chan RemoteMsg),servID:servID}
	return &d
}

// var (
// 	remToRep chan RemoteMsg
// )

func (d *Dfs) runAll() {
	on = true

	//channels
	uiTohier := make(chan UiToHier)
	hierTorep := make(chan HierToRep)

	repTohier := make(chan map[*replicationElement]string)
	hierToui := make(chan *DfsTreeElement)

	// remToRep = make(chan RemoteMsg)
	execOp := make(chan RemoteMsg)

	input:=make(chan bool)

	//go routines
	go d.hier.runDown(uiTohier, hierTorep) //run hier  top->down
	go d.hier.runUp(repTohier, hierToui)   //run hier  down -> top
	go d.rep.runLocally(execOp, hierTorep) //run rep local thread
	go d.rep.runRemotely(execOp, d.rem)
	go d.rep.pushUpState(repTohier, execOp)

	d.ui.recieveInitialRoot(hierToui)
	go d.ui.runRecieve(hierToui)
	go d.ui.run(uiTohier,input)

	d.manager = newClient(d)
	
	// time.Sleep(4*time.Second)//this enforce case(1) to occur
	
	d.manager.connectToClients(d)

	<-input  //Dfs gods offline

	d.closeClients()
	//get the data from DB
	d.rep.writeDB()

}
func (d *Dfs) getCurrentState() []interface{} {
	return d.rep.or.Values()
}

//triggers to send remote operation to other clients
func (d *Dfs) sendRemote(msg RemoteMsg) {
	(d.manager.broadcast) <- msg
}

func (d *Dfs) start() {
	d.rep.setDfs(d)
	d.hier.setDfs(d)
	
	/**
		connect to the testing server
	**/
	ch:=make(chan string)
	d.dfsServ=getClientUIServer(d,d.servID,d.id,ch)
	
	d.ui = newUserInteface(d.hier.root, d,ch)
}

//downwards

//User interface to Hier
// func (d *Dfs) updateAddHier(path string, n string, typ string) {
// 	d.UpdateAddReplication(path+n, typ)
// }
// func (d *Dfs) updateRemoveHier(path string, typ string) {
// 	d.UpdateRemoveReplication(path, typ)
// }

//update

// //Hier to replication
// func (d *Dfs) UpdateAddReplication(path string, typ string) {
// 	d.rep.add(path, typ)
// }
// func (d *Dfs) UpdateRemoveReplication(path string, typ string) {
// 	d.rep.remove(path, typ)
// }

//replication to other replicas (future)

//upwards

func (d *Dfs) updateHier(cmap map[*replicationElement]string) {
	d.hier.updateState(cmap) //infor hier layer
}

func (d *Dfs) updateInterface(root *DfsTreeElement) {
	d.ui.updateState(root)
}
func (d *Dfs) closeAll() {
	d.rep.writeDB()
}
func (d *Dfs) closeClients(){
	for _,v:=range d.clients{
		if(v!=nil){
		v.socket.Write([]byte("close"))
		}
		// v.socket.Close()
	}
}

