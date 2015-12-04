package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime/pprof"

	"github.com/eris-ltd/common/go/log"
)

var SAVE_RESULTS = false

func init() {
	saves := os.Getenv("GO_BENCHMARK_SAVE_RESULTS")
	if saves != "" {
		SAVE_RESULTS = true
	}

	go func() {
		logger.Println(http.ListenAndServe("localhost:6060", nil))
	}()

}

func openFileRunTest(f func(io.Writer), saveDir, p string) {
	var w *os.File
	var err error
	if saveDir != "" {
		w, err = os.Create(path.Join(saveDir, p))
		ifExit(err)
		defer w.Close()
	}
	f(w)
}

func main() {
	f, _ := os.Create("cpu_file")
	pprof.StartCPUProfile(f)

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "hash":
			var saveDir string
			if SAVE_RESULTS {
				if len(args) < 2 {
					exit("enter a directory to store the results")
				} else {
					saveDir = args[1]
				}
			}
			openFileRunTest(runRipemdHashBenchmarks, saveDir, "ripemd160")
			openFileRunTest(runSha256HashBenchmarks, saveDir, "sha256")
			openFileRunTest(runSha3HashBenchmarks, saveDir, "sha3")
		default:
			fmt.Println(args[0])
		}

	}

	//runCompactHexDecodeBenchmarks()
	//	runIAVLUpdateHashBenchmarks()

	//runIAVLHashBenchmarks()

	//	runIAVLGetBenchmarks()
	//	runIAVLGetLoadBenchmarks()

	//	runPatriciaGetBenchmarks()
	//runPatriciaGetLoadBenchmarks()

	//	runPatriciaGetLoadMemStats()
	//	runIAVLGetLoadMemStats()

	runIAVLSetRmBenchmarks()
	//	runPatriciaSetRmBenchmarks()
	//runIAVLSetRmLoadBenchmarks()
	//runLevelDBSetRmBenchmarks()
	//	runMemDBSetRmBenchmarks()

	//runCopyByteSliceBenchmarks()
	//runWriteByteSliceBenchmarks()
	//runBasicCodecEncodeBenchmarks()

	//	runRandBytesBenchmarks()

	log.Flush()
	pprof.StopCPUProfile()

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

func exit(s string) {
	fmt.Println(s)
	log.Flush()
	os.Exit(1)

}

func ifExit(err error) {
	if err != nil {
		fmt.Println(err)
		log.Flush()
		os.Exit(1)
	}
}
