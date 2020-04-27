package main

import (
	crdt "CRDTsGO"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

var dbPath string
var data string

//Structs and type
type replicationElement struct {
	Name        string
	ElementType string
}

type contentMap map[replicationElement]string

type replicationLayer struct {
	dfs    *Dfs
	or     *crdt.ORSet
	cmap   contentMap
}

//initalisation
func newReplicationLayer(id int,DB bool) *replicationLayer {
	dbPath = "./src/DFS/data.db"
	data = "data"
	dic,or :=readDB(id)



	l := replicationLayer{
		dfs:  new(Dfs),
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
		var el, u interface{}
		el = replicationElement{Name: msg.path, ElementType: msg.fileType}
		if msg.op == "add" {
			u = l.or.SrcAdd(el)
		} else if msg.op == "rm" {
			u = l.or.SrcRemove(el)
		}
		// rmsg:=RemoteMsg{SenderID:-1,Op:msg.op,Params:[]string{msg.path,msg.fileType},}
		rmsg := RemoteMsg{SenderID: -1, Op: msg.op,P1:el,P2:u}
		// fmt.Println(u)
		send <- rmsg
		go l.dfs.sendRemote(rmsg) //TODO: might miss local execution
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
		// fmt.Println(opMsg)
		if opMsg.Op == "add" {
			// el = opMsg.Params[0]
			// u = opMsg.Params[1]
			el=opMsg.P1
			u=opMsg.P2
			l.add(el.(replicationElement), u.(string))
		} else if opMsg.Op == "rm" {
			// el = opMsg.Params[0]
			// u = opMsg.Params[1]
			el=opMsg.P1
			u=opMsg.P2
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
func readDB(id int) (contentMap, *crdt.ORSet) {
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
	return contentMap, or
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
			token =GetTokens(l.or,k)
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

//GetTokens get set of tokens associated with specific element in the OR.set
func GetTokens(or *crdt.ORSet, el interface{}) string {
	items := []string{}

	it := or.Get(el).Iterator()
	for it.Next() {
		items = append(items, fmt.Sprintf("%v", it.Value()))
	}

	return strings.Join(items, ", ")
}

