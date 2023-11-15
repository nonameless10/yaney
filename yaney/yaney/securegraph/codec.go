package securegraph

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

const (
	vertexNameLen = 4
)

const (
	edgeRelationTypeLen = 4
	rankLen             = 4
)

type VertexCodec struct {
	vertex    Vertex
	keyBuffer bytes.Buffer
	valBuffer bytes.Buffer
	encoded   bool
	decoded   bool
}

func NewVertexCodecFromBytes(keyBytes, valBytes []byte) *VertexCodec {
	var vc VertexCodec
	vc.keyBuffer.Write(keyBytes)
	vc.valBuffer.Write(valBytes)
	vc.encoded = true

	return &vc
}

func NewVertexCodecFromVertex(vertex Vertex) *VertexCodec {
	var vc VertexCodec
	vc.vertex = vertex
	vc.decoded = true

	return &vc
}

func (vc *VertexCodec) ToVertex() Vertex {
	if vc.decoded {
		return vc.vertex
	}

	// get vertex name len
	nameLenBuf := make([]byte, vertexNameLen)
	if n, err := vc.keyBuffer.Read(nameLenBuf); err != nil || n != vertexNameLen {
		panic(err)
	}
	nameLen := binary.LittleEndian.Uint32(nameLenBuf)

	// get vertex name
	nameBuf := make([]byte, nameLen)
	if n, err := vc.keyBuffer.Read(nameBuf); err != nil || uint32(n) != nameLen {
		panic(err)
	}
	vc.vertex.name = string(nameBuf)

	// deserialize props
	d := gob.NewDecoder(&vc.valBuffer)
	if err := d.Decode(&vc.vertex.props); err != nil {
		panic(err)
	}

	vc.decoded = true

	return vc.vertex
}

func (vc *VertexCodec) ToBytes() ([]byte, []byte) {
	if vc.encoded {
		return vc.keyBuffer.Bytes(), vc.valBuffer.Bytes()
	}

	// encode vertex name
	nameLen := len(vc.vertex.name)
	nameLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(nameLenBytes, uint32(nameLen))
	vc.keyBuffer.Write(nameLenBytes)
	vc.keyBuffer.Write([]byte(vc.vertex.name))

	// encode vertex props
	e := gob.NewEncoder(&vc.valBuffer)
	err := e.Encode(vc.vertex.props)
	if err != nil {
		panic(err)
	}

	vc.encoded = true

	return vc.keyBuffer.Bytes(), vc.valBuffer.Bytes()
}

type EdgeCodec struct {
	edge      Edge
	keyBuffer bytes.Buffer
	valBuffer bytes.Buffer
	encoded   bool
	decoded   bool
}

func NewEdgeCodecFromBytes(keyBytes, valBytes []byte) *EdgeCodec {
	var ec EdgeCodec
	ec.keyBuffer.Write(keyBytes)
	ec.valBuffer.Write(valBytes)
	ec.encoded = true

	return &ec
}

func NewEdgeCodecFromEdge(edge Edge) *EdgeCodec {
	var ec EdgeCodec
	ec.edge = edge
	ec.decoded = true

	return &ec
}

