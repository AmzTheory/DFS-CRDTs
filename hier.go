package main

import (
	// "fmt"
	"fmt"
	"strings"

	lls "github.com/emirpasic/gods/stacks/linkedliststack"
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
	parent 	 *DfsTreeElement
	children []*DfsTreeElement
}

func (d DfsTreeElement) getPath() string {
	if d.name == "/" {
		return d.name
	}
	if d.fileType == "dir" {
		return d.path+d.name + "/"
	}
	return d.path + d.name
}

//slice  referring to children nodes
type children []*DfsTreeElement

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
	ro := DfsTreeElement{name: "/", fileType: "dir", path: "", content: "",parent:nil,}

	l := hierLayer{root: &ro,
		contentMap: make(map[string]string),
	}

	return &l
}

func (l *hierLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

//Update Interface

//add element
func (l *hierLayer) add(path string, name string, typ string) {
	l.dfs.UpdateAddReplication(path+"/"+name, typ)
}

//remove element
func (l *hierLayer) remove(path string, typ string) {
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
func (l *hierLayer) updateState(cmap map[*replicationElement]string) {
	//go through the map and build the tree
	l.root = &DfsTreeElement{name: "/", fileType: "dir", path: "", content: "",parent:nil,}
	stack := lls.New()


	//policy used here is skip

	stack.Push(l.root)
	// untill stack empty
	for !stack.Empty() {
		// 	pop stack call el
		ra, _ := stack.Pop()
		el := ra.(*DfsTreeElement)

		if el.fileType == "dir" {
			for _, i := range getChildren(el, cmap) {
				ii := i
				stack.Push(&ii)
				el.children = append(el.children, &ii)
			}
			// fmt.Println(el.getPath(), el.children)

		}

	}

	//last step is to send the interface layer with update state
	l.updateInterface()

}


func (l *hierLayer) reappear() {
	//request the content of reappear filed from the replication
}


//pass to interfac
func (l *hierLayer) updateInterface(){
	 l.dfs.updateInterface(l.root)
}

//return to interface
//user interface will be looking it up

//axulariy functions
func (l *hierLayer) printCurrentState() {
	l.printElement(*l.root, 0)
}

func (l *hierLayer) printElement(root DfsTreeElement, nt int) {
	for i := 0; i < nt; i++ {
		fmt.Printf("\t") //print tabs
	}
	// fmt.Println(root.name)
	val := (root.children)
	isDir :=""
	if(root.fileType=="dir"){
		isDir="+"
	}
	fmt.Println(isDir+root.name)
	for i := 0; i < len(val); i++ {
		l.printElement(*val[i], nt+1)
	}

}


func (l *hierLayer) runDown(ui chan UiToHier,rep chan HierToRep){

	for{ 
		msgu := <- ui //receiving from ui layer 
		msgR :=HierToRep{
						path: msgu.path+msgu.name, 
						fileType: msgu.fileType, 
						op:      msgu.op,
					}


		rep <-msgR  //sending message to replication layer
	}
}
func (l* hierLayer) runUp(rep chan map[*replicationElement]string ,ui chan *DfsTreeElement){
	for{
		msgr:=<-rep
		//apply the policies
		l.updateState(msgr)
		ui <-l.root //send the root to ui
	}
}

func skip(map[*replicationElement]string){
	//iterate thro
}


func findRoot(cmap map[*replicationElement]string) string {
	for k := range cmap {
		if !strings.Contains(k.name, "/") {
			return k.name //root found
		}
	}
	return ""
}
func pathAndName(str string) (string, string) {
	li := strings.LastIndex(str, "/")
	return str[:li+1], str[li+1:]
}
func getChildren(root *DfsTreeElement, cmap map[*replicationElement]string) []DfsTreeElement {
	path:=root.getPath()
	temp := []DfsTreeElement{}
	for k := range cmap {
		p, n := pathAndName(k.name)
		if p == path && k.name != "/" {
			el := DfsTreeElement{name: n,
				fileType: k.elementType,
				path:     p,
				children: []*DfsTreeElement{},
				parent: root,
			}
			temp = append(temp, el)
		}

	}
	return temp
}
