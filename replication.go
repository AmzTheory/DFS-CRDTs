package main

import (
	"fmt"

	set "github.com/emirpasic/gods/sets/linkedhashset"
)

/*
fields
	SET (assume its and ORSET) string
	MAP (assume its CRDT)  String->String
functions
	add(p,t)
	remove(p,t)
	update(p,t,u)

returns to upper layer
	MAP   (p,t)->String (content)
*/

//Structs and type
type replicationElement struct {
	name        string
	elementType string
}

// type elementSet []*replicationElement
type elementSet *set.Set
type contentMap map[*replicationElement]string

type replicationLayer struct {
	dfs  *Dfs
	set  elementSet
	cmap contentMap
}



//initalisation
func newReplicationLayer() *replicationLayer {
	el := replicationElement{name: "/",
		elementType: "dir"}
	// s := []*replicationElement{&el}
	s := set.New()
	s.Add(el)
	dic := make(map[*replicationElement]string)
	dic[&el] = ""

	l := replicationLayer{
		dfs:  new(Dfs),
		set:  s,
		cmap: dic,
	}

	return &l
}
func (l *replicationLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

func (l *replicationLayer) runLocally(send chan map[*replicationElement]string,recieve chan HierToRep){
	send<-l.returnCurrentSet() //send the initial state
	for{
		msg:=<-recieve
		if(msg.op=="add"){
			l.add(msg.path,msg.fileType)
		}else if(msg.op=="rm"){
			l.remove(msg.path,msg.fileType)
		}

		send<-l.returnCurrentSet() //send the updated set to hier
	}
}


//update inteface

func (l *replicationLayer) add(path string, typ string) {
	el := replicationElement{name: path, elementType: typ}
	// l.set = append(l.set, &el)
	(*l.set).Add(el) //element get added
	l.cmap[&el] = "" //initate with an empty content
	// l.updateDfs()
	fmt.Println("added", path)
}

func (l *replicationLayer) remove(path string, typ string) {
	//remove an element from the slice
	// temp := set.New()
	for _, i := range (*l.set).Values() {
		ii := i.(replicationElement)
		if (ii.name == path && ii.elementType == typ) {
			(*l.set).Remove(ii)
		}
	}
	// l.set = temp
	// fmt.Println((*l.set).Size(), temp.Size())
	// l.updateDfs()
	fmt.Println("removed", path)
}

// func (l *replicationLayer) udpate(path string,typ string){
// 	fmt.Println("element has been added")
// }

//update hier by through dfs
func (l *replicationLayer) updateDfs() {

	l.dfs.updateHier(l.returnCurrentSet()) //select only one the exist in the setS
}

func (l *replicationLayer) returnCurrentSet() map[*replicationElement]string {
	temp := make(map[*replicationElement]string)
	for _, k := range (*l.set).Values() {
		kk := (k.(replicationElement))
		temp[&kk] = l.cmap[&kk]
	}
	return temp
}

func (l *replicationLayer) printCurrentState() {
	fmt.Println("\nCRDT_Set\n-------------")
	// for _, k := range l.set {
	// 	v := l.cmap[k]
	// 	fmt.Println("", k.name, "content", v)
	// }
	for _, k := range (*l.set).Values() {
		kk := (k.(replicationElement))
		v := l.cmap[&kk]
		fmt.Println("", kk.name, "content", v)
	}
	fmt.Println()
}
