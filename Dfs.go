package main

//"fmt"

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
	hier *hierLayer
	rep  *replicationLayer
}

func newDfs() *Dfs {
	d := Dfs{hier: newhierLayer(), rep: newReplicationLayer()}
	return &d
}
func (d *Dfs) start() {
	(*d).rep.setDfs(d)
	(*d).hier.setDfs(d)

	//add root
	//pass it to to interface

}

//downwards

//User interface to Hier
func (dfs Dfs) updateAddHier(path string, n string, typ string) {}
func (dfs Dfs) updateRemoveHier(path string, typ string)        {}

//update

//Hier to replication
func (dfs Dfs) UpdateAddReplication(path string, typ string) {
	dfs.rep.add(path, typ)
}
func (dfs Dfs) UpdateRemoveReplication(path string, typ string) {
	dfs.rep.remove(path, typ)
}

//replication to other replicas (future)

//upwards

func (dfs Dfs) updateHier(cmap map[string]string) {
	//	dfs.
}
func (dfs Dfs) updateInterface(tree DfsTree) {}
