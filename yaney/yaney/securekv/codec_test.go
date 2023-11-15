package securekv

import (
	"reflect"
	"testing"
)

func TestKVCodec(t *testing.T) {
	var key = []byte("test key")
	var val = []byte("test val")

	if len(defaultValSecret) != 16 || len(defaultKeySecret) != 16 {
		t.Fail()
	}

	content1 := newContentFromKV(key, val, defaultKeySecret, defaultValSecret)
	bytes, _ := content1.toBytes()
	content2 := newContentFromBytes(bytes, defaultKeySecret, defaultValSecret)
	k, v, err := content2.toKV()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(k, key) || !reflect.DeepEqual(v, val) {
		t.Fatal()
	}
}
