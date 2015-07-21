package goctree_test

import (
	"testing"

	octree "github.com/azr/goctree"
)

func testQueryReturns(t *testing.T, ot octree.Tree, bmin, bmax octree.Vector3D, datas ...*TestPoint) {
	times := 0
	ot.Walk(bmin, bmax, func(d octree.Data) octree.WalkChoice {
		if datas[times] != d {
			t.Fatalf("Item is different expected %v, got %v", datas[times], d)
		}
		times++
		return octree.ContinueWalking
	})
	if len(datas) != times {
		t.Fatalf("got: %d items expected: %v", times, datas)
	}
}

type testQuery struct {
	bmin, bmax octree.Vector3D
	returns    []*TestPoint
}

type testStruct struct {
	tree    octree.Tree
	tps     []*TestPoint
	size    int
	queries []testQuery
}

func tests(t *testing.T, ts []testStruct) {
	for _, test := range ts {
		for _, point := range test.tps {
			test.tree.Insert(point)
		}
		if test.tree.Size() != test.size {
			t.Error("Sizes differ: expected %d, got %d", test.size, test.tree.Size())
		}
		for _, query := range test.queries {
			testQueryReturns(t, test.tree, query.bmin, query.bmax, query.returns...)
		}
	}
}

func TestInsertGetAndSize(t *testing.T) {
	backleftbottom := &TestPoint{
		p:    octree.Vector3D{-3, -2, -2},
		name: "backleftbottom",
	}
	frontrighttop := &TestPoint{
		p:    octree.Vector3D{3, 2, 2},
		name: "frontrighttop",
	}
	center := &TestPoint{
		p:    octree.Vector3D{0, 0, 0},
		name: "center",
	}

	// ocr :=

	ts := []testStruct{
		{
			tree: octree.NewRecursive(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2}),
			tps:  []*TestPoint{backleftbottom, frontrighttop, center},
			size: 3,
			queries: []testQuery{
				{
					bmin:    octree.Vector3D{-3, -2, -2},
					bmax:    octree.Vector3D{0, 0, 0},
					returns: []*TestPoint{backleftbottom, center},
				},
				{
					bmin:    octree.Vector3D{0, 0, 0},
					bmax:    octree.Vector3D{3, 2, 2},
					returns: []*TestPoint{center, frontrighttop},
				},
			},
		},
		{
			tree: octree.NewIterative(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2}),
			tps:  []*TestPoint{backleftbottom, frontrighttop, center},
			size: 3,
			queries: []testQuery{
				{
					bmin:    octree.Vector3D{-3, -2, -2},
					bmax:    octree.Vector3D{0, 0, 0},
					returns: []*TestPoint{center, backleftbottom},
				},
				{
					bmin:    octree.Vector3D{0, 0, 0},
					bmax:    octree.Vector3D{3, 2, 2},
					returns: []*TestPoint{frontrighttop, center},
				},
			},
		},
	}

	tests(t, ts)
}

type TestPoint struct {
	p    octree.Vector3D
	name string
}

func (t *TestPoint) GetPosition() octree.Vector3D {
	return t.p
}
