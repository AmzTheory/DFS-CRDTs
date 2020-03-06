package main

import (
// 		"os"
		
		crdt "DFS/CRDTsGO"
// 		"strings"
// 		"bufio"
	)


type person struct {
	name string
	age  *int
}

func main() {
	

	
	// dfs := newDfs()
	// dfs.start()
	// dfs.updateAddHier("/", "1st", "txt")
	// dfs.updateAddHier("/", "2nd", "txt")
	// dfs.updateAddHier("/", "3rd", "txt")
	// dfs.updateAddHier("/", "folder", "dir")
	// dfs.updateAddHier("/folder/", "one", "txt")
	// dfs.updateAddHier("/folder/", "folder2", "dir")
	// dfs.updateAddHier("/folder/folder2/", "rand", "txt")
	// dfs.updateRemoveHier("/1st","txt")

	or:=crdt.NewORSet()
	id:=or.AddSrc("a")
	or.AddDownStream(id,"a")
	// or.PrintElements()
}
