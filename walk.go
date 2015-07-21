package goctree

type WalkChoice bool

const (
	ContinueWalking = WalkChoice(false)
	StopWalking     = WalkChoice(true)
)

//GetPointsInsideBox lists breadth first points
func (o *root) GetPointsInsideBoxIterative(bmin, bmax Vector3D, fn func(Data) WalkChoice) (seen int, choice WalkChoice) {
	return o.tree.getPointsInsideBoxIterative(bmin, bmax, fn)
}

func (o *node) getPointsInsideBoxIterative(bmin, bmax Vector3D, fn func(Data) WalkChoice) (seen int, choice WalkChoice) {
	nodes := make([]*node, 1, 50)
	nodes[0] = o

	for len(nodes) != 0 && choice == ContinueWalking {
		o := nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
		seen++
		// If we're at a leaf node, just see if the current data point is inside
		// the query bounding box
		if o.IsLeafNode() {
			if o.data != nil {
				p := o.data.GetPosition()
				if p[x] > bmax[x] || p[y] > bmax[y] || p[z] > bmax[z] {
					continue
				}
				if p[x] < bmin[x] || p[y] < bmin[y] || p[z] < bmin[z] {
					continue
				}
				choice = fn(o.data)
			}
		} else {
			// We're at an interior node of the tree. We will check to see if
			// the query bounding box lies outside the octants of this node.
			for i := 0; i < 8 && choice == ContinueWalking; i++ {
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
				nodes = append(nodes, o.children[i])
			}
		}
	}
	return
}

// GetPointsInsideBox is a really simple routine for querying the tree for points
// within a bounding box defined by min/max points (bmin, bmax)
// sent to fn
// Beware, recursivity means stack overflow and depth first !
func (o *root) GetPointsInsideBoxRecursive(bmin, bmax Vector3D, fn func(Data) WalkChoice) (int, WalkChoice) {
	return o.tree.getPointsInsideBoxRecursive(bmin, bmax, fn)
}

func (o *node) getPointsInsideBoxRecursive(bmin, bmax Vector3D, fn func(Data) WalkChoice) (int, WalkChoice) {
	// If we're at a leaf node, just see if the current data point is inside
	// the query bounding box
	if o.IsLeafNode() {
		if o.data != nil {
			p := o.data.GetPosition()
			if p[x] > bmax[x] || p[y] > bmax[y] || p[z] > bmax[z] {
				return 1, ContinueWalking
			}
			if p[x] < bmin[x] || p[y] < bmin[y] || p[z] < bmin[z] {
				return 1, ContinueWalking
			}
			return 1, fn(o.data)
		}
		return 1, ContinueWalking
	} else {
		walked := 1
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
			n, c := o.children[i].getPointsInsideBoxRecursive(bmin, bmax, fn)
			walked += n
			if c == StopWalking {
				return walked, StopWalking
			}
		}
		return walked, ContinueWalking
	}
}
