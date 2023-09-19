package topicsdb

import (
	"context"
	"sync"

	"github.com/unicornultrafoundation/go-u2u/libs/common"
)

type logHandler func(rec *logrec) (gonext bool, err error)

func (tt *index) searchParallel(ctx context.Context, pattern [][]common.Hash, blockStart, blockEnd uint64, onMatched logHandler, onDbIterator func()) error {
	if ctx == nil {
		ctx = context.Background()
	}

	var (
		syncing      = newSynchronizator()
		mu           sync.Mutex
		foundByBlock = make(map[uint64]map[ID]*logrec)
	)

	aggregator := func(pos, num int) logHandler {
		return func(rec *logrec) (gonext bool, err error) {
			if rec == nil {
				syncing.FinishThread(pos, num)
				return
			}

			err = ctx.Err()
			if err != nil {
				return
			}

			block := rec.ID.BlockNumber()
			if blockEnd > 0 && block > blockEnd {
				return
			}
			if rec.topicsCount < uint8(len(pattern)-1) {
				gonext = true
				return
			}

			var prevBlock uint64
			prevBlock, gonext = syncing.GoNext(block)
			if !gonext {
				return
			}

			mu.Lock()
			defer mu.Unlock()

			if prevBlock > 0 {
				delete(foundByBlock, prevBlock)
			}

			found, ok := foundByBlock[block]
			if !ok {
				found = make(map[ID]*logrec)
				foundByBlock[block] = found
			}

			if before, ok := found[rec.ID]; ok {
				rec = before
			} else {
				found[rec.ID] = rec
			}
			rec.matched++
			if rec.matched == syncing.PositionsCount() {
				gonext, err = onMatched(rec)
				if !gonext {
					syncing.Halt()
					return
				}
			}

			return
		}
	}

	// start the threads
	var preparing sync.WaitGroup
	preparing.Add(1)
	for pos := range pattern {
		if len(pattern[pos]) == 0 {
			continue
		}
		for i, variant := range pattern[pos] {
			syncing.StartThread(pos, i)
			go func(pos, i int, variant common.Hash) {
				onMatched := aggregator(pos, i)
				preparing.Wait()
				tt.scanPatternVariant(uint8(pos), variant, blockStart, onMatched, onDbIterator)
			}(pos, i, variant)
		}
	}
	preparing.Done()

	syncing.WaitForThreads()

	return ctx.Err()
}

func (tt *index) scanPatternVariant(pos uint8, variant common.Hash, start uint64, onMatched logHandler, onDbIterator func()) {
	prefix := append(variant.Bytes(), posToBytes(pos)...)

	onDbIterator()
	it := tt.table.Topic.NewIterator(prefix, uintToBytes(start))
	defer it.Release()
	for it.Next() {
		id := extractLogrecID(it.Key())
		topicCount := bytesToPos(it.Value())
		rec := newLogrec(id, topicCount)

		gonext, _ := onMatched(rec)
		if !gonext {
			break
		}
	}
	onMatched(nil)
}
