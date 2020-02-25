package main

import (
	//"fmt"
	"fmt"
	"strings"
)

/*
fields
	TREE (n,t,p,content)
update Operations
	add(p,n,t)
	remove(p,t)
	update(p,n,t,u)
		for each particular policy must be adopted

obtains from the upper layer
	Map (p,t) -> content
		for all (p,t) that are in the set

returns to upper layer  (assume hier is upper layer)
	Tree with elements
		(name,type,path,content)
*/

//struct definitions

//an element of the tree

type DfsTreeElement struct {
	name     string
	fileType string
	path     string
	content  string
	children children
}

//slice  referring to children nodes
type children []DfsTreeElement

//store the root
type DfsTree DfsTreeElement

//hier layer
type hierLayer struct {
	dfs        *Dfs
	root       *DfsTreeElement
	contentMap map[string]string
}

//initalisation
func newhierLayer() *hierLayer {
	ro := DfsTreeElement{name: "root", fileType: "dir", path: "", content: ""}

	l := hierLayer{root: &ro,
		contentMap: make(map[string]string),
	}

	return &l
}

func (l hierLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

//Update Interface

//add element
func (l hierLayer) add(path string, name string, typ string) {
	l.dfs.UpdateAddReplication(path+"/"+name, typ)
}

//remove element
func (l hierLayer) remove(path string, typ string) {
	l.dfs.UpdateRemoveReplication(path, typ)
}

// func (tree DfsTree) update(path string,name string,typ string){
// 	fmt.Println("Element has been updated")
// }

//update lower layer
func updateReplation() {
	//communicate it to the Dfs instance
}

//modify the state based on new info from replication
func (l hierLayer) updateState(cmap map[string]string) {
	//go through the map and build the tree
	fmt.Print()
	// (Depth first)
	// find the root , instantiate and add into stack
	// ro := findRoot(cmap)
	// //ignore type
	// p,n:=pathAndName(ro)
	// rootEl:=DfsTreeElement{name:n,fileType:"uf",path:p,children:[]DfsTreeElement{}}
	// then add root to stack
	// untill stack empty
	// 	pop stack call el
	// 	iterate throu children
	// 		instantiate
	// 		add reference to stack
	// 		add reference to el children

	//pass to interface
}

//return to interface
//user interface will be looking it up

//axulariy
func findRoot(cmap map[string]string) string {
	for k,_ :=range cmap{
		if(!strings.Contains(k,"/")){
			return k  //root found
		}
	}
	return ""
}
func pathAndName(str string)(string,string) {
	li := strings.LastIndex(str, "/")
	return str[:li+1],str[li+1:]
}

