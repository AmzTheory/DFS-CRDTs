package main

import (
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

type DfsNode struct {
	name     string
	fileType string
	path     string
	content  string
	parent 	 *DfsNode
	children map[string]*DfsNode
}

func (d DfsNode) getPath() string {
	if d.name == "/" {
		return d.name
	}
	if d.fileType == "dir" {
		return d.path+d.name + "/"
	}
	return d.path + d.name
}

//slice  referring to children nodes
type children []*DfsNode

//store the root
type DfsTree DfsNode

//hier layer
type HierLayer struct {
	dfs        *Dfs
	root       *DfsNode
	contentMap map[string]string
}

const (
	skip int=1
	compact int=2
)

//initalisation
func newHierLayer() *HierLayer {
	ro := DfsNode{name: "/", fileType: "dir", path: "", content: "",parent:nil,}

	l := HierLayer{root: &ro,
				   contentMap: make(map[string]string),
	}

	return &l
}

func (l *HierLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}


//update lower layer
func updateReplation() {
	//communicate it to the Dfs instance
}

//modify the state based on new info from replication(Defualt implement skip )
func buildTree(cmap map[*RepElem]string) *DfsNode {
	//go through the map and build the tree

	root := &DfsNode{name: "/", fileType: "dir", path: "", content: "",parent:nil,children:map[string]*DfsNode{}}
	stack := lls.New()

	stack.Push(root)
	// until stack empty
	for !stack.Empty() {
		// 	pop stack call el
		ra, _ := stack.Pop()
		el := ra.(*DfsNode)
		if el.fileType == "dir" {
			for _, i := range getChildren(el, cmap) {
				ii := i
				stack.Push(&ii)
				el.children[ii.name]=&ii
			}
			// fmt.Println(el.getPath(), el.children)

		}

	}

	return root

}


func (l *HierLayer) reappearP(mapping map[*RepElem]string) {
	/*
		find orphan elements
		generated the needed directories and files 
	*/
}

func (l *HierLayer) compactP(mapping map[*RepElem]string) {
	/*
		find orphan elements
		connect the element the first parent that exist 
	*/
	
}

func (l *HierLayer) skipP(mapping map[*RepElem]string) {
	root:=buildTree(mapping)
	l.root=root
}




//axulariy functions
func (l *HierLayer) runDown(ui chan UiToHier,rep chan HierToRep){

	for{ 
		msgu := <- ui //receiving from ui layer 
		msgR :=HierToRep{
						path: msgu.path+msgu.name, 
						fileType: msgu.fileType, 
						op:      msgu.op,
						cancel:  msgu.cancel,
					}


		rep <-msgR  //sending message to replication layer
	}
}
func (l* HierLayer) runUp(rep chan map[*RepElem]string ,ui chan *DfsNode){
	for{
		msgr:=<-rep

		l.skipP(msgr)

		ui <-l.root //send the root to ui
	}
}

func findRoot(cmap map[*RepElem]string) string {
	for k := range cmap {
		if !strings.Contains(k.Name, "/") {
			return k.Name //root found
		}
	}
	return ""
}
func pathAndname(str string) (string, string) {
	li := strings.LastIndex(str, "/")
	return str[:li+1], str[li+1:]
}

func getChildren(root *DfsNode, cmap map[*RepElem]string) []DfsNode {
	path:=root.getPath()
	temp := []DfsNode{}
	for k := range cmap {
		p, n := pathAndname(k.Name)
		if p == path && k.Name != "/" {
			el := DfsNode{name: n,
				fileType: k.ElemType,
				path:     p,
				children: map[string]*DfsNode{},
				parent: root,
			}
			temp = append(temp, el)
		}

	}
	return temp
}
