package goctree

type Data interface {
	GetPosition() Vector3D
}

type Node struct {
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
	children [8]*Node
	data     Data
	root     bool
	size     int
}

func New(origin Vector3D, halfDimension Vector3D) *Node {
	o := new(origin, halfDimension)
	o.root = true
	return o
}

func new(origin Vector3D, halfDimension Vector3D) *Node {
	return &Node{
		origin:        origin,
		halfDimension: halfDimension,
	}
}

func (o *Node) GetOctantContainingPoint(point Vector3D) (oct int) {
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

func (o *Node) IsRootNode() bool {
	return o.root
}

func (o *Node) IsLeafNode() bool {
	// This is correct, but overkill. See below.
	// for i := 0; i < 8; i++ {
	// 	if children[i] != nil {
	// 		return false
	// 	}
	// }
	// return true

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

func (o *Node) Size() int {
	return o.size
}

func (o *Node) Insert(point Data) {
	if o.IsRootNode() {
		o.size++
	}

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
			oldPoint := o.data
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
			o.children[o.GetOctantContainingPoint(oldPoint.GetPosition())].Insert(oldPoint)
			o.children[o.GetOctantContainingPoint(point.GetPosition())].Insert(point)
		}
	} else {
		// We are at an interior node. Insert recursively into the
		// appropriate child octant
		o.children[o.GetOctantContainingPoint(point.GetPosition())].Insert(point)
	}
}
