package main

import (
	"os"
	"runtime"
	"testing"

	"github.com/tendermint/tendermint/merkle"
	"github.com/tendermint/tendermint/wire"
)

//------------------------------------------------------------------------------------------
// core benchmark functions

func benchmarkIAVLGet(b *testing.B, data *Data, t *merkle.IAVLTree) {
	keys := data.Keys
	runtime.GC()
	var r interface{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		_, r = t.Get(keys[im])
	}
	resultI = r
}

func benchmarkIAVLSetRm(b *testing.B, data *Data, t *merkle.IAVLTree) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		t.Set(keys[im], values[im])
		t.Remove(keys[im])
	}
}

func benchmarkIAVLHash(b *testing.B, data *Data, t *merkle.IAVLTree) {
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// XXX: t.HashNoSet()
	}
}

func benchmarkIAVLUpdateHash(b *testing.B, data *Data, t *merkle.IAVLTree, updates int) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < updates; j++ {
			jm := (j + i) % data.N()
			t.Set(keys[jm], values[jm])
			t.Hash()
		}
		b.StopTimer()
		for j := 0; j < updates; j++ {
			jm := (j + i) % data.N()
			t.Remove(keys[jm])
			t.Hash()
		}
		b.StartTimer()

	}
}

//------------------------------------------------------------------------------------------
// 1) benchmark IAVL get

func runIAVLGetBenchmarks() {
	logger.Println("\nIAVLGet")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				ic := 0 // cache size 0
				x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
				data := initData(&dataParams{x, y, z})
				tree := initIAVLTree(data, c, nil)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkIAVLGet(b, data, tree)
				})
				logger.Println(x, y, z, r)
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 2) benchmark IAVL get after save and load

func runIAVLGetLoadBenchmarks() {
	logger.Println("\nIAVLGetLoad")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					preData := initData(&dataParams{0, 0, 0})
					db := initLevelDB(preData, "bench")
					//db := initMemDB(0, 0, 0)
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
					data := initData(&dataParams{x, y, z})
					tree := initIAVLTree(data, c, db)
					rootHash := tree.Save()
					tree = merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, c, db)
					tree.Load(rootHash)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkIAVLGet(b, data, tree)
					})
					logger.Println(x, y, z, r)
					db.Close()
					os.RemoveAll("bench")
				}
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 3) benchmark IAVL set+rm

func runIAVLSetRmBenchmarks() {
	logger.Println("\nIAVLSetRm")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				ic := 0 // cache size 0
				x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
				data := initData(&dataParams{TestDataSize, y, z})
				preData := initData(&dataParams{x, y, z})
				tree := initIAVLTree(preData, c, nil)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkIAVLSetRm(b, data, tree)
				})
				logger.Println(x, y, z, r)
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 4) benchmark IAVL set+rm after save and load

func runIAVLSetRmLoadBenchmarks() {
	logger.Println("\nIAVLSetRmLoad")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					dbData := initData(&dataParams{0, 0, 0})
					db := initLevelDB(dbData, "bench")
					//db := initMemDB(0, 0, 0)
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
					data := initData(&dataParams{TestDataSize, y, z})
					preData := initData(&dataParams{x, y, z})
					tree := initIAVLTree(preData, c, db)
					rootHash := tree.Save()
					tree = merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, c, db)
					tree.Load(rootHash)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkIAVLSetRm(b, data, tree)
					})
					logger.Println(x, y, z, c, r)
					db.Close()
					os.RemoveAll("bench")
				}
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 5) benchmark IAVL compute merkle root

func runIAVLHashBenchmarks() {
	logger.Println("\nIAVLHash")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					dbData := initData(&dataParams{0, 0, 0})
					db := initLevelDB(dbData, "bench")
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
					data := initData(&dataParams{x, y, z})
					tree := initIAVLTree(data, c, db)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkIAVLHash(b, nil, tree)
					})
					logger.Println(x, y, z, c, r)
					db.Close()
					os.RemoveAll("bench")
				}
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 6) benchmark IAVL update and compute merkle root

func runIAVLUpdateHashBenchmarks() {
	logger.Println("\nIAVLUpdateHash")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					for iw := 0; iw < numW; iw++ {
						dbData := initData(&dataParams{0, 0, 0})
						db := initLevelDB(dbData, "bench")
						x, y, z, c, w := X(ix), Y(iy), Z(iz), C(ic), W(iw)
						data := initData(&dataParams{TestDataSize, y, z})
						preData := initData(&dataParams{x, y, z})
						tree := initIAVLTree(preData, c, db)
						rootHash := tree.Save()
						tree = merkle.NewIAVLTree(wire.BasicCodec, wire.BasicCodec, c, db)
						tree.Load(rootHash)
						r := testing.Benchmark(func(b *testing.B) {
							benchmarkIAVLUpdateHash(b, data, tree, w)
						})
						logger.Println(x, y, z, c, w, r)
						db.Close()
						os.RemoveAll("bench")
					}
				}
			}
		}
	}
}

// for forcing results
// TODO: do we really need it?
var resultI interface{}
