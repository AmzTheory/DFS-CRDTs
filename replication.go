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
type set []replicationElement

type replication
