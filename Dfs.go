package main

import (
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
	clients []int
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

func newDfs(id int, clients []int) *Dfs {
	d := Dfs{id: id, clients: clients, hier: newhierLayer(), rep: newReplicationLayer(id)}
	return &d
}

func (d *Dfs) runAll(ch chan bool) {
	on = true

	//channels
	uiTohier := make(chan UiToHier)
	hierTorep := make(chan HierToRep)

	repTohier := make(chan map[*replicationElement]string)
	hierToui := make(chan *DfsTreeElement)

	remToRep:=make(chan RemoteMsg)
	execOp:=make(chan RemoteMsg)

	//go routines
	go d.ui.run(uiTohier, hierToui)           //run ui
	go d.hier.runDown(uiTohier, hierTorep)    //run hier  top->down
	go d.hier.runUp(repTohier, hierToui)      //run hier  down -> top
	go d.rep.runLocally(execOp, hierTorep) //run rep local thread
	go d.rep.runRemotely(execOp,remToRep)
	go d.rep.pushUpState(repTohier,execOp)
	
	
	d.manager = newClient(d.id)
	ch <- true
	//Wait for ever
	for on {
		//break when DFS closed

	}

	//get the data from DB
	d.rep.writeDB()

}
func (d *Dfs) startConnecting(){
	d.manager.connectToClients(d)
}

//triggers to send remote operation to other clients
func (d *Dfs) sendRemote(msg RemoteMsg) {
	(d.manager.broadcast) <- msg
}
func (d *Dfs) sendRemoteToRep (msg RemoteMsg) {
	//locked
	
	
}

func (d *Dfs) printInstanceRef() {
	//	fmt.Println("calling reference",&d)
}
func (d *Dfs) start() {
	d.rep.setDfs(d)
	d.hier.setDfs(d)
	d.ui = newUserInteface(d.hier.root, d)

	d.ui.printDfs()
}

//downwards

//User interface to Hier
func (d *Dfs) updateAddHier(path string, n string, typ string) {
	d.UpdateAddReplication(path+n, typ)
}
func (d *Dfs) updateRemoveHier(path string, typ string) {
	d.UpdateRemoveReplication(path, typ)
}

//update

//Hier to replication
func (d *Dfs) UpdateAddReplication(path string, typ string) {
	d.rep.add(path, typ)
}
func (d *Dfs) UpdateRemoveReplication(path string, typ string) {
	d.rep.remove(path, typ)
}

//replication to other replicas (future)

//upwards

func (d *Dfs) updateHier(cmap map[*replicationElement]string) {
	d.hier.updateState(cmap) //infor hier layer
}

func (d *Dfs) updateInterface(root *DfsTreeElement) {
	d.ui.updateState(root)
}
func (d *Dfs) closeAll() {
	on = false
}
