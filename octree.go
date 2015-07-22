package goctree

type Data interface {
	GetPosition() Vector3D
}

type Tree interface {
	Insert(...Data)
	Walk(bmin, bmax Vector3D, fn func(Data) WalkChoice) (seen int, choice WalkChoice)
	Size() int
}

type root struct {
	tree     *node
	size     int
	maxItem  int
	insertFn func(...Data)
	walkFn   func(bmin, bmax Vector3D, fn func(Data) WalkChoice) (seen int, choice WalkChoice)
}

type node struct {
	// Physical position/size. This implicitly defines the bounding
	// box of this node
	origin        Vector3D // The physical center of this node
	halfDimension Vector3D // Half the width/height/depth of this node

	/*
		Children follow a predictable pattern to make accesses simple.
		Here, - means less than 'origin' in that dimension, + means greater than.
		child:	0 1 2 3 4 5 6 7
		x:      - - - - + + + +
		y:      - - + + - - + +
		z:      - + - + - + - +
	*/
	children [8]*node
	datas    []Data
}

func NewRecursive(origin Vector3D, halfDimension Vector3D, maxItem int) Tree {
	r := &root{
		maxItem: maxItem,
		tree:    new(origin, halfDimension, maxItem),
	}
	r.insertFn = r.InsertRecursive
	r.walkFn = r.GetPointsInsideBoxRecursive
	return r
}
func NewIterative(origin Vector3D, halfDimension Vector3D, maxItem int) Tree {
	r := &root{
		maxItem: maxItem,
		tree:    new(origin, halfDimension, maxItem),
	}
	r.insertFn = r.InsertIterative
	r.walkFn = r.GetPointsInsideBoxIterative
	return r
}

func (r *root) Insert(pts ...Data) {
	r.insertFn(pts...)
}
func (r *root) Walk(bmin, bmax Vector3D, fn func(Data) WalkChoice) (seen int, choice WalkChoice) {
	return r.walkFn(bmin, bmax, fn)
}

func new(origin Vector3D, halfDimension Vector3D, items int) *node {
	return &node{
		origin:        origin,
		halfDimension: halfDimension,
		datas:         make([]Data, 0, items),
	}
}

func (o *node) GetOctantContainingPoint(point Vector3D) (oct int) {
	if !o.contains(point) {
		panic("Point out of tree: " + point.String())
	}
	if point[x] >= o.origin[x] {
		oct |= 4
	}
	if point[y] >= o.origin[y] {
		oct |= 2
	}
	if point[z] >= o.origin[z] {
		oct |= 1
	}
	return
}

func (n *node) findChildrenContainingPoint(point Vector3D) *node {
	return n.children[n.GetOctantContainingPoint(point)]
}

func (o *root) findLeafContainingPoint(node *node, point Vector3D) *node {
	for !node.IsLeafNode() {
		node = node.findChildrenContainingPoint(point)
	}
	return node
}

func (o *node) IsLeafNode() bool {
	// We are a leaf if we have no children. Since we either have none, or
	// all eight, it is sufficient to just check the first.
	return o.children[0] == nil
}

func side(b bool) float64 {
	if b {
		return 0.5
	}
	return -0.5
}

func (o *root) Size() int {
	return o.size
}

func (o *root) InsertRecursive(points ...Data) {
	for _, point := range points {
		o.size++
		if !o.tree.contains(point.GetPosition()) {
			panic("Point out of tree: " + point.GetPosition().String() + " " + o.tree.origin.String() + " " + o.tree.halfDimension.String())
		}
		o.insertRecursive(o.tree, point)
	}
}

func (o *root) insertRecursive(node *node, point Data) {
	leaf := o.findLeafContainingPoint(node, point.GetPosition())

	if len(leaf.datas) < o.maxItem {
		leaf.datas = append(leaf.datas, point)
		return
	} else {

		oldDatas := leaf.datas
		leaf.datas = leaf.datas[:0]

		o.Split(leaf)

		for _, oldData := range oldDatas {
			o.insertRecursive(leaf.findChildrenContainingPoint(oldData.GetPosition()), oldData)
		}
		o.insertRecursive(leaf.findChildrenContainingPoint(point.GetPosition()), point)
	}
}

// Split a leaf node and create new empty trees for each
// child octant.
func (o *root) Split(leaf *node) {
	for i := 0; i < 8; i++ {
		// Compute new bounding box for this child
		newOrigin := leaf.origin
		newOrigin[x] += leaf.halfDimension[x] * side(i&4 != 0)
		newOrigin[y] += leaf.halfDimension[y] * side(i&2 != 0)
		newOrigin[z] += leaf.halfDimension[z] * side(i&1 != 0)
		leaf.children[i] = new(newOrigin, leaf.halfDimension.Imul(0.5), o.maxItem)
	}
}

func (o *node) contains(p Vector3D) bool {
	return !(p[x] > o.origin[x]+o.halfDimension[x] ||
		p[x] < o.origin[x]-o.halfDimension[x] ||

		p[y] > o.origin[y]+o.halfDimension[y] ||
		p[y] < o.origin[y]-o.halfDimension[y] ||

		p[z] > o.origin[z]+o.halfDimension[z] ||
		p[z] < o.origin[z]-o.halfDimension[z])
}

type insertPair struct {
	node   *node
	points []Data
}

func (o *root) InsertIterative(datas ...Data) {
	o.size++

	for len(datas) > 0 {
		data := datas[len(datas)-1]
		datas = datas[:len(datas)-1]
		leaf := o.findLeafContainingPoint(o.tree, data.GetPosition())
		if len(leaf.datas) < o.maxItem {
			leaf.datas = append(leaf.datas, data)
		} else {

			oldDatas := leaf.datas
			leaf.datas = nil

			o.Split(leaf)
			datas = append(datas, oldDatas...)
			datas = append(datas, data)
		}
	}
}
