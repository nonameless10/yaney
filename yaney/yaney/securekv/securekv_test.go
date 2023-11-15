package securekv

import (
	"encoding/binary"
	"github.com/jo3yzhu/yaney/securekv/sse"
	"strconv"
	"testing"
)

func TestBuildOrder(t *testing.T) {
	b := newBadgerKV()
	err := b.init("builder_test")
	if err != nil {
		t.Fail()
	}

	bytes := make([]byte, 8)

	for i := 0; i < 1024; i += 16 {
		binary.BigEndian.PutUint64(bytes, uint64(i))
		err = b.put([]byte("bucket"), bytes, []byte(strconv.Itoa(i)))
	}

	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, uint64(0))

	for {
		if key == nil || err != nil {
			break
		}
		val, err := b.get([]byte("bucket"), key)
		if err != nil {
			t.Fail()
		}

		t.Log(string(val))

		key, err = b.next([]byte("bucket"), key)
	}
}

func getSecret(passphrase, salt string, iter int) []byte {
	return sse.Key([]byte(passphrase), []byte(salt), iter)
}

func TestSecureKV(t *testing.T) {
	masterSecret := getSecret("master", "secret", 4096)
	keySecret := getSecret("key", "secret", 4096)
	valSecret := getSecret("val", "secret", 4096)

	keys := make([][]byte, 0)
	vals := make([][]byte, 0)

	for i := 0; i < 100; i++ {
		keys = append(keys, []byte(strconv.Itoa(i)))
		vals = append(vals, []byte(strconv.Itoa(i)))
	}

	kv := NewSecureKV("secure_kv_test", masterSecret, keySecret, valSecret)
	kv.Build(keys, vals)

	key := []byte(strconv.Itoa(0))
	hasNext := true

	for {
		val := kv.Get(key)
		t.Log(string(key), string(val))

		hasNext = kv.HasNext(key)
		if !hasNext {
			break
		}

		key = kv.Next(key)
	}

}
