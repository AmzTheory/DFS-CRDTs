package main

import (
		"fmt"
		"github.com/disiqueira/gotree"
		)

/*

	main functionality
		display the view
		communicate operataions to heir layer

	fields
	operations
		print out the DFS Interface
		all DFS operations


*/
type UserInterface struct {
	root *DfsTreeElement
	dfs  *Dfs
}

func newUserInteface(r *DfsTreeElement, d *Dfs) *UserInterface {
	return &UserInterface{root: r, dfs: d}
}

func (l UserInterface) printDfs() {
	dfsTree:=gotree.New(format(*l.root))
	printDfsHelper(&dfsTree, l.root.children)
	fmt.Println(dfsTree.Print())
}
func printDfsHelper(root *gotree.Tree,children []*DfsTreeElement){
	for _, i := range children {
		subTree:=(*root).Add(format(*i))
		if(len(i.children)!=0){
			printDfsHelper(&subTree,i.children)
		}
	}
}


func (l *UserInterface) wait() {

}
func (l *UserInterface) updateState(root *DfsTreeElement){
	l.root=root
	l.printDfs()
}
func format(el DfsTreeElement) string {
	if(el.fileType=="dir"){
		return el.name
	}
	return el.name+"."+el.fileType
}
