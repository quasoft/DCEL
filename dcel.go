// Package dcel implements a doubly connected edge list data structure, and is used to represent
// planar graphs in a plane. This implementation is intended to be used in 2D space only.
package dcel

import (
	"fmt"
	"strconv"
)

// DCEL stores the state of the data structure and provides methods for linking of three sets of
// objects: vertecies, edges and faces.
type DCEL struct {
	Vertices  []*Vertex
	Faces     []*Face
	HalfEdges []*HalfEdge
}

// Vertex represents a node in the DCEL structure. Each vertex has 2D coordinates and a pointer
// to an arbitrary half edge that has this vertex as its target (origin). Annotations (user data)
// can be stored in the Data field.
type Vertex struct {
	X, Y     int
	HalfEdge *HalfEdge
	Data     interface{}
}

// Face represents a subdivision of the plane. Each face has a pointer to one of the half edges
// at its boundary. Faces can have user specified IDs and annotations.
type Face struct {
	HalfEdge *HalfEdge
	ID       int64
	Data     interface{}
}

// HalfEdge represents one of the half-edges in an edge pair. Each half-edge has a pointer to its
// target vertex (origin), the face to which it belongs, its twin edge (a reversed half-edge, pointing
// to a neighbour face) and pointers to the next and previous half-edges at the boundary of its face.
// Half-edges can also store user data.
type HalfEdge struct {
	Target *Vertex
	Face   *Face
	Twin   *HalfEdge
	Next   *HalfEdge
	Prev   *HalfEdge
	Data   interface{}
}

func (v *Vertex) String() string {
	return fmt.Sprintf("{Vertex %p; X,Y: %d,%d; Edge: %p}", v, v.X, v.Y, v.HalfEdge)
}

func (f *Face) String() string {
	return fmt.Sprintf("{Face #%d %p}", f.ID, f)
}

func (he *HalfEdge) String() string {
	faceID := "nil"
	if he.Face != nil {
		faceID = "#" + strconv.FormatInt(he.Face.ID, 10)
	}

	return fmt.Sprintf("{Edge %p; Target: %d,%d; Twin: %p; Face: %s}", he, he.Target.X, he.Target.Y, he.Twin, faceID)
}

// IsClosed returns true if both half-edges in the pair have a target vertex.
func (he *HalfEdge) IsClosed() bool {
	return he.Target != nil && he.Twin != nil && he.Twin.Target != nil
}

// NewDCEL creates a new DCEL data structure.
func NewDCEL() *DCEL {
	return &DCEL{}
}

// NewFace creates a new face and stores it in the DCEL structure.
func (d *DCEL) NewFace() *Face {
	face := &Face{}
	d.Faces = append(d.Faces, face)
	return face
}

// NewVertex creates a new vertex with the given coordinates and stores it in the structure.
func (d *DCEL) NewVertex(x, y int) *Vertex {
	vertex := &Vertex{
		X: x,
		Y: y,
	}
	d.Vertices = append(d.Vertices, vertex)
	return vertex
}

// NewHalfEdge creates a new half-edge starting at the given vertex and stores it in the structure.
func (d *DCEL) NewHalfEdge(face *Face, vertex *Vertex) *HalfEdge {
	halfEdge := &HalfEdge{
		Face:   face,
		Target: vertex,
	}
	if face.HalfEdge == nil {
		face.HalfEdge = halfEdge
	}
	if vertex.HalfEdge == nil {
		vertex.HalfEdge = halfEdge
	}
	d.HalfEdges = append(d.HalfEdges, halfEdge)
	return halfEdge
}

// NewEdge creates a pair of half-edges, one of them starting at the given vertex.
func (d *DCEL) NewEdge(face1, face2 *Face, vertex *Vertex) (*HalfEdge, *HalfEdge) {
	halfEdge := d.NewHalfEdge(face1, vertex)
	twin := &HalfEdge{
		Face: face2,
		Twin: halfEdge,
	}
	halfEdge.Twin = twin
	if face2.HalfEdge == nil {
		face2.HalfEdge = twin
	}
	d.HalfEdges = append(d.HalfEdges, twin)
	return halfEdge, twin
}
