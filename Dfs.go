package main

// import "fmt"

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
	ui  *UserInterface
	hier *hierLayer
	rep  *replicationLayer
}

func newDfs() *Dfs {
	d := Dfs{hier: newhierLayer(), rep: newReplicationLayer()}
	return &d
}
func (d *Dfs) printInstanceRef() {
	//	fmt.Println("calling reference",&d)
}
func (d *Dfs) start() {
	d.rep.setDfs(d)
	d.hier.setDfs(d)
	d.ui=newUserInteface(d.hier.root,d)

	d.ui.printDfs()

}

//downwards

//User interface to Hier
func (d *Dfs) updateAddHier(path string, n string, typ string) {
	d.UpdateAddReplication(path+n, typ)
}
func (d *Dfs) updateRemoveHier(path string, typ string) {
	d.UpdateRemoveReplication(path,typ)
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
