package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/eris-ltd/common/go/log"
)

func init() {

	go func() {
		logger.Println(http.ListenAndServe("localhost:6060", nil))
	}()

}

func main() {

	//runCompactHexDecodeBenchmarks()
	//	runIAVLUpdateHashBenchmarks()

	//runIAVLHashBenchmarks()

	runIAVLGetBenchmarks()
	//	runIAVLGetLoadBenchmarks()

	runPatriciaGetBenchmarks()
	//runPatriciaGetLoadBenchmarks()

	//	runPatriciaGetLoadMemStats()
	//	runIAVLGetLoadMemStats()

	runIAVLSetRmBenchmarks()
	runPatriciaSetRmBenchmarks()
	//runIAVLSetRmLoadBenchmarks()
	//runLevelDBSetRmBenchmarks()
	//runMemDBSetRmBenchmarks()

	//runCopyByteSliceBenchmarks()
	//runWriteByteSliceBenchmarks()
	//runBasicCodecEncodeBenchmarks()

	//	runRandBytesBenchmarks()

	/*
		runRipemdHashBenchmarks()
		runSha256HashBenchmarks()
		runSha3HashBenchmarks()
	*/
	log.Flush()
}

/*
LevelDBGet
1 32 256  1000000             1206 ns/op
10 32 256  1000000            1423 ns/op
100 32 256  1000000           1832 ns/op
1000 32 256  1000000          2129 ns/op
10000 32 256   500000         3189 ns/op
16384 32 256   200000         7565 ns/op
25000 32 256   300000         7368 ns/op
32768 32 256   100000        10798 ns/op
50000 32 256   200000        11100 ns/op
65536 32 256    30000        49901 ns/op
75000 32 256   100000        13356 ns/op
100000 32 256    10000      138344 ns/op
250000 32 256    10000      180822 ns/op
500000 32 256    10000      249448 ns/op
LevelDBEthGet
1 32 256  1000000             1410 ns/op
10 32 256  1000000            1608 ns/op
100 32 256  1000000           2090 ns/op
1000 32 256   500000          2273 ns/op
10000 32 256   500000         3300 ns/op
16384 32 256   200000         7974 ns/op
25000 32 256   200000         6248 ns/op
32768 32 256   200000         8949 ns/op
50000 32 256   100000        10920 ns/op
65536 32 256   100000        13232 ns/op
75000 32 256    20000        65056 ns/op
100000 32 256    10000      128287 ns/op
250000 32 256    10000      198661 ns/op
500000 32 256    10000      217917 ns/op
*/
func dbBenchmarks() {
	runLevelDBGetBenchmarks()
	runLevelDBEthGetBenchmarks()
}

////////
/*

From running 1 Get query on a DB with 100 entries

IAVLGet
100 32 256 <nil> 1.421µs

IAVLGetLoad
LevelDBMint GET
LevelDBMint GET
LevelDBMint GET
LevelDBMint GET
LevelDBMint GET
LevelDBMint GET
LevelDBMint GET
100 32 256 <nil> 57.919µs

PatriciaGet
100 32 256 [] 20.817µs

PatriciaGetLoad
LevelDBEth GET
LevelDBEth GET
LevelDBEth GET
100 32 256 0 [] 124.788µs

*/
