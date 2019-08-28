package merkleTree

import (
	"crypto/rand"
	"encoding/hex"
)

func genBinary(l int) []byte {
	b := make([]byte, l)
	_, _ = rand.Read(b)
	return b
}

type testValue struct {
	_struct bool `codec:",toarray"` //nolint
	Seqno   int
	Key     string
	KeyRaw  []byte
}

// TestObjFactory generates a bunch of test objects for debugging
type TestObjFactory struct {
	objs  map[string]KeyValuePair
	seqno int
}

// NewTestObjFactor makes a new object factory for testing
func NewTestObjFactory() *TestObjFactory {
	return &TestObjFactory{
		objs: make(map[string]KeyValuePair),
	}
}

// Produce one test object
func (of *TestObjFactory) Produce() KeyValuePair {
	key := genBinary(8)
	keyString := hex.EncodeToString(key)
	val := testValue{Seqno: of.seqno, Key: keyString, KeyRaw: key}
	of.seqno++
	kvp := KeyValuePair{Key: key, Value: val}
	of.objs[keyString] = kvp
	return kvp
}

// Mproduce makes many test objects.
func (of *TestObjFactory) Mproduce(n int) []KeyValuePair {
	var ret []KeyValuePair
	for i := 0; i < n; i++ {
		ret = append(ret, of.Produce())
	}
	return ret
}

func (of *TestObjFactory) Construct() interface{} {
	return testValue{}
}
