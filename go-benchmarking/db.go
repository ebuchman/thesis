package main

import (
	"os"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
	. "github.com/tendermint/go-common"
	dbm "github.com/tendermint/go-db"
	"github.com/tendermint/go-merkle"
	"github.com/tendermint/go-wire"
)

//------------------------------------------------------------------------------------------
// data to fill trees with

var TestDataSize = 10000

type Data struct {
	n      int
	Keys   [][]byte
	Values [][]byte
}

func (d *Data) N() int {
	return d.n
}

// TODO: we should make missData so we can bench misses too
func initData(dp *dataParams) *Data {
	keys := make([][]byte, dp.N())
	for i := 0; i < dp.N(); i++ {
		keys[i] = RandBytes(dp.KeyLength())
	}

	values := make([][]byte, dp.N())
	for i := 0; i < dp.N(); i++ {
		values[i] = RandBytes(dp.ValueLength())
	}

	return &Data{dp.N(), keys, values}
}

var emptyData = initData(&dataParams{0, 0, 0})

//------------------------------------------------------------------------------------------
// params for making data
// TODO: add distributions over lengths

type dataParams struct {
	numElem     int
	keyLength   int
	valueLength int
}

func (dp *dataParams) N() int {
	return dp.numElem
}

func (dp *dataParams) KeyLength() int {
	return dp.keyLength
}

func (dp *dataParams) ValueLength() int {
	return dp.valueLength
}

//------------------------------------------------------------------------------------------
// init funcs (optionally fill with numElem items)

func initLevelDB(data *Data, name string) *dbm.LevelDB {
	db, err := dbm.NewLevelDB(name)
	if err != nil {
		panic(err)
	}
	for i := 0; i < data.N(); i++ {
		db.Set(data.Keys[i], data.Values[i])
	}
	return db
}

func initLevelDBEth(data *Data, cacheSize int, name string) *ethdb.LDBDatabase {
	db, err := ethdb.NewLDBDatabase(name, cacheSize)
	if err != nil {
		panic(err)
	}
	for i := 0; i < data.N(); i++ {
		db.Put(data.Keys[i], data.Values[i])
	}
	return db
}

func initMemDB(data *Data) *dbm.MemDB {
	db := dbm.NewMemDB()
	for i := 0; i < data.N(); i++ {
		db.Set(data.Keys[i], data.Values[i])
	}
	return db
}

func initIAVLTree(data *Data, cacheSize int, db dbm.DB) *merkle.IAVLTree {
	t := merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, cacheSize, db)
	for i := 0; i < data.N(); i++ {
		t.Set(data.Keys[i], data.Values[i])
	}
	return t
}

func initPatriciaTree(data *Data, db *ethdb.LDBDatabase) *trie.Trie {
	var t *trie.Trie
	// no error for empty root
	if db == nil { // or we'll be bit by the nil dog! TODO: shouldnt be true
		t, _ = trie.New(common.Hash{}, nil)
	} else {
		t, _ = trie.New(common.Hash{}, db)
	}
	for i := 0; i < data.N(); i++ {
		t.Update(data.Keys[i], data.Values[i])
	}
	return t
}

//------------------------------------------------------------------------------------------
// core benchmark funcs

func benchmarkLevelDBGet(b *testing.B, data *Data, db *dbm.LevelDB) {
	keys := data.Keys
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		db.Get(keys[im])
	}
}

func benchmarkLevelDBSetRm(b *testing.B, data *Data, db *dbm.LevelDB) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % TestDataSize
		db.Set(keys[im], values[im])
		db.Delete(keys[im])
	}
}

func benchmarkLevelDBEthGet(b *testing.B, data *Data, db *ethdb.LDBDatabase) {
	keys := data.Keys
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		db.Get(keys[im])
	}
}

func benchmarkLevelDBEthSetRm(b *testing.B, data *Data, db *ethdb.LDBDatabase) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % TestDataSize
		db.Put(keys[im], values[im])
		db.Delete(keys[im])
	}
}

func benchmarkMemDBSetRm(b *testing.B, data *Data, db *dbm.MemDB) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % TestDataSize
		db.Set(keys[im], values[im])
		db.Delete(keys[im])
	}
}

//------------------------------------------------------------------------------------------
// benchmark LevelDB get

func runLevelDBGetBenchmarks() {
	logger.Println("LevelDBGet")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{x, y, z})
				db := initLevelDB(data, "bench")
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkLevelDBGet(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// benchmark LevelDB set+rm

func runLevelDBSetRmBenchmarks() {
	logger.Println("LevelDBSetRm")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{TestDataSize, y, z})
				preData := initData(&dataParams{x, y, z})
				db := initLevelDB(preData, "bench")
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkLevelDBSetRm(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// benchmark LevelDBEth get

func runLevelDBEthGetBenchmarks() {
	logger.Println("LevelDBEthGet")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{x, y, z})
				db := initLevelDBEth(data, 0, "bench")
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkLevelDBEthGet(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// benchmark LevelDBEth set+rm

func runLevelDBEthSetRmBenchmarks() {
	logger.Println("LevelDBEthSetRm")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{TestDataSize, y, z})
				preData := initData(&dataParams{x, y, z})
				db := initLevelDBEth(preData, 0, "bench")
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkLevelDBEthSetRm(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// benchmark MemDB get

func runMemDBGetBenchmarks() {
	logger.Println("MemDBGet")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{TestDataSize, y, z})
				db := initMemDB(data)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkMemDBSetRm(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// benchmark MemDB set+rm

func runMemDBSetRmBenchmarks() {
	logger.Println("MemDBSetRm")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{TestDataSize, y, z})
				preData := initData(&dataParams{x, y, z})
				db := initMemDB(preData)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkMemDBSetRm(b, data, db)
				})
				logger.Println(x, y, z, r)
				db.Close()
				os.RemoveAll("bench")
			}
		}
	}
}
