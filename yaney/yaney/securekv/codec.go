package securekv

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
)

var defaultKeySecret = []byte("abcdabcdabcdabcd")
var defaultValSecret = []byte("abcdabcdabcdabcd")

const cipherLenSize = 4

type content struct {
	bytes     []byte
	keySecret []byte
	valSecret []byte
	key       []byte
	val       []byte
	offset    int
	decoded   bool
	encoded   bool
}

func encryptWithCTR(plaintext, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// initial vendor iv randomly for security
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, err
}

func decryptWithCTR(ciphertext, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext[aes.BlockSize:]))
	iv := ciphertext[:aes.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext, nil
}

func newContentFromKV(key, val, keySecret, valSecret []byte) *content {
	var c content

	c.key = key
	c.val = val
	c.keySecret = keySecret
	c.valSecret = valSecret
	c.decoded = true

	return &c
}

func newContentFromBytes(bytes, keySecret, valSecret []byte) *content {
	var c content

	c.bytes = bytes
	c.keySecret = keySecret
	c.valSecret = valSecret
	c.encoded = true

	return &c
}

func (c *content) read(p []byte) (n int, err error) {
	readLen := len(p)

	if readLen == 0 {
		return 0, nil
	}

	if readLen > len(c.bytes)-c.offset {
		return 0, errors.New("no enough bytes to read")
	}

	copy(p, c.bytes[c.offset:c.offset+readLen])
	c.offset += readLen

	return n, nil
}

func (c *content) toBytes() ([]byte, error) {
	if c.encoded {
		return c.bytes, nil
	}

	c.bytes = make([]byte, 0)
	// encrypt key
	keyCipher, err := encryptWithCTR(c.key, c.keySecret)
	if err != nil {
		panic(err)
	}
	keyCipherLen := make([]byte, cipherLenSize)
	binary.LittleEndian.PutUint32(keyCipherLen, uint32(len(keyCipher)))

	// encrypt val
	valCipher, err := encryptWithCTR(c.val, c.valSecret)
	if err != nil {
		panic(err)
	}
	valCipherLen := make([]byte, cipherLenSize)
	binary.LittleEndian.PutUint32(valCipherLen, uint32(len(valCipher)))

	// append cipher len and ciphertext to bytes
	c.bytes = append(c.bytes, keyCipherLen...)
	c.bytes = append(c.bytes, keyCipher...)
	c.bytes = append(c.bytes, valCipherLen...)
	c.bytes = append(c.bytes, valCipher...)

	return c.bytes, nil
}

func (c *content) toKV() ([]byte, []byte, error) {
	if c.decoded {
		return c.key, c.val, nil
	}

	// get cipher len
	keyCipherLen := make([]byte, cipherLenSize)
	if _, err := c.read(keyCipherLen); err != nil {
		return nil, nil, err
	}

	// get cipher
	keyCipher := make([]byte, binary.LittleEndian.Uint32(keyCipherLen))
	if _, err := c.read(keyCipher); err != nil {
		return nil, nil, err
	}

	var err error

	// decrypt cipher
	c.key, err = decryptWithCTR(keyCipher, c.keySecret)
	if err != nil {
		return nil, nil, err
	}

	// get cipher len
	valCipherLen := make([]byte, cipherLenSize)
	if _, err := c.read(valCipherLen); err != nil {
		return nil, nil, err
	}

	// get cipher
	valCipher := make([]byte, binary.LittleEndian.Uint32(valCipherLen))
	if _, err := c.read(valCipher); err != nil {
		return nil, nil, err
	}

	// decrypt cipher
	c.val, err = decryptWithCTR(valCipher, c.valSecret)
	if err != nil {
		return nil, nil, err
	}

	c.decoded = true

	return c.key, c.val, nil
}
