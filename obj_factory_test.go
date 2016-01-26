package merkleTree

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func genString(l int) string {
	return hex.EncodeToString(genBinary((l + 2) / 1))[0:l]
}

func genBinary(l int) []byte {
	b := make([]byte, l)
	rand.Read(b)
	return b
}

type testValue struct {
	seqno int
	key   string
}

type objFactory struct {
	objs  map[string]KeyValuePair
	seqno int
}

func newObjFactory() *objFactory {
	return &objFactory{
		objs: make(map[string]KeyValuePair),
	}
}

func (of *objFactory) dumpAll() []KeyValuePair {
	var ret []KeyValuePair
	for _, v := range of.objs {
		ret = append(ret, v)
	}
	return ret
}

func (of *objFactory) produce() KeyValuePair {
	key := genBinary(8)
	keyString := hex.EncodeToString(key)
	val := testValue{seqno: of.seqno, key: keyString}
	of.seqno++
	kvp := KeyValuePair{Key: key, Value: val}
	of.objs[keyString] = kvp
	return kvp
}

func (of *objFactory) modifySome(mod int) {
	i := 0
	for _, v := range of.objs {
		tv, ok := v.Value.(testValue)
		if !ok {
			panic(fmt.Sprintf("Got value of wrong type: %T", v))
		}
		if (i % mod) == 0 {
			tv.seqno *= 2
			tv.key += "-yo-dawg"
		}
		i++
	}
}
