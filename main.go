package main

// 		"os"

// 		"strings"
// 		"bufio"

import (

)

type person struct {
	name string
	age  *int
}

func main() {

	dfs := newDfs(1)
	//init DFS
	dfs.start()
	// //start the DFS with concurrency
	dfs.runAll()


}

