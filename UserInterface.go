package main

import (
	// "bufio"
	"fmt"
	// "os"
	"strings"
	"context"
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
	root       *DfsNode
	dfs        *Dfs
	currentDir DfsNode
	intchan    chan string
}

func newUserInteface(r *DfsNode, d *Dfs,ch chan string) *UserInterface {
	return &UserInterface{root: r, dfs: d, currentDir: DfsNode{},intchan:ch}
}

func (l UserInterface) printDfs() {
	dfsTree := gotree.New(format(*l.root))
	printDfsHelper(&dfsTree, l.root.children)
	fmt.Println(dfsTree.Print())
}
func printDfsHelper(root *gotree.Tree, children map[string]*DfsNode) {
	for _, i := range children {
		subTree := (*root).Add(format(*i))
		if len(i.children) != 0 {
			printDfsHelper(&subTree, i.children)
		}
	}
}

func (l *UserInterface) recieveInitialRoot(recieve chan *DfsNode) {
	l.root = <-recieve //recieve the inital root
	//run or runBackground can be executed as gorotuines after this method get invoked
	l.currentDir = *l.root
}

func (l *UserInterface) run(send chan UiToHier, status chan bool) {	
	/**
	cd [name/-r]change current Directory
	ls  show files in current directory
	mk [name] [type]  create file/folder in the current directory
	rm [name] [type]  remove file/folder  //
	printFs	print the entire Dfs
	quit   close the program(go offline)
	*/
	// reader := bufio.NewReader(os.Stdin)
	//infinite loop
	var ctx context.Context
	var cancel context.CancelFunc
	var text string
	for {

		fmt.Print(l.currentDir.getPath() + "->")
		// text, _ := reader.ReadString('\n')
		
		// text = strings.Replace(text, "\n", "", 1)
		text=<-l.intchan
		
		fmt.Println(text)
		words := strings.Split(text, " ")

		command := words[0]
		if command == "ls" {
			for _,v:=range l.currentDir.children {
				fmt.Print("\t")
				if (v.fileType=="dir"){
					fmt.Print("+")
				}
			
				fmt.Println(v.name + "\t" + v.fileType)
			}
			
		} else if command == "cd" {
			dirName := words[1]
			found := false
			//check if change directory to parent and current has parent
			if dirName == ".." && l.currentDir.parent != nil {
				l.currentDir = *l.currentDir.parent
				found = true
			}

			if dirName == "-c"  {
				l.currentDir = *l.currentDir.parent
				found = true
			}

			// for i := 0; i < len(l.currentDir.children) && !found; i++ {
			// 	el = l.currentDir.children[i]
			// 	if el.fileType == "dir" && el.name == dirName {
			// 		l.currentDir = *el
			// 		found = true
			// 	}
			// }

			if !found {
				fmt.Println("\t make sure " + dirName + " is directory and does exist")
			}

		} else if command == "mk" {
			if len(words) != 3 {
				fmt.Println("\tmk is defined as : mk  name fileType")
				continue
			}

			name := words[1]
			fileType := words[2]
			if exists(&l.currentDir, name, fileType) {
				fmt.Println("\t" + name + " does exist at the current directory!")
				continue
			}

			ctx, cancel=context.WithCancel(context.Background())
			// l.dfs.updateAddHier(currentDir.getPath(), name, fileType)
			send <- UiToHier{
				path:     l.currentDir.getPath(),
				name:     name,
				fileType: fileType,
				op:       "add",
				cancel:	  cancel,
			}

			<-ctx.Done() //the op has been executed

			//wait for an update
			// l.root = <-recieve
			// currentDir = l.updateNodePointer(currentDir.getPath())

		} else if command == "rm" {
			if len(words) != 3 {
				fmt.Println("\trm is defined as : rm  name fileType")
				continue
			}

			
			name := words[1]
			fileType := words[2]
			if !exists(&l.currentDir, name, fileType) {
				fmt.Println("\t" + name + " does not exist at the current directory!")
				continue
			}

			ctx, cancel=context.WithCancel(context.Background())

			// l.dfs.updateRemoveHier(currentDir.getPath()+name, fileType)
			send <- UiToHier{
				path:     l.currentDir.getPath() + name,
				name:     "",
				fileType: fileType,
				op:       "rm",
				cancel:   cancel,
			}
			//wait for an update
			<-ctx.Done() //the op has been executed
		} else if command == "printfs" {
			l.printDfs()
		} else if command == "help" {
			fmt.Println("\tcd      change current Directory\tcd dir")
			fmt.Println("\tls      show files in current directory\tls")
			fmt.Println("\tmk      create new file/directory\tmk name filetype")
			fmt.Println("\trm  	   remove file/directory\trm name filetype")
			fmt.Println("\tprintfs print the file system tree")
			fmt.Println("\tquit    turn off access mode")
			fmt.Println("\toffline   go offline")
		} else if command == "quit" {
			l.currentDir = *l.root
			status <- true
			ctx, cancel=context.WithCancel(context.Background())
			
			send <- UiToHier{
				op:       "quit",
				cancel:	  cancel,
			}

			<-ctx.Done() //the op has been executed
			break
		}else {
			fmt.Println("->" + command + " Unknown command")
		}
	}

}

func (l *UserInterface) runRecieve(recieve chan *DfsNode) {
	for {
		l.root = <-recieve
		l.currentDir = l.updateNodePointer(l.currentDir.getPath())
		// fmt.Println("UI Recieved ", l.dfs.id)
	}
}

func (l *UserInterface) updateState(root *DfsNode) {
	l.root = root
	//l.printDfs()
}
func (l *UserInterface) updateNodePointer(path string) DfsNode {
	r := l.root
	dirs := strings.Split(path[1:], "/")

	if path == "/" {
		return *r // TODO: improve the current method
	}

	for i := 0; i < len(dirs)-1; i++ {
		// fmt.Println(r.name)
		if a:=findNode(r, dirs[i]);a!=nil{
			r=a
		}else{
			return *r
		}
		
	}
	return *r

}
func findNode(root *DfsNode, dir string) *DfsNode {
	for k,v:=range root.children{
		if k==dir{
			return v 
		} 
	}
	return nil
}

//TODO: expand to go deeper in the tree
func exists(root *DfsNode, path, fileType string) bool {
	for _,el:=range root.children{
		if el.name == path && el.fileType == fileType {
			return true
		}
	}
	return false
}
func format(el DfsNode) string {
	if el.fileType == "dir" {
		return el.name
	}
	return el.name + "." + el.fileType
}
