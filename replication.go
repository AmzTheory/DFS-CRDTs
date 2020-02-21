package main

import (
	"fmt"
)


/*
fields	
	SET (assume its and ORSET) string
	MAP (assume its CRDT)  String->String
functions
	add(p,t)
	remove(p,t)
	update(p,t,u)

returns to upper layer
	MAP   (p,t)->String (content)
*/



//Structs and type
type replicationElement struct{
	name 		string
	elementType string
}
type elementSet []replicationElement
type contentMap map[string]string

type replicationLayer struct{
	dfs  Dfs
	set elementSet
	cmap contentMap
}


//update inteface

func (l replicationLayer) add(path string,typ string){
	fmt.Println("element has been added")
}

func (l replicationLayer) remove(path string,typ string){
	fmt.Println("element has been removed")
}

// func (l replicationLayer) udpate(path string,typ string){
// 	fmt.Println("element has been added")
// }


//update hier by through dfs
func updateHier(){
		//iclude only keys that exist in the set
		//pass the map
}
