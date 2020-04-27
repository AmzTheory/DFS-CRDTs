package main

import (
	"os"
	"strconv"
	"strings"
	"fmt"
)

// func main() {

// 	dfs := newDfs(1)
// 	//init DFS
// 	dfs.start()
// 	// //start the DFS with concurrency
// 	dfs.runAll()

// }

func main() {
	//id  list of clients TESTServer
	id,_:=strconv.Atoi(os.Args[1])
	cls:=os.Args[2]
	clients:=getMapClients(strings.Split(cls,","))
	serv,_:=strconv.Atoi(os.Args[3])//
	fmt.Println(clients)
	dfs:=newDfs(id,clients,serv)
	dfs.start()
	dfs.runAll()

	

}

func getMapClients(ls []string) map[int]*Client{
	clients:=map[int]*Client{}
	var a int
	for _,v:=range ls{
		a,_=strconv.Atoi(v)
		clients[a]=nil
	}
	return clients
}