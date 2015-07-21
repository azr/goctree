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
	data     Data
}

func NewRecursive(origin Vector3D, halfDimension Vector3D) Tree {
	r := &root{
		tree: new(origin, halfDimension),
	}
	r.insertFn = r.InsertRecursive
	r.walkFn = r.GetPointsInsideBoxRecursive
	return r
}
func NewIterative(origin Vector3D, halfDimension Vector3D) Tree {
	r := &root{
		tree: new(origin, halfDimension),
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

func new(origin Vector3D, halfDimension Vector3D) *node {
	return &node{
		origin:        origin,
		halfDimension: halfDimension,
	}
}

func (o *node) GetOctantContainingPoint(point Vector3D) (oct int) {
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
		if !o.tree.contains(point) {
			panic("Point out of tree: " + point.GetPosition().String() + " " + o.tree.origin.String() + " " + o.tree.halfDimension.String())
		}
		o.tree.InsertRecursive(point)
	}
}

func (o *node) InsertRecursive(point Data) {
	// If this node doesn't have a data point yet assigned
	// and it is a leaf, then we're done!
	if o.IsLeafNode() {
		if o.data == nil {
			o.data = point
			return
		} else {
			// We're at a leaf, but there's already something here
			// We will split this node so that it has 8 child octants
			// and then insert the old data that was here, along with
			// this new data point

			// Save this data point that was here for a later re-insert
			oldData := o.data
			o.data = nil

			// Split the current node and create new empty trees for each
			// child octant.
			for i := 0; i < 8; i++ {
				// Compute new bounding box for this child
				newOrigin := o.origin
				newOrigin[x] += o.halfDimension[x] * side(i&4 != 0)
				newOrigin[y] += o.halfDimension[y] * side(i&2 != 0)
				newOrigin[z] += o.halfDimension[z] * side(i&1 != 0)
				o.children[i] = new(newOrigin, o.halfDimension.Imul(0.5))
			}

			// Re-insert the old point, and insert this new point
			// (We wouldn't need to insert from the root, because we already
			// know it's guaranteed to be in this section of the tree)
			o.children[o.GetOctantContainingPoint(oldData.GetPosition())].InsertRecursive(oldData)
			o.children[o.GetOctantContainingPoint(point.GetPosition())].InsertRecursive(point)
		}
	} else {
		// We are at an interior node. Insert recursively into the
		// appropriate child octant
		o.children[o.GetOctantContainingPoint(point.GetPosition())].InsertRecursive(point)
	}
}

func (o *node) contains(point Data) bool {
	p := point.GetPosition()
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

func (o *root) InsertIterative(points ...Data) {
	o.size++

	pairs := []insertPair{
		{
			points: points,
			node:   o.tree,
		},
	}

	for len(pairs) > 0 {
		pair := pairs[len(pairs)-1]
		pairs = pairs[:len(pairs)-1]

		node := pair.node
		for len(pair.points) > 0 {
			point := pair.points[len(pair.points)-1]
			pair.points = pair.points[:len(pair.points)-1]

			if node.IsLeafNode() {
				if node.data == nil {
					node.data = point
					continue
				} else {
					// We're at a leaf, but there's already something here
					// We will split this node so that it has 8 child octants
					// and then insert the old data that was here, along with
					// this new data point

					// Save this data point that was here for a later re-insert
					oldData := node.data
					node.data = nil

					// Split the current node and create new empty trees for each
					// child octant.
					for i := 0; i < 8; i++ {
						// Compute new bounding box for this child
						newOrigin := node.origin
						newOrigin[x] += node.halfDimension[x] * side(i&4 != 0)
						newOrigin[y] += node.halfDimension[y] * side(i&2 != 0)
						newOrigin[z] += node.halfDimension[z] * side(i&1 != 0)
						node.children[i] = new(newOrigin, node.halfDimension.Imul(0.5))
					}

					// Re-insert the old point, and insert this new point
					// (We wouldn't need to insert from the root, because we already
					// know it's guaranteed to be in this section of the tree)
					pairs = append(pairs,
						insertPair{
							node:   node.children[node.GetOctantContainingPoint(oldData.GetPosition())],
							points: []Data{oldData},
						})
					pairs = append(pairs,
						insertPair{
							node:   node.children[node.GetOctantContainingPoint(point.GetPosition())],
							points: []Data{point},
						})
				}
			} else {
				pairs = append(pairs,
					insertPair{
						node:   node.children[node.GetOctantContainingPoint(point.GetPosition())],
						points: []Data{point},
					})
			}
		}
	}
}
