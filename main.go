package main

import (
	crdt "CRDTsGO"

	"time"

	set "github.com/emirpasic/gods/sets/linkedhashset"
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

	// rmap:=generateReplicasMap(5)
	// cont:=NewController()
	// cont.create(rmap)
	// cont.SetUpConnection()
	// cont.SetUpCommunication()

	//DFS is ready for operation execution/communication between replicas

	aContent := []string{"a", "b", "c"}
	bContent := []string{"d", "a", "e", "g"}

	aRmContent := []string{}
	bRmContent := []string{"g"}

	a := crdt.NewORSet()
	b := crdt.NewORSet()

	// //channels
	aAdd := make(chan Apair)
	aRm := make(chan Rpair)

	bAdd := make(chan Apair)
	bRm := make(chan Rpair)

	go listenOR(a, aAdd, aRm)
	go listenOR(b, bAdd, bRm)

	//pass to other set
	go setOp(a, bAdd, bRm, aContent, aRmContent)
	go setOp(b, aAdd, aRm, bContent, bRmContent)

	time.Sleep(5 * time.Second)

	a.PrintElements()
	b.PrintElements()

}

//Add Pair
type Apair struct {
	u string
	e string
}

//rm Pair
type Rpair struct {
	R *set.Set
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
