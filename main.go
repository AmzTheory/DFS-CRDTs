package main

import (
	"fmt"
)

func main() {
	//actual DFS running
	ls := []replicationElement{
		replicationElement{name: "A", elementType: "t1"},
		replicationElement{name: "A/B", elementType: "t2"},
	}

	for _, el := range ls {
		fmt.Println("%s", el.elementType)
	}
}
