package merkleTree

import (
	"sync"
)

type Tree struct {
	sync.RWMutex
	eng Engine
	cfg Config
}
