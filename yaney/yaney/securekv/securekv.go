package securekv

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/jo3yzhu/yaney/securekv/sse"
	"math/big"
	"reflect"
	"sort"
)

type KVPair struct {
	Key []byte
	Val []byte
}

type ByKey []KVPair

func (kvs ByKey) Len() int {
	return len(kvs)
}

func (kvs ByKey) Swap(i, j int) {
	kvs[i], kvs[j] = kvs[j], kvs[i]
}

func (kvs ByKey) Less(i, j int) bool {
	return bytes.Compare(kvs[i].Key, kvs[j].Key) < 0
}

type SecureKV struct {
	kv           *badgerKV
	indexBucket  []byte
	dataBucket   []byte
	masterSecret []byte
	keySecret    []byte
	valSecret    []byte
}

func NewSecureKV(dbPath string, masterSecret, keySecret, valSecret []byte) *SecureKV {
	var s SecureKV
	s.kv = newBadgerKV()
	s.indexBucket = []byte("0")
	s.dataBucket = []byte("1")
	s.masterSecret = masterSecret
	s.keySecret = keySecret
	s.valSecret = valSecret

	err := s.kv.init(dbPath)
	if err != nil {
		panic(err)
	}

	return &s
}

func (s *SecureKV) HasNext(key []byte) bool {
	index := sse.HMAC(append(key, sse.One[:]...), s.masterSecret)
	secret := sse.HMAC(append(key, sse.Two[:]...), s.masterSecret)

	hasNext, err := s.innerServerHasNext(index, secret)
	if err != nil {
		panic(err)
	}

	return hasNext
}

func (s *SecureKV) Next(key []byte) []byte {
	index := sse.HMAC(append(key, sse.One[:]...), s.masterSecret)
	secret := sse.HMAC(append(key, sse.Two[:]...), s.masterSecret)

	nextContent, err := s.innerServerNext(index, secret)
	if err != nil {
		panic(err)
	}

	content := newContentFromBytes(nextContent, s.keySecret, s.valSecret)
	k, _, err := content.toKV()
	return k
}

func (s *SecureKV) Exist(key []byte) bool {
	index := sse.HMAC(append(key, sse.One[:]...), s.masterSecret)
	secret := sse.HMAC(append(key, sse.Two[:]...), s.masterSecret)

	if exist, err := s.innerServerExist(index, secret); !exist || err != nil {
		return false
	} else {
		return true
	}
}

func (s *SecureKV) Get(key []byte) []byte {
	index := sse.HMAC(append(key, sse.One[:]...), s.masterSecret)
	secret := sse.HMAC(append(key, sse.Two[:]...), s.masterSecret)

	b, _ := s.innerServerGet(index, secret)
	content := newContentFromBytes(b, s.keySecret, s.valSecret)
	k, v, err := content.toKV()
	if err != nil || !reflect.DeepEqual(k, key) {
		panic(err)
	}

	return v
}

// Build import key-value pairs to SecureKV sequentially
func (s *SecureKV) Build(keys, vals [][]byte) {
	if len(keys) != len(vals) {
		panic("length of keys and vals are not matched")
	}

	// sort kv pairs by bitwise order of key
	kvs := make([]KVPair, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		kv := KVPair{
			Key: keys[i],
			Val: vals[i],
		}
		kvs = append(kvs, kv)
	}

	sort.Sort(ByKey(kvs))

	var seq uint64 = 0
	pos := make([]byte, 8)

	for i := 0; i < len(kvs); i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(15))
		seq = seq + n.Uint64() + 1
		binary.BigEndian.PutUint64(pos, seq)

		if err := s.put(kvs[i].Key, kvs[i].Val, pos); err != nil {
			panic(err)
		}
	}
}

func (s *SecureKV) Close(){
	s.kv.db.Close()
}

func (s *SecureKV) put(key, val, pos []byte) error {
	// start 
	index := sse.HMAC(append(key, sse.One[:]...), s.masterSecret)
	secret := sse.HMAC(append(key, sse.Two[:]...), s.masterSecret)
	content := newContentFromKV(key, val, s.keySecret, s.valSecret)
	// end
	b, _ := content.toBytes()

	return s.innerServerPut(index, secret, b, pos)
}

func (s *SecureKV) innerServerPut(index, secret, content, pos []byte) error {
	// put content in data bucket
	if err := s.kv.put(s.dataBucket, pos, content); err != nil {
		return err
	}

	// build index
	h := sse.HMAC([]byte("COUNT"), index)
	if epos, err := sse.Encrypt(pos, secret); err != nil {
		return err
	} else {
		if err := s.kv.put(s.indexBucket, h, epos); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (s *SecureKV) innerServerExist(index, secret []byte) (bool, error) {
	// get index
	h := sse.HMAC([]byte("COUNT"), index)
	if exist, err := s.kv.exist(s.indexBucket, h); err != nil || !exist {
		return false, err
	} else {
		epos, _ := s.kv.get(s.indexBucket, h)
		if pos, err := sse.Decrypt(epos, secret); err != nil {
			panic(err)
		} else {
			return s.kv.exist(s.dataBucket, pos)
		}
	}
}

func (s *SecureKV) innerServerGet(index, secret []byte) ([]byte, error) {
	// get index
	h := sse.HMAC([]byte("COUNT"), index)
	if epos, err := s.kv.get(s.indexBucket, h); err != nil {
		return nil, err
	} else {
		// get data by index
		if pos, err := sse.Decrypt(epos, secret); err != nil {
			panic(err)
		} else {
			if content, err := s.kv.get(s.dataBucket, pos); err != nil {
				return nil, err
			} else {
				return content, nil
			}
		}
	}
}

func (s *SecureKV) innerServerHasNext(index, secret []byte) (bool, error) {
	// get index
	h := sse.HMAC([]byte("COUNT"), index)
	if epos, err := s.kv.get(s.indexBucket, h); err != nil {
		return false, err
	} else {
		// get data by index
		if pos, err := sse.Decrypt(epos, secret); err != nil {
			return false, err
		} else {
			// get next content and return
			hasNext, err := s.kv.hasNext(s.dataBucket, pos)
			return hasNext, err
		}
	}
}

func (s *SecureKV) innerServerNext(index, secret []byte) ([]byte, error) {
	// get index
	h := sse.HMAC([]byte("COUNT"), index)
	if epos, err := s.kv.get(s.indexBucket, h); err != nil {
		return nil, err
	} else {
		// get data by index
		if pos, err := sse.Decrypt(epos, secret); err != nil {
			return nil, err
		} else {
			// get next content and return
			if nextPos, err := s.kv.next(s.dataBucket, pos); err != nil {
				return nil, err

			} else {
				nextContent, err := s.kv.get(s.dataBucket, nextPos)
				if err != nil {
					return nil, err
				}
				return nextContent, err
			}
		}
	}
}
