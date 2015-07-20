package goctree

type WalkChoice bool

const (
	ContinueWalking = WalkChoice(true)
	StopWalking     = WalkChoice(false)
)

// GetPointsInsideBox is a really simple routine for querying the tree for points
// within a bounding box defined by min/max points (bmin, bmax)
// All results are pushed into 'results'
func (o *Node) GetPointsInsideBox(bmin, bmax Vector3D, fn func(Data) WalkChoice) {
	o.getPointsInsideBox(bmin, bmax, fn)
}

func (o *Node) getPointsInsideBox(bmin, bmax Vector3D, fn func(Data) WalkChoice) WalkChoice {
	// If we're at a leaf node, just see if the current data point is inside
	// the query bounding box
	if o.IsLeafNode() {
		if o.data != nil {
			p := o.data.Position()
			if p[x] > bmax[x] || p[y] > bmax[y] || p[z] > bmax[z] {
				return ContinueWalking
			}
			if p[x] < bmin[x] || p[y] < bmin[y] || p[z] < bmin[z] {
				return ContinueWalking
			}
			return fn(o.data)
		}
	} else {
		// We're at an interior node of the tree. We will check to see if
		// the query bounding box lies outside the octants of this node.
		for i := 0; i < 8; i++ {
			// Compute the min/max corners of this child octant
			cmax := o.children[i].origin.Add(o.children[i].halfDimension)
			cmin := o.children[i].origin.Sub(o.children[i].halfDimension)

			// If the query rectangle is outside the child's bounding box,
			// then continue
			if cmax[x] < bmin[x] || cmax[y] < bmin[y] || cmax[z] < bmin[z] {
				continue
			}
			if cmin[x] > bmax[x] || cmin[y] > bmax[y] || cmin[z] > bmax[z] {
				continue
			}

			// At this point, we've determined that this child is intersecting
			// the query bounding box
			if o.children[i].getPointsInsideBox(bmin, bmax, fn) == StopWalking {
				return StopWalking
			}
		}
	}

	return ContinueWalking
}
