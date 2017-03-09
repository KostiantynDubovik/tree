package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

const selectQuery = "SELECT `node_id`, `image_name`, `level` FROM my_tree ORDER BY left_key"
const insertQuery = "CALL `images_tree`.`insert_node`(?, ?, ?)"
const deleteQuery = "CALL `images_tree`.`delete_node`(?)"
const deleteAllQuery = "CALL `images_tree`.`delete_all_nodes`()"
const driverName = "mysql"
const dataSourceName = "root:root:@/images_tree"




func GetNodes() {

	database, err := sql.Open(driverName, dataSourceName)
	defer database.Close()
	checkErr(err)
	rows, err := database.Query(selectQuery)
	checkErr(err)

	for rows.Next() {
		var node Node
		var nodeId int
		var leftKey int
		var rightKey int
		var level int
		var parentId int
		var imageName string
		err = rows.Scan(&nodeId, &leftKey, &rightKey, &level, &parentId, &imageName)
		checkErr(err)
		node.NodeId = nodeId
		node.PrentId = parentId
		node.ImageName = imageName
		Nodes = append(Nodes, node)
	}
}

func AddNode(nodeId int, parentId int, imageName string) {
	database, err := sql.Open(driverName, dataSourceName)
	defer database.Close()
	checkErr(err)
	database.Exec(insertQuery, parentId, imageName)
	GetNodes()
}

func DeleteNode(nodeId int) {
	database, err := sql.Open(driverName, dataSourceName)
	defer database.Close()
	checkErr(err)
	database.Exec(deleteQuery, nodeId)
	GetNodes()
}

func DeleteAllNodes()  {
	database, err := sql.Open(driverName, dataSourceName)
	defer database.Close()
	checkErr(err)
	database.Exec(deleteAllQuery)
	GetNodes()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
