package main

import (
	crdt "CRDTsGO"
	// "fmt"
	// set "github.com/emirpasic/gods/sets/linkedhashset"
)

// func main() {

// 	dfs := newDfs(1)
// 	//init DFS
// 	dfs.start()
// 	// //start the DFS with concurrency
// 	dfs.runAll()

// }

func main() {

	rmap := generateReplicasMap(2)
	cont := NewController()
	cont.create(rmap)
	cont.SetUpConnection()
	cont.SetUpCommunication()

	// fmt.Println(rmap)
	cont.run()

	

}

//Add Pair
type Apair struct {
	u string
	e string
}

//rm Pair
type Rpair struct {
	R []interface{}
	e string
}

func listenOR(set *crdt.ORSet, add chan Apair, rm chan Rpair) {
	for {
		select {
		case p := <-add:
			set.Add(p.u, p.e)
		case q := <-rm:
			set.Remove(q.R, q.e)

		}
	}
}

func setOp(set *crdt.ORSet, aSend chan Apair, rmSend chan Rpair, cotent []string, rmcontent []string) {

	//add elements
	for _, el := range cotent {
		u := set.AddL(el)
		aSend <- Apair{u: u, e: el}
	}

	//remove elements
	for _, el := range rmcontent {
		r := set.RemoveL(el)

		rmSend <- Rpair{R: r, e: el}
	}
}
