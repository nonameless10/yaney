package securekv

import (
	"reflect"
	"testing"
)

func TestDbPutGet(t *testing.T) {
	bucket := []byte("1")
	key := []byte("key")
	val := []byte("val")

	b := newBadgerKV()
	err := b.init("badger_test")
	if err != nil {
		t.Fail()
	}

	err = b.put(bucket, key, val)
	if err != nil {
		t.Fatal(err)
	}

	v, err := b.get(bucket, key)
	if err != nil || !reflect.DeepEqual(val, v) {
		t.Fail()
	}

	// if put on the same key, just overwrite it
	err = b.put(bucket, key, append(val, []byte(val)...))
	if err != nil {
		t.Fatal(err)
	}

	v, err = b.get(bucket, key)
	if err != nil {
		t.Fail()
	}
	t.Log(v)
}

func TestDbNext(t *testing.T) {
	b := newBadgerKV()
	err := b.init("badger_test")
	if err != nil {
		t.Fail()
	}

	bucket := []byte("1")
	err = b.put(bucket, []byte("11"), []byte("111"))
	err = b.put(bucket, []byte("12"), []byte("122"))
	err = b.put(bucket, []byte("13"), []byte("133"))
	err = b.put(bucket, []byte("14"), []byte("144"))
	err = b.put(bucket, []byte("15"), []byte("155"))

	key := []byte("11")

	for {
		if key == nil || err != nil {
			break
		}
		t.Log(key)
		key, err = b.next(bucket, key)
	}
}

func TestDbScan(t *testing.T) {
	b := newBadgerKV()
	err := b.init("badger_test")
	if err != nil {
		t.Fail()
	}

	bucket := []byte("1")
	err = b.put(bucket, []byte("11"), []byte("111"))
	err = b.put(bucket, []byte("12"), []byte("122"))
	err = b.put(bucket, []byte("13"), []byte("133"))
	err = b.put(bucket, []byte("14"), []byte("144"))
	err = b.put(bucket, []byte("15"), []byte("155"))

	keys, err := b.scan(bucket, []byte("1"))
	t.Log(keys)

	vals := make([][]byte, 0)
	for _, key := range keys {
		val, err := b.get(bucket, key)
		if err != nil {
			t.Fail()
			return
		}
		vals = append(vals, val)
	}
	t.Log(vals)

}
