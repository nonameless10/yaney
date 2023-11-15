package securegraph

type KVStorage interface {
	Build(keys, vals [][]byte)
	Exist(key []byte) bool
	Get(key []byte) []byte
	HasNext(key []byte) bool
	Next(key []byte) []byte
	Close()
}
