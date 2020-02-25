package main

import (
	"fmt"
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
type elementSet []replicationElement
type contentMap map[string]string

type replicationLayer struct {
	dfs  *Dfs
	set  elementSet
	cmap contentMap
}

//initalisation
func newReplicationLayer() *replicationLayer{
	el:=replicationElement{name:"root",
							elementType: "dir",}
	s:=[]replicationElement{el}
	
	dic:=make(map[string]string)
	dic[el.name+","+el.elementType]=""

	l:= replicationLayer{
		set:  s,
		cmap: dic,
	}

	return &l
}
func (l replicationLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

//update inteface

func (l replicationLayer) add(path string, typ string) {
	el := replicationElement{name:path,elementType: typ,}
	l.set = append(l.set, el)
	l.cmap[el.name+""+el.elementType] = "" //initate with an empty content
	l.updateDfs()
}

func (l replicationLayer) remove(path string, typ string) {
	//remove an element from the slice
	temp := []replicationElement{}
	for _,i := range l.set {
		if !(i.name == path && i.elementType == typ) {
			temp = append(temp, i)
		}
		fmt.Println(i)
	}
	l.set = temp
}

// func (l replicationLayer) udpate(path string,typ string){
// 	fmt.Println("element has been added")
// }

//update hier by through dfs
func (l replicationLayer) updateDfs() {
	l.dfs.updateHier(l.cmap) //select only one the exist in the setS
}
