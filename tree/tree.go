package tree

type node struct {
	id int
	leftKey int
	rightKey int
	prentId int
	level int
	originalImageName string
}

type tree struct {
	 nodes [] node
}
