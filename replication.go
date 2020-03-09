package main

import (
	"database/sql"
	"fmt"

	set "github.com/emirpasic/gods/sets/linkedhashset"
	_ "github.com/mattn/go-sqlite3"
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

var dbPath string
var data string

//Structs and type
type replicationElement struct {
	name        string
	elementType string
}

// type elementSet []*replicationElement
type elementSet *set.Set
type contentMap map[*replicationElement]string

type replicationLayer struct {
	dfs  *Dfs
	set  elementSet
	cmap contentMap
}

//initalisation
func newReplicationLayer() *replicationLayer {
	dbPath = "./src/DFS/data.db"
	data = "data"
	s, dic := readDB()

	// s := set.New()
	// s.Add(el)
	// dic := make(map[*replicationElement]string)
	// dic[&el] = ""

	l := replicationLayer{
		dfs:  new(Dfs),
		set:  s,
		cmap: dic,
	}

	return &l
}
func (l *replicationLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

func (l *replicationLayer) runLocally(send chan map[*replicationElement]string, recieve chan HierToRep) {
	send <- l.returnCurrentSet() //send the initial state
	for {
		msg := <-recieve
		if msg.op == "add" {
			l.add(msg.path, msg.fileType)
		} else if msg.op == "rm" {
			l.remove(msg.path, msg.fileType)
		}

		send <- l.returnCurrentSet() //send the updated set to hier
	}
}

//update inteface

func (l *replicationLayer) add(path string, typ string) {
	el := replicationElement{name: path, elementType: typ}
	// l.set = append(l.set, &el)
	(*l.set).Add(el) //element get added
	l.cmap[&el] = "" //initate with an empty content
	// l.updateDfs()
	fmt.Println("added", path)
}

func (l *replicationLayer) remove(path string, typ string) {
	//remove an element from the slice
	// temp := set.New()
	for _, i := range (*l.set).Values() {
		ii := i.(replicationElement)
		if ii.name == path && ii.elementType == typ {
			(*l.set).Remove(ii)
		}
	}
	// l.set = temp
	// fmt.Println((*l.set).Size(), temp.Size())
	// l.updateDfs()
	fmt.Println("removed", path)
}

// func (l *replicationLayer) udpate(path string,typ string){
// 	fmt.Println("element has been added")
// }

//update hier by through dfs
func (l *replicationLayer) updateDfs() {

	l.dfs.updateHier(l.returnCurrentSet()) //select only one the exist in the setS
}

func (l *replicationLayer) returnCurrentSet() map[*replicationElement]string {
	temp := make(map[*replicationElement]string)
	for _, k := range (*l.set).Values() {
		kk := (k.(replicationElement))
		temp[&kk] = l.cmap[&kk]
	}
	return temp
}

func (l *replicationLayer) printCurrentState() {
	fmt.Println("\nCRDT_Set\n-------------")
	// for _, k := range l.set {
	// 	v := l.cmap[k]
	// 	fmt.Println("", k.name, "content", v)
	// }
	for _, k := range (*l.set).Values() {
		kk := (k.(replicationElement))
		v := l.cmap[&kk]
		fmt.Println("", kk.name, "content", v)
	}
	fmt.Println()
}

//read the databse
func readDB() (*set.Set, contentMap) {
	s := set.New()
	contentMap := make(map[*replicationElement]string)

	database, err := sql.Open("sqlite3", dbPath)
	checkErr(err)
	rows, err := database.Query("SELECT path,type,content,used from " + data)
	checkErr(err)
	var path string
	var elementType string
	var content string
	var used int
	for rows.Next() {
		rows.Scan(&path, &elementType, &content, &used)
		el := replicationElement{name: path, elementType: elementType}
		contentMap[&el] = content
		if used == 1 {
			s.Add(el)
		}

	}
	rows.Close()
	return s, contentMap
}

func (l *replicationLayer) writeDB() {
	database, _ := sql.Open("sqlite3", dbPath)
	//dropping the table
	statement, err := database.Prepare("Drop table " + data)
	statement.Exec()
	checkErr(err)
	//creating data table
	statement, err = database.Prepare("create table " + data + " (id INTEGER PRIMARY KEY, path TEXT ,type TEXT ,content TEXT,used INTEGER)")
	statement.Exec()

	statement, err = database.Prepare("INSERT INTO " + data + " (path, type,content,used) VALUES (?,?,?,?)")

	//insert data points
	for k, v := range l.cmap {
		used := 0
		if (*l.set).Contains(*k) {
			used = 1
		}
		checkErr(err)
		statement.Exec(k.name, k.elementType, v, used)

	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
