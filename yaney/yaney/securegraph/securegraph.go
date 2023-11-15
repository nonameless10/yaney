package securegraph

import (
	"bytes"
	"errors"
	"github.com/nonameless10/yaney/securekv"
)

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

type SecureGraph struct {
	kv KVStorage
}

func NewSecureGraph(dbPath string, masterSecret, keySecret, valSecret []byte) *SecureGraph {
	var sg SecureGraph

	sg.kv = securekv.NewSecureKV(dbPath, masterSecret, keySecret, valSecret)

	return &sg
}



func (sg *SecureGraph) BuildGraph(vgs []VertexGraph) {
	keys := make([][]byte, 0)
	vals := make([][]byte, 0)

	for _, vg := range vgs {
		vertexCodec := NewVertexCodecFromVertex(vg.vertex)
		vk, vv := vertexCodec.ToBytes()
		keys = append(keys, vk)
		vals = append(vals, vv)

		for _, ie := range vg.inEdges {
			inEdgeCodec := NewEdgeCodecFromEdge(ie)
			iek, iev := inEdgeCodec.ToBytes()
			keys = append(keys, iek)
			vals = append(vals, iev)
		}

		for _, oe := range vg.outEdges {
			outEdgeCodec := NewEdgeCodecFromEdge(oe)
			oek, oev := outEdgeCodec.ToBytes()
			keys = append(keys, oek)
			vals = append(vals, oev)
		}
	}

	sg.kv.Build(keys, vals)
}

func (sg *SecureGraph) QueryGraph(vertexName string) (*VertexGraph, error) {
	var vg VertexGraph

	vertexCodec := NewVertexCodecFromVertex(Vertex{
		name:  vertexName,
		props: nil,
	})

	vertexKey, _ := vertexCodec.ToBytes()
	if !sg.kv.Exist(vertexKey) {
		return nil, errors.New("vertex doesn't exist")
	}

	// get vertex
	vertexVal := sg.kv.Get(vertexKey)
	vg.vertex = NewVertexCodecFromBytes(vertexKey, vertexVal).ToVertex()

	edgeKey := vertexKey
	for {
		// get edges
		if !sg.kv.HasNext(edgeKey) {
			break
		} else {
			edgeKey = sg.kv.Next(edgeKey)
			if !bytes.HasPrefix(edgeKey, vertexKey) {
				break
			}

			edgeVal := sg.kv.Get(edgeKey)
			edge := NewEdgeCodecFromBytes(edgeKey, edgeVal).ToEdge()
			if edge.inout > 0 {
				vg.outEdges = append(vg.outEdges, edge)
			} else if edge.inout < 0 {
				vg.inEdges = append(vg.inEdges, edge)
			} else {
				panic("never have a edge with 0 inout")
			}
		}
	}

	return &vg, nil
}
