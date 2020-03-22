package main

import (
	"fmt"
	"time"
)

// func main() {

// 	dfs := newDfs(1)
// 	//init DFS
// 	dfs.start()
// 	// //start the DFS with concurrency
// 	dfs.runAll()

// }

func main() {
	d1 := newDfs(1001, []int{1002,1003})
	d2 := newDfs(1002, []int{1001,1003})
	d3  := newDfs(1003, []int{1001,1002})
	// d3 := newDfs(1003, []int{1001, 1002})

	msg := RemoteMsg{ClientID: 1, Msg: "This is Ahmed"}
	ds := []*Dfs{d1, d2,d3}

	//initate dfs and listener
	for _, d := range ds { 		
		d.start()
		b := make(chan bool)
		go d.runAll(b)

		//indicate the dfs listener is open
		<-b
	}

	//connect to other clients
	for _, d := range ds {
		fmt.Printf("connect %d \n",d.id)
		d.startConnecting()
	}


	time.Sleep(3*time.Second)//wait for all connections
	for _,d :=range ds{
		d.sendRemote(msg)
		time.Sleep(3*time.Second)//wait for message to be sent
	}
	//l
	// for true {

	// }

}
