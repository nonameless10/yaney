package securekv

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
)

// other key-value storage such as rocksdb is also supported only if function below can be implemented

type badgerKV struct {
	db *badger.DB
}

func newBadgerKV() *badgerKV {
	var b badgerKV
	return &b
}

func (b *badgerKV) init(path string) error {
	var err error
	if b.db, err = badger.Open(badger.DefaultOptions(path)); err != nil {
		return err
	}
	return nil
}

func (b *badgerKV) put(bucket, key, val []byte) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(append(bucket, key...), val)
		return err
	})

	return err
}

func (b *badgerKV) exist(bucket, key []byte) (bool, error) {
	exist := true
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(append(bucket, key...))
		if err == nil {
			return nil
		}

		if err == badger.ErrKeyNotFound {
			exist = false
			return nil
		} else {
			return err
		}
	})

	return exist, err
}

func (b *badgerKV) get(bucket, key []byte) ([]byte, error) {
	var val []byte

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(bucket, key...))
		if err != nil {
			return err
		}

		err = item.Value(func(v []byte) error {
			val = v
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	return val, err
}

func (b *badgerKV) hasNext(bucket, key []byte) (bool, error) {
	hasNext := false
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iter := txn.NewIterator(opts)
		defer iter.Close()

		innerPrefix := make([]byte, 0, len(bucket)+len(key))
		innerPrefix = append(innerPrefix, bucket...)
		innerPrefix = append(innerPrefix, key...)

		iter.Seek(innerPrefix)
		if !iter.Valid() {
			return nil
		}

		iter.Next()
		if !iter.Valid() {
			return nil
		}

		hasNext = true
		return nil
	})

	return hasNext, err
}

func (b *badgerKV) next(bucket, key []byte) ([]byte, error) {
	nextKey := make([]byte, 0)
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iter := txn.NewIterator(opts)
		defer iter.Close()

		innerPrefix := make([]byte, 0, len(bucket)+len(key))
		innerPrefix = append(innerPrefix, bucket...)
		innerPrefix = append(innerPrefix, key...)

		iter.Seek(innerPrefix)
		if !iter.Valid() {
			return errors.New("not a valid key")
		}

		iter.Next()
		if !iter.Valid() {
			return errors.New("already at the end")
		}

		nextKey = iter.Item().Key()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nextKey[len(bucket):], nil
}

func (b *badgerKV) scan(bucket, prefix []byte) ([][]byte, error) {
	keySet := make([][]byte, 0)
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		iter := txn.NewIterator(opts)
		defer iter.Close()

		innerPrefix := make([]byte, 0, len(bucket)+len(prefix))
		innerPrefix = append(innerPrefix, bucket...)
		innerPrefix = append(innerPrefix, prefix...)

		for iter.Seek(innerPrefix); iter.ValidForPrefix(innerPrefix); iter.Next() {
			item := iter.Item()
			key := item.Key()
			keySet = append(keySet, key[len(bucket):])
		}

		return nil
	})

	return keySet, err
}

func (b *badgerKV) Close() {
	b.Close()
}