func (ec *EdgeCodec) ToEdge() Edge {
	if ec.decoded {
		return ec.edge
	}

	// get self vertex name len
	selfVertexNameLenBuf := make([]byte, vertexNameLen)
	if n, err := ec.keyBuffer.Read(selfVertexNameLenBuf); err != nil || n != vertexNameLen {
		panic(err)
	}
	selfVertexNameLen := binary.LittleEndian.Uint32(selfVertexNameLenBuf)

	// get self vertex name
	selfVertexNameBuf := make([]byte, selfVertexNameLen)
	if n, err := ec.keyBuffer.Read(selfVertexNameBuf); err != nil || uint32(n) != selfVertexNameLen {
		panic(err)
	}
	ec.edge.selfVertexName = string(selfVertexNameBuf)

	// get edge type
	edgeType, err := ec.keyBuffer.ReadByte()
	if err != nil || edgeType == 0 {
		panic(err)
	}
	ec.edge.inout = int8(edgeType)

	// get relation type len
	relationNameLenBuf := make([]byte, edgeRelationTypeLen)
	if n, err := ec.keyBuffer.Read(relationNameLenBuf); err != nil || n != edgeRelationTypeLen {
		panic(err)
	}
	relationNameLen := binary.LittleEndian.Uint32(relationNameLenBuf)

	// get relation type
	relationNameBuf := make([]byte, relationNameLen)
	if n, err := ec.keyBuffer.Read(relationNameBuf); err != nil || uint32(n) != relationNameLen {
		panic(err)
	}
	ec.edge.relationName = string(relationNameBuf)

	// get rank
	rankBuf := make([]byte, rankLen)
	if n, err := ec.keyBuffer.Read(rankBuf); err != nil || n != rankLen {
		panic(err)
	}
	ec.edge.rank = binary.LittleEndian.Uint32(rankBuf)

	// get peer vertex name len
	peerVertexNameLenBuf := make([]byte, vertexNameLen)
	if n, err := ec.keyBuffer.Read(peerVertexNameLenBuf); err != nil || n != vertexNameLen {
		panic(err)
	}
	peerVertexNameLen := binary.LittleEndian.Uint32(peerVertexNameLenBuf)

	// get self vertex name
	peerVertexNameBuf := make([]byte, peerVertexNameLen)
	if n, err := ec.keyBuffer.Read(peerVertexNameBuf); err != nil || uint32(n) != peerVertexNameLen {
		panic(err)
	}
	ec.edge.peerVertexName = string(peerVertexNameBuf)

	// decode properties
	d := gob.NewDecoder(&ec.valBuffer)
	if err := d.Decode(&ec.edge.props); err != nil {
		panic(err)
	}

	ec.decoded = true

	return ec.edge
}

func (ec *EdgeCodec) ToBytes() ([]byte, []byte) {
	if ec.encoded {
		return ec.keyBuffer.Bytes(), ec.valBuffer.Bytes()
	}

	// encode self vertex name
	selfVertexNameLen := len(ec.edge.selfVertexName)
	selfVertexNameBytes := make([]byte, vertexNameLen)
	binary.LittleEndian.PutUint32(selfVertexNameBytes, uint32(selfVertexNameLen))
	ec.keyBuffer.Write(selfVertexNameBytes)
	ec.keyBuffer.Write([]byte(ec.edge.selfVertexName))

	// encode edge type
	ec.keyBuffer.WriteByte(byte(ec.edge.inout))

	// encode relation name
	relationTypeLen := len(ec.edge.relationName)
	relationTypeLenBytes := make([]byte, vertexNameLen)
	binary.LittleEndian.PutUint32(relationTypeLenBytes, uint32(relationTypeLen))
	ec.keyBuffer.Write(relationTypeLenBytes)
	ec.keyBuffer.Write([]byte(ec.edge.relationName))

	// encode rank
	rankBytes := make([]byte, rankLen)
	binary.LittleEndian.PutUint32(rankBytes, ec.edge.rank)
	ec.keyBuffer.Write(rankBytes)

	// encode peer vertex name
	peerVertexNameLen := len(ec.edge.peerVertexName)
	peerVertexNameBytes := make([]byte, vertexNameLen)
	binary.LittleEndian.PutUint32(peerVertexNameBytes, uint32(peerVertexNameLen))
	ec.keyBuffer.Write(peerVertexNameBytes)
	ec.keyBuffer.Write([]byte(ec.edge.peerVertexName))

	// encode properties
	e := gob.NewEncoder(&ec.valBuffer)
	err := e.Encode(ec.edge.props)
	if err != nil {
		panic(err)
	}

	ec.encoded = true

	return ec.keyBuffer.Bytes(), ec.valBuffer.Bytes()
}
