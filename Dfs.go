package main

// "time"

// "fmt"
import (
	"context"
)

/*
	instance of
		UserInterface
		Hier
		replication	x
	Model communication
		between layers

	View of the DFS

	Comunication between replicas
		assign operations
*/

type Dfs struct {
	id      int
	ui      *UserInterface
	hier    *HierLayer
	rep     *RepLayer
	manager *ClientManager
	clients map[int]*Client
	uiID    int
	uiServ  *Client
	rem     chan RemoteMsg
}


//messages type
type UiToHier struct {
	path     string
	name     string
	fileType string
	op       string
	cancel   context.CancelFunc
}

type HierToRep struct {
	path     string
	fileType string
	op       string
	cancel   context.CancelFunc
}

const (
	channellen = 5
)

var on bool

func newDfs(id int, clients map[int]*Client, servID int) *Dfs {
	h := newHierLayer()
	r := newRepLayer(id)
	d := Dfs{id: id, clients: clients, hier: h, rep: r, rem: make(chan RemoteMsg, channellen), uiID: servID}
	d.hier.dfs = &d
	d.rep.dfs = &d
	return &d
}

func (d *Dfs) runAll() {

	//channels
	uiTohier := make(chan UiToHier, channellen)
	hierTorep := make(chan HierToRep, channellen)

	repTohier := make(chan map[*RepElem]string, channellen)
	hierToui := make(chan *DfsNode, channellen)

	// remToRep = make(chan RemoteMsg)
	execOp := make(chan RemoteMsg, channellen)

	input := make(chan bool)

	//go routines
	go d.hier.runDown(uiTohier, hierTorep) //run hier  top->down
	go d.hier.runUp(repTohier, hierToui)   //run hier  down -> top
	go d.rep.runLocally(execOp, hierTorep) //run rep local thread
	go d.rep.runRemotely(execOp, d.rem)
	go d.rep.executeOp(repTohier, execOp)

	d.ui.recieveInitialRoot(hierToui)
	go d.ui.runRecieve(hierToui)
	go d.ui.run(uiTohier, input)

	d.manager = newClient(d)

	// time.Sleep(4*time.Second)//this enforce case(1) to occur

	d.manager.connectToClients(d)

	<-input //Dfs gods offline

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
	/**
		connect to the testing server
	**/
	ch := make(chan string, channellen)
	d.uiServ = getClientUIServer(d, d.uiID, d.id, ch)

	d.ui = newUserInteface(d.hier.root, d, ch)
}

func (d *Dfs) closeAll() {
	d.rep.writeDB()
}
func (d *Dfs) closeClients() {
	for _, v := range d.clients {
		if v != nil {
			v.socket.Write([]byte("close"))
		}
		// v.socket.Close()
	}
}
