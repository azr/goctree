package goctree_test

import (
	"testing"

	octree "github.com/azr/goctree"
)

func testQueryReturns(t *testing.T, n, q int, ot octree.Tree, bmin, bmax octree.Vector3D, datas ...*TestPoint) {
	times := 0
	ot.Walk(bmin, bmax, func(d octree.Data) octree.WalkChoice {
		if datas[times] != d {
			t.Fatalf("[%d][%d]Item is different expected %v, got %v", n, q, datas[times], d)
		}
		times++
		return octree.ContinueWalking
	})
	if len(datas) != times {
		t.Errorf("[%d][%d]got: %d items expected: %v", n, q, times, datas)
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
	for n, test := range ts {
		for _, point := range test.tps {
			test.tree.Insert(point)
		}
		if test.tree.Size() != test.size {
			t.Errorf("Sizes differ: expected %d, got %d", test.size, test.tree.Size())
		}
		for q, query := range test.queries {
			testQueryReturns(t, n, q, test.tree, query.bmin, query.bmax, query.returns...)
		}
	}
}

func TestInsertGetAndSize(t *testing.T) {
	backleftbottom := &TestPoint{
		p:    octree.Vector3D{-3, -2, -2},
		name: "backleftbottom",
	}
	backcenterbottom := &TestPoint{
		p:    octree.Vector3D{-3, 0, -2},
		name: "backcenterbottom",
	}
	backrightbottom := &TestPoint{
		p:    octree.Vector3D{-3, 2, -2},
		name: "backrightbottom",
	}
	backlefttop := &TestPoint{
		p:    octree.Vector3D{-3, -2, 2},
		name: "backlefttop",
	}
	backcentertop := &TestPoint{
		p:    octree.Vector3D{-3, 0, 2},
		name: "backcentertop",
	}
	backrighttop := &TestPoint{
		p:    octree.Vector3D{-3, 2, 2},
		name: "backrighttop",
	}
	frontrighttop := &TestPoint{
		p:    octree.Vector3D{3, 2, 2},
		name: "frontrighttop",
	}
	frontcentertop := &TestPoint{
		p:    octree.Vector3D{3, 0, 2},
		name: "frontcentertop",
	}
	frontlefttop := &TestPoint{
		p:    octree.Vector3D{3, -2, 2},
		name: "frontlefttop",
	}
	center := &TestPoint{
		p:    octree.Vector3D{0, 0, 0},
		name: "center",
	}

	// ocr :=

	ts := []testStruct{
		{
			tree: octree.NewRecursive(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2}, 8),
			tps:  []*TestPoint{backleftbottom, frontrighttop, center, backcenterbottom, backrightbottom, backlefttop, backcentertop, backrighttop, frontcentertop, frontlefttop},
			size: 10,
			queries: []testQuery{
				{
					bmin:    octree.Vector3D{-3, -2, -2},
					bmax:    octree.Vector3D{0, 0, 0},
					returns: []*TestPoint{backleftbottom, backcenterbottom, center},
				},
				{
					bmin:    octree.Vector3D{0, 0, 0},
					bmax:    octree.Vector3D{3, 2, 2},
					returns: []*TestPoint{frontrighttop, center, frontcentertop},
				},
			},
		},
		{
			tree: octree.NewIterative(octree.Vector3D{0, 0, 0}, octree.Vector3D{3, 2, 2}, 8),
			tps:  []*TestPoint{backleftbottom, frontrighttop, center, backcenterbottom, backrightbottom, backlefttop, backcentertop, backrighttop, frontcentertop, frontlefttop},
			size: 10,
			queries: []testQuery{
				{
					bmin:    octree.Vector3D{-3, -2, -2},
					bmax:    octree.Vector3D{0, 0, 0},
					returns: []*TestPoint{center, backcenterbottom, backleftbottom},
				},
				{
					bmin:    octree.Vector3D{0, 0, 0},
					bmax:    octree.Vector3D{3, 2, 2},
					returns: []*TestPoint{frontcentertop, center, frontrighttop},
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
