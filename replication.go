package main

import (
	crdt "CRDTsGO"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	set "github.com/emirpasic/gods/sets/linkedhashset"
	_ "github.com/mattn/go-sqlite3"
)

var dbPath string
var data string

//Structs and type
type replicationElement struct {
	Name        string
	ElementType string
}

// type elementSet []*replicationElement
type elementSet *set.Set
type contentMap map[replicationElement]string

type replicationLayer struct {
	dfs    *Dfs
	set    elementSet
	or     *crdt.ORSet
	cmap   contentMap
	opLock sync.Mutex
}

//initalisation
func newReplicationLayer(id int) *replicationLayer {
	dbPath = "./src/DFS/data.db"
	data = "data"
	s, dic, or := readDB(id)

	// s := set.New()
	// s.Add(el)
	// dic := make(map[*replicationElement]string)
	// dic[&el] = ""

	l := replicationLayer{
		dfs:  new(Dfs),
		set:  s,
		or:   or,
		cmap: dic,
	}

	return &l
}
func (l *replicationLayer) setDfs(dfs *Dfs) {
	l.dfs = dfs
}

//run locally & remotely comunicate theire messages to pushupState, which in turn get excuted by pushupstate, then passed to upper layers

func (l *replicationLayer) runLocally(send chan RemoteMsg, recieve chan HierToRep) {
	for {
		msg := <-recieve
		// if msg.op == "add" {
		// 	l.add(msg.path, msg.fileType)
		// } else if msg.op == "rm" {
		// 	l.remove(msg.path, msg.fileType)
		// }
		var el, u interface{}
		el = replicationElement{Name: msg.path, ElementType: msg.fileType}
		if msg.op == "add" {
			u = l.or.SrcAdd(el)
		} else if msg.op == "rm" {
			u = l.or.SrcRemove(el)
		}
		// rmsg:=RemoteMsg{SenderID:-1,Op:msg.op,Params:[]string{msg.path,msg.fileType},}
		rmsg := RemoteMsg{SenderID: -1, Op: msg.op, Params: []interface{}{el, u}}
		send <- rmsg
		go l.dfs.sendRemote(rmsg) //broadcast to others
		// send <- l.returnCurrentSet() //send the updated set to hier
	}
}

//execute operation and update hier
func (l *replicationLayer) pushUpState(send chan map[*replicationElement]string, recieve chan RemoteMsg) {
	send <- l.returnCurrentSet() //send the initial state
	var opMsg RemoteMsg
	var el, u interface{}
	var r []interface{}
	for { //wait for operation to be executed local/remotely
		opMsg = <-recieve
		if opMsg.Op == "add" {
			el = opMsg.Params[0]
			u = opMsg.Params[1]
			l.add(el.(replicationElement), u.(string))
		} else if opMsg.Op == "rm" {
			el = opMsg.Params[0]
			u = opMsg.Params[1]
			r = u.([]interface{})
			l.remove(r, el)
		}
		send <- l.returnCurrentSet()
	}
}

//listen remotely
func (l *replicationLayer) runRemotely(send chan RemoteMsg, recieve chan RemoteMsg) {
	var rmsg RemoteMsg
	for {
		rmsg = <-recieve

		send <- rmsg

		// fmt.Println("rep Recieved ", l.dfs.id)
	}
}

//update inteface

func (l *replicationLayer) add(el replicationElement, u string) {
	l.or.Add(u, el)
	if _, ok := l.cmap[el]; !ok { //start with an empty content in the case it hasn't be creatd
		l.cmap[el] = ""
	}

}

func (l *replicationLayer) remove(r []interface{}, el interface{}) {
	l.or.Remove(r, el)
}

//update hier by through dfs
func (l *replicationLayer) updateDfs() {
	l.dfs.updateHier(l.returnCurrentSet()) //select only one the exist in the setS
}

func (l *replicationLayer) returnCurrentSet() map[*replicationElement]string {
	temp := make(map[*replicationElement]string)
	for _, k := range l.or.Values() {
		kk := (k.(replicationElement))
		temp[&kk] = l.cmap[kk]
	}
	return temp
}

//compare current instance OR set with the orSET passed
func (l *replicationLayer) RepEqual(or *crdt.ORSet) bool {
	return l.or.Equal(or)
}

//
func (l *replicationLayer) printCurrentState() {
	fmt.Println("\nCRDT_Set\n-------------")

	for _, k := range l.or.Values() {
		kk := (k.(replicationElement))
		v := l.cmap[kk]
		fmt.Println("", kk.Name, "content", v)
	}
	fmt.Println()
}

//read the databse
func readDB(id int) (*set.Set, contentMap, *crdt.ORSet) {
	s := set.New()
	or := crdt.NewORSet()
	contentMap := make(map[replicationElement]string)

	database, err := sql.Open("sqlite3", dbPath)
	checkErr(err)

	rows, err := database.Query("SELECT path,type,content,used,token from " + data + " where dfsId=" + strconv.Itoa(id))
	checkErr(err)
	var path, content, elementType, token string
	var used int
	for rows.Next() {
		rows.Scan(&path, &elementType, &content, &used, &token)
		el := replicationElement{Name: path, ElementType: elementType}
		contentMap[el] = content
		if used == 1 {

			///Add for each token into OR
			for _, v := range strings.Split(token, ",") {
				or.Add(v, el)
			}
		}

	}
	rows.Close()
	return s, contentMap, or
}

func (l *replicationLayer) writeDB() {
	database, _ := sql.Open("sqlite3", dbPath)
	//delete dfs records
	id := (*l.dfs).id
	statement, err := database.Prepare("delete from " + data + " where dfsID=" + strconv.Itoa(id))
	statement.Exec()
	checkErr(err)
	// //creating data table
	// statement, err = database.Prepare("create table " + data + " (id INTEGER PRIMARY KEY, path TEXT ,type TEXT ,content TEXT,used INTEGER,dfsID Integer)")
	// statement.Exec()

	statement, err = database.Prepare("INSERT INTO " + data + " (path, type,content,used,dfsID,token) VALUES (?,?,?,?,?,?)")

	var used int
	var token string
	//insert data points
	for k, v := range l.cmap {
		used = 0
		token = ""
		if (*l.or).Contains(k) {
			used = 1
			token = l.or.GetTokens(k)
		}

		checkErr(err)
		statement.Exec(k.Name, k.ElementType, v, used, l.dfs.id, token)

	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
