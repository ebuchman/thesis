package main

import (
	"os"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"

	//dbm "github.com/tendermint/tendermint/db"
	"github.com/tendermint/tendermint/merkle"
	"github.com/tendermint/tendermint/wire"
)

var MemStatsN = 100000
var MemStatsUnique = 100

var memStats = new(runtime.MemStats)

func reportMem(s string) {
	runtime.GC()
	runtime.ReadMemStats(memStats)
	logger.Println(s, float64(memStats.Alloc)/float64(1000000), "MB")
}

func reportBenchmark(duration time.Duration, count int) {
	logger.Printf("%f us/op, %d db hits\n", float64(duration.Nanoseconds()/int64(MemStatsN))/float64(1000), count)
}

/*
Notes
- a full run through the IAVL generates about half as much garbage as the patricia

*/

/*

NOTE: we generate 10000 data items with 32 byte keys and size 256 byte values
thats 32+256 say size, cap, pointer for each byte slice, so + 12 + 12 = 312
= 3.12 MB,
but we get an increase below of  3.58-0.21 = 3.37
How to account for the 0.25 MB ?

The data is only about 3 MB, but filling the trie adds more than twice that.
Part of this is because the length of the keys are doubled by CompactDecode (there's no uint4, alas)

After saving and clearing we only have the keys left + runtime, a respectable 0.78
Fetching items brings this back over 4MB - its cachining nodes in the tree, implicitly, even if you disable the cache.store
In tendermint this doesn't happen

PatriciaGetLoad
before generating data we have 0.215208 MB
after generating data we have 3.587456 MB
after initing leveldb we have 7.90152 MB
initializing tree 10000
after populating the tree we have 15.690408 MB
after clearing and running GC we have 0.789632 MB
after fetching from the tree we have 4.603048 MB
929 ns/op
10000 32 256 0
*/

func runPatriciaGetLoadMemStats() {
	logger.Println("\nPatriciaGetLoadMemStats")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)

					reportMem("before generating data we have")
					data := initData(&dataParams{x, y, z})
					reportMem("after generating data we have")

					db := initLevelDBEth(emptyData, c, "bench")
					reportMem("after initing leveldb we have")

					tree := initPatriciaTree(data, db)
					reportMem("after populating the tree we have")

					rootHash := tree.Root()
					reportMem("after hashing the tree we have")
					tree.Commit()
					reportMem("after committing the tree we have")
					db.Close()
					keys := data.Keys
					data = nil
					reportMem("after clearing we have")
					db = initLevelDBEth(emptyData, c, "bench")
					reportMem("after reloading the db we have")
					tree = trie.New(rootHash, db)
					reportMem("after reloading the tree we have")
					time.Sleep(time.Second * 2)
					reportMem("after sleeping for 2 seconds we have")

					startTime := time.Now()
					for i := 0; i < MemStatsN; i++ {
						im := i % MemStatsUnique
						tree.Get(keys[im])
					}
					duration := time.Since(startTime)
					reportMem("after fetching from the tree we have")

					reportBenchmark(duration, ethdb.COUNT)
					logger.Println(x, y, z, c)
					db.Close()
					os.RemoveAll("bench")
				}
			}
		}
	}
}

func runIAVLGetLoadMemStats() {
	logger.Println("\nIAVLGetLoadMemStats")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)

					reportMem("before generating data we have")
					data := initData(&dataParams{x, y, z})
					reportMem("after generating data we have")

					db := initLevelDB(emptyData, "bench")
					reportMem("after initing leveldb we have")

					tree := initIAVLTree(data, 10000, db)
					reportMem("after populating the tree we have")

					rootHash := tree.Hash()
					reportMem("after hashing the tree we have")
					tree.Save()
					reportMem("after committing the tree we have")
					db.Close()
					keys := data.Keys
					data = nil
					reportMem("after clearing we have")
					db = initLevelDB(emptyData, "bench")
					reportMem("after reloading leveldb we have")
					tree = merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, 50000, db)
					tree.Load(rootHash)
					reportMem("after reloading the tree we have")

					startTime := time.Now()
					for i := 0; i < MemStatsN; i++ {
						im := i % MemStatsUnique
						tree.Get(keys[im])
					}
					duration := time.Since(startTime)
					reportMem("after fetching from the tree we have")

					reportBenchmark(duration, -1) //dbm.COUNT)
					logger.Println(x, y, z, c)
					db.Close()
					os.RemoveAll("bench")
				}
			}
		}
	}
}
