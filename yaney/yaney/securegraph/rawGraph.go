package securegraph

import (
	"bytes"
	"errors"
	"github.com/jo3yzhu/yaney/securekv"
)

type RawGraph struct {
	kv KVStorage
}

func NewRawGraph(dbPath string) *RawGraph{
	var rg RawGraph
	rg.kv = securekv.NewRawKV(dbPath)
	return &rg
}


func (rg *RawGraph) BuildGraph(vgs []VertexGraph) {
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
	rg.kv.Build(keys, vals)
}

func (rg *RawGraph) QueryGraph(vertexName string) (*VertexGraph, error) {
	var vg VertexGraph

	vertexCodec := NewVertexCodecFromVertex(Vertex{
		name:  vertexName,
		props: nil,
	})
	vertexKey, _ := vertexCodec.ToBytes()
	if !rg.kv.Exist(vertexKey) {
		return nil, errors.New("vertex doesn't exist")
	}
	// get vertex
	vertexVal := rg.kv.Get(vertexKey)
	vg.vertex = NewVertexCodecFromBytes(vertexKey, vertexVal).ToVertex()

	edgeKey := vertexKey
	for {
		// get edges
		if !rg.kv.HasNext(edgeKey) {
			break
		} else {
			edgeKey = rg.kv.Next(edgeKey)
			if !bytes.HasPrefix(edgeKey, vertexKey) {
				break
			}

			edgeVal := rg.kv.Get(edgeKey)
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