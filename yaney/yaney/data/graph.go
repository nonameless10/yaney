package data
type Vertex struct {
	name  string
	props map[string]string // property of vertex
}

type Edge struct {
	selfVertexName string // vertex name
	inout          int8   // edge type, out edge when inout > 0, in edge when inout < 0
	relationName   string // relation name of a certain edge
	rank           uint32
	peerVertexName string            // vertex name
	props          map[string]string // property of edge
}

type VertexGraph struct {
	vertex   Vertex
	inEdges  []Edge
	outEdges []Edge
}


