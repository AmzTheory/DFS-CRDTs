package main

import (
)

// func main() {

// 	dfs := newDfs(1)
// 	//init DFS
// 	dfs.start()
// 	// //start the DFS with concurrency
// 	dfs.runAll()

// }

func main() {

	rmap:=generateReplicasMap(5)
	cont:=NewController()
	cont.create(rmap)
	cont.SetUpConnection()
	cont.SetUpCommunication()

	//DFS is ready for operation execution/communication between replicas

}
