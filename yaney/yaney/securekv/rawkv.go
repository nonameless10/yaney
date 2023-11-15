package securekv

import "sort"

type RawKV struct {
	kv *badgerKV
	defaultBucket []byte
}

func NewRawKV(dbPath string) *RawKV{
	var s RawKV
	s.kv = newBadgerKV()
	s.defaultBucket = []byte("0")

	err := s.kv.init(dbPath)
	if err != nil {
		panic(err)
	}
	return &s
}

func (r *RawKV)HasNext(key []byte) bool{
	rtn,_ := r.kv.hasNext(r.defaultBucket, key)
	return rtn
}

func (r *RawKV)Next(key []byte) []byte{
	rtn,_ := r.kv.next(r.defaultBucket,key)
	return rtn
}

func (r *RawKV)Exist(key []byte) bool{
	rtn, _ := r.kv.exist(r.defaultBucket, key)
	return rtn
}

func (r* RawKV)Get(key []byte) []byte{
	rtn, _ := r.kv.get(r.defaultBucket, key);
	return rtn
}

func (r *RawKV)Build(keys, vals [][]byte){
	if len(keys) != len(vals) {
		panic("length of keys and vals are not matched")
	}
	kvs := make([]KVPair, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		kv := KVPair{
			Key: keys[i],
			Val: vals[i],
		}
		kvs = append(kvs, kv)
	}
	sort.Sort(ByKey(kvs))

	for i := 0; i < len(kvs); i++{
		r.kv.put(r.defaultBucket, kvs[i].Key, kvs[i].Val)
	}
}

func (s *RawKV) Close(){
	s.kv.db.Close()
}