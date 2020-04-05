package main

import (
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"os"
)

/*

create number of DFS    (number of replicas)    

establish connections among replicas

manual testing (switch between replicas)

Perform operations   (randomly generate operations perform testting)

Testing convergence    when all replicas have recieved 


*/


type DfsController struct{
	replicas map[int]*Dfs //map of clients  ports/ids
	input    map[int]chan bool
}


func NewController() *DfsController{
	return &DfsController{replicas:make(map[int]*Dfs),input:make(map[int]chan bool),}
}

func (inst *DfsController) create(rep map[int][]int){
	for k,v:=range rep{
		inst.replicas[k]=newDfs(k,v)
	}
}
func (inst *DfsController) SetUpConnection(){
	for k, v := range inst.replicas { 		
		v.start()
		ch:=make(chan bool) //request the access mode
		inst.input[k]=ch	
		b := make(chan bool)
		go v.runAll(b,ch)

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

func (inst *DfsController) run(){
	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Print("type the replica you wish to access\n")
		text, _ := reader.ReadString('\n')

		text = strings.Replace(text, "\n", "", 1)
		replica,_:=strconv.Atoi(text)

		inst.input[replica]<-true  //start access mode

		<-inst.input[replica]  //wait for quit request


	}
}