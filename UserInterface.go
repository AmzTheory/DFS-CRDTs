package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	dfsTree := gotree.New(format(*l.root))
	printDfsHelper(&dfsTree, l.root.children)
	fmt.Println(dfsTree.Print())
}
func printDfsHelper(root *gotree.Tree, children []*DfsTreeElement) {
	for _, i := range children {
		subTree := (*root).Add(format(*i))
		if len(i.children) != 0 {
			printDfsHelper(&subTree, i.children)
		}
	}
}

func (l *UserInterface) wait() {
	/**
	cd  change current Directory
	ls  show files in current directory
	mk [name] [type]  create file/folder in the current directory
	rm [name] [type]  remove file/folder  //
	printFs	print the entire Dfs
	quit   close the program(go offline)
	*/

	currentDir := l.root
	reader := bufio.NewReader(os.Stdin)

	//infinite loop
	for {
		fmt.Print(currentDir.getPath() + "->")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", 1)
		// fmt.Println(text)

		words := strings.Split(text, " ")

		var el DfsTreeElement
		command := words[0]
		if command == "ls" {

			for i := 0; i < len(currentDir.children); i++ {
				el = *currentDir.children[i]
				if el.fileType == "dir" {
					fmt.Println("\t+" + el.name + "\t" + el.fileType)
				} else {
					fmt.Println("\t" + el.name + "\t" + el.fileType)
				}
			}
		} else if command == "cd" {
			dirName := words[1]
			found := false
			//check if change directory to parent and current has parentt
			if dirName == ".." && currentDir.parent != nil {
				currentDir = currentDir.parent
				found = true
			}

			for i := 0; i < len(currentDir.children) && !found; i++ {
				el = *currentDir.children[i]
				if el.fileType == "dir" && el.name == dirName {
					currentDir = &el
					found = true
				}

			}

			if !found {
				fmt.Println("\t make sure " + dirName + " is directory and does exist")
			}

		} else if command == "mk" {
			name:= words[1]
			fileType:=words[2]
			l.dfs.updateAddHier(currentDir.getPath(),name,fileType)
			fmt.Println(currentDir.children,l.dfs.hier.root)
		} else if command == "rm" {
			// name:= words[1]
			// fileType:=words[2]
			// l.dfs.updateAddHier(currentDir.getPath(),name,fileType)
		} else if command == "printFs" {
			l.printDfs()
		} else if command == "quit" {
			fmt.Println("DFS is closed")
			break
		} else {
			fmt.Println("->" + command + " Unknown command")
		}
	}

}
func (l *UserInterface) updateState(root *DfsTreeElement) {
	l.root = root
	//l.printDfs()
}
func format(el DfsTreeElement) string {
	if el.fileType == "dir" {
		return el.name
	}
	return el.name + "." + el.fileType
}
func getChildrenOfCurrent() {

}
