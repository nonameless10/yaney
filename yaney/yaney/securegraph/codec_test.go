package securegraph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

func TestGobMap(t *testing.T) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	props := make(map[string]string)
	props["1"] = "1"
	props["2"] = "2"
	props["3"] = "3"
	props["4"] = "4"
	err := encoder.Encode(props)
	if err != nil {
		panic(err)
	}

	var newProps map[string]string
	var newBuffer bytes.Buffer
	newBuffer.Write(buffer.Bytes())
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&newProps)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", newProps)
}

func TestGobNil(t *testing.T) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)

	var props map[string]string
	err := encoder.Encode(props)
	if err != nil {
		panic(err)
	}

	var newProps map[string]string = nil

	var newBuffer bytes.Buffer
	newBuffer.Write(buffer.Bytes())
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&newProps)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", newProps)
}

func TestVertexCodec(t *testing.T) {
	var vertex Vertex
	vertex.name = "jo3yzhu"

	props := make(map[string]string)
	props["1"] = "1"
	props["2"] = "2"
	props["3"] = "3"
	props["4"] = "4"
	vertex.props = props

	vc := NewVertexCodecFromVertex(vertex)
	keyBytes, valBytes := vc.ToBytes()

	newVc := NewVertexCodecFromBytes(keyBytes, valBytes)
	newVertex := newVc.ToVertex()

	t.Log(newVertex)
}

func TestEdgeCode(t *testing.T) {
	var edge Edge
	edge.selfVertexName = "self"
	edge.inout = 3
	edge.relationName = "goto"
	edge.rank = 5
	edge.peerVertexName = "peer"

	props := make(map[string]string)
	props["1"] = "1"
	props["2"] = "2"
	props["3"] = "3"
	props["4"] = "4"
	edge.props = props

	ec := NewEdgeCodecFromEdge(edge)
	keyBytes, valBytes := ec.ToBytes()

	newEc := NewEdgeCodecFromBytes(keyBytes, valBytes)
	newEdge := newEc.ToEdge()

	t.Log(newEdge)
}
