package main

/*

testing
---------

skip policy

*/
import(
	"testing"
	"reflect"
)


func TestBuildTree(t *testing.T){
//just a root
//no orphans
//with orphans
//complicated case(so many elements)
	// a:=newRE("/","dir")

	
	c:=newRE("/c/","dir")
	b:=newRE("/b","txt")

	var tests = []struct {
		cmap 	map[*RepElem]string
		want    DfsNode
		}{
			{map[*RepElem]string{}, 
			 DfsNode{name: "/", fileType: "dir", path: "", content: "",parent:nil,children:map[string]*DfsNode{}}},
			 {map[*RepElem]string{&c:"",&b:""}, 
			 DfsNode{name: "/", fileType: "dir", path: "", content: "",parent:nil,children:map[string]*DfsNode{}}},
		}

	
		var got DfsNode
		for _,test:=range tests{
			got=*buildTree(test.cmap)
			

			if !reflect.DeepEqual(test.want,got){
				t.Errorf("buildTree()  got= %+v, want %+v",got,test.want)
			}
		}



}