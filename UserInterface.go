package main

import (
	"fmt"
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
type interfaceLayer struct {
	root *DfsTreeElement
}

func newUserInteface(root *DfsTreeElement) {
	fmt.Println("Ahmed")
}

func (l interfaceLayer) printDfs() {

}

func (l *interfaceLayer) wait() {

}
