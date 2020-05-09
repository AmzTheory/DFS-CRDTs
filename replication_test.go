package main



/*

test replication layer
----------------------
adding element   (dic & OR)

removing element  (dic & OR)


reading & writing from DB


*/

import (
	"testing"
	// "reflect"
)

func TestAdd(t *testing.T){
	rep:=newRepLayer(1)


	var tests = []struct {
		name string
		ty 	 string
		token string
		wantDic []interface{}
		wantSet []interface{}
		}{
			{"/","dir","1", []interface{}{newRE("/","dir")},[]interface{}{newRE("/","dir")}}, 
			{"/a", "txt","2",[]interface{}{newRE("/","dir"),newRE("/a","txt")},[]interface{}{newRE("/","dir"),newRE("/a","txt")}},
			{"/a", "txt","3",[]interface{}{newRE("/","dir"),newRE("/a","txt")},[]interface{}{newRE("/","dir"),newRE("/a","txt")}},  //duplicate elements
			{"/c/", "dir","3",[]interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/c/","dir")},[]interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/c/","dir")}},  //add dir
			
		}
		var el RepElem
		for _, test := range tests {
			el=newRE(test.name,test.ty)
			
			//modify state
			rep.add(el,test.token)


			
			//retrieve update state
			gotSet:=rep.or.Values()
			gotDic:=getMapKeys(rep.cmap)
			
			if len(difference(gotSet,test.wantSet))>0{
				t.Errorf("add(%s,%s)  Set= %s, want %s", el.name,el.elemType, gotSet,test.wantSet)
			}

			if len(difference(gotDic,test.wantDic))>0{
				t.Errorf("add(%s,%s)  Dic= %s, want %s", el.name,el.elemType, gotDic,test.wantDic)
			}
		}
}

func TestRemove(t *testing.T){
	rep:=newRepLayer(1)
	
	rep.add(newRE("/","dir"),"1")
	rep.add(newRE("/a","txt"),"2")
	rep.add(newRE("/a","txt"),"3")
	rep.add(newRE("/d/","dir"),"4")

	var tests = []struct {
		name string
		ty 	 string
		r []interface{}
		wantSet []interface{}
		wantDic []interface{}
		}{
			{"/d/","dir",[]interface{}{"4"}, []interface{}{newRE("/","dir"),newRE("/a","txt")}, []interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/d/","dir")}},  //remove directory
			{"/a","dir",[]interface{}{"2","3"}, []interface{}{newRE("/","dir"),newRE("/a","txt")}, []interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/d/","dir")}},  //remove an element (same name but differen type)
			{"/a","txt",[]interface{}{"2","3"}, []interface{}{newRE("/","dir")}, []interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/d/","dir")}},  //remove and element added twice
			{"/c","dir",[]interface{}{"2","3"}, []interface{}{newRE("/","dir")}, []interface{}{newRE("/","dir"),newRE("/a","txt"),newRE("/d/","dir")}},  //remove and element that doesnt exist
		} 
		var el RepElem
		for _, test := range tests {
			el=newRE(test.name,test.ty)
			
			//modify state
			rep.remove(test.r,el)


			
			//retrieve update state
			gotSet:=rep.or.Values()
			gotDic:=getMapKeys(rep.cmap)
			
			if len(difference(gotSet,test.wantSet))>0{
				t.Errorf("remove((%s,%s),%s)  Set= %s, want %s", el.name,el.elemType,test.r, gotSet,test.wantSet)
			}

			if len(difference(gotDic,test.wantDic))>0{
				t.Errorf("remove((%s,%s),%s)  dic= %s, want %s", el.name,el.elemType,test.r, gotDic,test.wantDic)
			}
		}
}






//helperFunctions
func getMapKeys(mymap map[RepElem]string) []interface{}{
	keys := []interface{}{}
	for k := range mymap {
    	keys = append(keys, k)
	}
	return keys
}

//
func newRE(name,typ string) RepElem{
	return RepElem{name:name,elemType:typ}
}

func difference(slice1 []interface{}, slice2 []interface{}) ([]interface{}){
    diffStr := []interface{}{}
    m :=map [interface{}]int{}

    for _, s1Val := range slice1 {
        m[s1Val] = 1
    }
    for _, s2Val := range slice2 {
        m[s2Val] = m[s2Val] + 1
    }

    for mKey, mVal := range m {
        if mVal==1 {
            diffStr = append(diffStr, mKey)
        }
    }

    return diffStr
}