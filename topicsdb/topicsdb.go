package topicsdb

import (
	"context"
	"fmt"

	"github.com/unicornultrafoundation/go-helios/native/idx"
	"github.com/unicornultrafoundation/go-helios/u2udb"
	"github.com/unicornultrafoundation/go-u2u/common"
	"github.com/unicornultrafoundation/go-u2u/core/types"
)

const maxTopicsCount = 5 // count is limited hard to 5 by EVM (see LOG0...LOG4 ops)

var (
	ErrEmptyTopics  = fmt.Errorf("empty topics")
	ErrTooBigTopics = fmt.Errorf("too many topics")
)

// Index is a specialized indexes for log records storing and fetching.
type Index interface {
	FindInBlocks(ctx context.Context, from, to idx.Block, pattern [][]common.Hash) (logs []*types.Log, err error)
	ForEachInBlocks(ctx context.Context, from, to idx.Block, pattern [][]common.Hash, onLog func(*types.Log) (gonext bool)) error
	Push(recs ...*types.Log) error
	Close()

	WrapTablesAsBatched() (unwrap func())
}

// New Index instance.
func New(dbs u2udb.DBProducer) Index {
	tt := newIndex(dbs)

	return tt
}

// New Index instance consumes limited threads count.
func NewWithThreadPool(dbs u2udb.DBProducer) Index {
	tt := newIndex(dbs)
	return &withThreadPool{tt}
}

func limitPattern(pattern [][]common.Hash) (limited [][]common.Hash, err error) {
	if len(pattern) > (maxTopicsCount + 1) {
		limited = make([][]common.Hash, (maxTopicsCount + 1))
	} else {
		limited = make([][]common.Hash, len(pattern))
	}
	copy(limited, pattern)

	ok := false
	for i, variants := range limited {
		ok = ok || len(variants) > 0
		if len(variants) > 1 {
			limited[i] = uniqOnly(variants)
		}
	}
	if !ok {
		err = ErrEmptyTopics
		return
	}

	return
}

func uniqOnly(hh []common.Hash) []common.Hash {
	index := make(map[common.Hash]struct{}, len(hh))
	for _, h := range hh {
		index[h] = struct{}{}
	}

	uniq := make([]common.Hash, 0, len(index))
	for h := range index {
		uniq = append(uniq, h)
	}
	return uniq
}
