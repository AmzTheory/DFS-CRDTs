package main

import (
	// set "github.com/emirpasic/gods/sets/linkedhashset"
)

/*

create number of DFS    (number of replicas)    

establish connections among replicas

Perform operations   (randomly generate operations perform testting)

Testing convergence    when all replicas have recieved 


*/


type DfsController struct{
	replicas map[int]*Dfs //map of clients  ports/ids
}


func NewController() *DfsController{
	return &DfsController{replicas:make(map[int]*Dfs)}
}

func (inst *DfsController) create(rep map[int][]int){
	for k,v:=range rep{
		inst.replicas[k]=newDfs(k,v)
	}
}
func (inst *DfsController) SetUpConnection(){
	for _, v := range inst.replicas { 		
		v.start()	
		b := make(chan bool)
		go v.runAll(b)

		//indicate the dfs listener is open
		<-b
	}
}
func (inst *DfsController) SetUpCommunication(){
	for _, v := range inst.replicas { 		
		v.startConnecting()
	}
}

func generateReplicasMap(n int) (map[int][]int) { 
	//n^2
	repMap:=make(map[int][]int)
	init:=1001
	
	for i:=init;i<=init+n;i++{
		repMap[i]=[]int{}
		for j:=1001;j<=init+n;j++{
			if(j!=i){repMap[i]=append(repMap[i],j)}
		}
	}

	return repMap
}



