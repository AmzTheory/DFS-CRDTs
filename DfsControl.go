package main

import (
	// "bufio"
	// "fmt"
	// "os"
	// "strconv"
	// "strings"
)

/*

create number of DFS    (number of replicas)

establish connections among replicas

manual testing (switch between replicas)

Perform operations   (randomly generate operations perform testting)

Testing convergence    when all replicas have recieved


*/

// type DfsController struct {
// 	replicas map[int]*Dfs //map of clients  ports/ids
// 	input    map[int]chan bool
// }

// func NewController() *DfsController {
// 	return &DfsController{replicas: make(map[int]*Dfs), input: make(map[int]chan bool)}
// }

// func (inst *DfsController) create(rep map[int][]int) {
// 	for k, v := range rep {
// 		inst.replicas[k] = newDfs(k, v)
// 	}
// }
// func (inst *DfsController) SetUpConnections() {
// 	for k, v := range inst.replicas {
// 		inst.setUpConnection(k,v)
// 	}
// }

// func (inst *DfsController) setUpConnection(k int,v *Dfs){
// 	v.start()
// 		ch := make(chan bool) //request the access mode
// 		inst.input[k] = ch
// 		b := make(chan bool)
// 		go v.runAll(b, ch)

// 		//indicate the dfs listener is open
// 		<-b
// }

// func (inst *DfsController) SetUpCommunication() {
// 	for _, v := range inst.replicas {

// 		v.startConnecting()
// 	}
// }

// func generateReplicasMap(n int) map[int][]int {
// 	//n^2
// 	repMap := make(map[int][]int)
// 	init := 1001

// 	for i := init; i <= init+n; i++ {
// 		repMap[i] = []int{}
// 		for j := 1001; j <= init+n; j++ {
// 			if j != i {
// 				repMap[i] = append(repMap[i], j)
// 			}
// 		}
// 	}

// 	return repMap
// }

// func (inst *DfsController) run() {
// 	reader := bufio.NewReader(os.Stdin)
// 	var con []string
// 	for {
// 		fmt.Print("type the replica to access\n")
// 		text, _ := reader.ReadString('\n')
// 		text = strings.Replace(text, "\n", "", 1)


// 		if(text=="close"){ 
// 			fmt.Println("Replicas is being written into DB")
// 			inst.closeAll()
// 			fmt.Println("all DFS are closed")	
// 			break
// 		}

// 		if(strings.HasPrefix(text,"offline")){ 
// 			con=strings.Split(text," ")
// 			replica,_ := strconv.Atoi(con[1])
// 			fmt.Println("Replicas",replica ,"is getting closed")
// 			// inst.close(replica)
// 			fmt.Println("DFS is closed")	
// 			break
// 		}

// 		replica, err := strconv.Atoi(text)
		
// 		if (err!=nil) {
// 			continue
// 		}
		
// 		inst.input[replica] <- true //start access mode

// 		<-inst.input[replica] //wait for quit request
// 		// fmt.Println("Checking Convergence", inst.checkConvergence())

// 	}
// }
// func (inst *DfsController) closeAll(){
// 	for _,v:=range inst.replicas{
// 		v.rep.writeDB()
// 	}
// }


// func (inst *DfsController) checkConvergence() bool {
// 	fs := inst.replicas[1001]
// 	res := true
// 	for _, v := range inst.replicas {
// 		if fs.id != v.id {
// 			res = fs.rep.RepEqual(v.rep.or)
// 			if !res {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }
