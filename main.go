package main

// 		"os"

// 		"strings"
// 		"bufio"

type person struct {
	name string
	age  *int
}

func main() {

	dfs := newDfs()
	dfs.start()
	dfs.updateAddHier("/", "1st", "txt")
	dfs.updateAddHier("/", "2nd", "txt")
	dfs.updateAddHier("/", "3rd", "txt")
	dfs.updateAddHier("/", "folder", "dir")
	dfs.updateAddHier("/folder/", "one", "txt")
	dfs.updateAddHier("/folder/", "folder2", "dir")
	dfs.updateAddHier("/folder/folder2/", "rand", "txt")
	// dfs.updateRemoveHier("/1st", "txt")

	//start the DFS with concurrency
	dfs.runAll()
}
