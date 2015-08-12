package main

import (
	"os"
	"runtime"
	"testing"

	_ "github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

func benchmarkPatriciaGet(b *testing.B, data *Data, t *trie.Trie) {
	keys := data.Keys
	runtime.GC()
	var r []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		r = t.Get(keys[im])
	}
	resultI = r
}

func benchmarkPatriciaSetRm(b *testing.B, data *Data, t *trie.Trie) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		im := i % data.N()
		t.Update(keys[im], values[im])
		t.Delete(keys[im])
	}
}

func benchmarkPatriciaHash(b *testing.B, data *Data, t *trie.Trie) {
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Hash()
		// TODO ?
	}
}

func benchmarkPatriciaUpdateHash(b *testing.B, data *Data, t *trie.Trie, updates int) {
	keys := data.Keys
	values := data.Values
	runtime.GC()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < updates; j++ {
			jm := (j + i) % data.N()
			t.Update(keys[jm], values[jm])
			t.Hash()
		}
		b.StopTimer()
		for j := 0; j < updates; j++ {
			jm := (j + i) % data.N()
			t.Delete(keys[jm])
			t.Hash()
		}
		b.StartTimer()
	}
}

//------------------------------------------------------------------------------------------
// 1) benchmark Patricia get

func runPatriciaGetBenchmarks() {
	logger.Println("\nPatriciaGet")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{x, y, z})
				tree := initPatriciaTree(data, nil)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkPatriciaGet(b, data, tree)
				})
				logger.Println(x, y, z, r)
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 2) benchmark Patricia get after save and load

func runPatriciaGetLoadBenchmarks() {
	logger.Println("\nPatriciaGetLoad")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)

					data := initData(&dataParams{x, y, z})
					dbData := initData(&dataParams{0, 0, 0})
					db := initLevelDBEth(dbData, c, "bench")
					tree := initPatriciaTree(data, db)
					tree.Commit()
					rootHash := tree.Root()
					db.Close()
					data = nil
					db = initLevelDBEth(dbData, c, "bench")
					tree = trie.New(rootHash, db)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkPatriciaGet(b, data, tree)
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
// 3) benchmark Patricia set+rm

func runPatriciaSetRmBenchmarks() {
	logger.Println("\nPatriciaSetRm")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				x, y, z := X(ix), Y(iy), Z(iz)
				data := initData(&dataParams{TestDataSize, y, z})
				preData := initData(&dataParams{x, y, z})
				tree := initPatriciaTree(preData, nil)
				r := testing.Benchmark(func(b *testing.B) {
					benchmarkPatriciaSetRm(b, data, tree)
				})
				logger.Println(x, y, z, r)
			}
		}
	}
}

//------------------------------------------------------------------------------------------
// 4) benchmark Patricia set+rm after save and load

func runPatriciaSetRmLoadBenchmarks() {
	logger.Println("\nPatriciaSetRmLoad")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
					dbData := initData(&dataParams{0, 0, 0})
					db := initLevelDBEth(dbData, c, "bench")
					data := initData(&dataParams{TestDataSize, y, z})
					preData := initData(&dataParams{x, y, z})
					tree := initPatriciaTree(preData, db)
					tree.Commit()
					rootHash := tree.Root()
					tree = trie.New(rootHash, db)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkPatriciaSetRm(b, data, tree)
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
// 5) benchmark Patricia compute merkle root

func runPatriciaHashBenchmarks() {
	logger.Println("\nPatriciaHash")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					x, y, z, c := X(ix), Y(iy), Z(iz), C(ic)
					dbData := initData(&dataParams{0, 0, 0})
					db := initLevelDBEth(dbData, c, "bench")
					data := initData(&dataParams{x, y, z})
					tree := initPatriciaTree(data, db)
					r := testing.Benchmark(func(b *testing.B) {
						benchmarkPatriciaHash(b, nil, tree)
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
// 6) benchmark Patricia update and compute merkle root

func runPatriciaUpdateHashBenchmarks() {
	logger.Println("\nPatriciaUpdateHash")
	for ix := 0; ix < numX; ix++ {
		for iy := 0; iy < numY; iy++ {
			for iz := 0; iz < numZ; iz++ {
				for ic := 0; ic < numC; ic++ {
					for iw := 0; iw < numW; iw++ {
						x, y, z, c, w := X(ix), Y(iy), Z(iz), C(ic), W(iw)
						dbData := initData(&dataParams{0, 0, 0})
						db := initLevelDBEth(dbData, c, "bench")
						data := initData(&dataParams{TestDataSize, y, z})
						preData := initData(&dataParams{x, y, z})
						tree := initPatriciaTree(preData, db)
						tree.Commit()
						rootHash := tree.Root()
						tree = trie.New(rootHash, db)
						r := testing.Benchmark(func(b *testing.B) {
							benchmarkPatriciaUpdateHash(b, data, tree, w)
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
