package main

import (
	"fmt"
)

/*
	instance of
		UserInterface
		Hier
		replication
	Model communication
		between layers		
		
	View of the DFS

	Comunication between replicas
		assign operations
*/


type Dfs stuct{
	hier hierLayer
	rep  replicationLayer
	//userInterface
}



//downwards

//User interface to Hier
func (Dfs dfs) updateAddHier(path string,n string,typ string){}
func (Dfs dfs) updateRemoveHier(path string,typ string){}
//update

//Hier to replication
func (Dfs dfs) updateAddReplication(path string,typ string){}
func (Dfs dfs) updateRemoveReplication(path string,typ string){}

//replication to other replicas (future)


//upwards

func (Dfs dfs) updateHier(cmap map[string]string){}
func (Dfs dfs) updateInterface(tree DfsTree){}