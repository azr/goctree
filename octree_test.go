package goctree_test

import (
	"testing"

	octree "github.com/azr/goctree"
)

func TestInsertGet(t *testing.T) {
	ot := octree.New(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2})

	backleftbottom := &TestPoint{
		p:    octree.Vector3D{-3, -2, -2},
		name: "backleftbottom",
	}
	ot.Insert(backleftbottom)

	frontrighttop := &TestPoint{
		p:    octree.Vector3D{3, 2, 2},
		name: "frontrighttop",
	}
	ot.Insert(frontrighttop)

	{
		times := 0
		ot.GetPointsInsideBox(octree.Vector3D{-3, -2, -2}, octree.Vector3D{0, 0, 0}, func(d octree.Data) octree.WalkChoice {
			times++
			if d != backleftbottom {
				t.Error("Item is different")
			}
			return octree.ContinueWalking
		})
		if times != 1 {
			t.Errorf("Incorrect number of found items: %d", times)
		}
	}

	{
		times := 0
		ot.GetPointsInsideBox(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2}, func(d octree.Data) octree.WalkChoice {
			times++
			if d != frontrighttop {
				t.Error("Item is different")
			}
			return octree.ContinueWalking
		})
		if times != 1 {
			t.Errorf("Incorrect number of found items: %d", times)
		}
	}
}

type TestPoint struct {
	p    octree.Vector3D
	name string
}

func (t *TestPoint) Position() octree.Vector3D {
	return t.p
}
