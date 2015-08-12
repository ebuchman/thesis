package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"

	. "github.com/tendermint/tendermint/common"

	"code.google.com/p/go.crypto/ripemd160"
	"github.com/tendermint/tendermint/vm/sha3"
	"github.com/tendermint/tendermint/wire"

	_ "github.com/ethereum/go-ethereum/trie"
)

var result []byte // so the compiler doesnt get away with optimizing code away

//------------------------------------------------------------------------------------------
// benchmark functions that take a size

func runBytesBenchmarksFunc(f func(b *testing.B, size int), name string) {
	fmt.Println("\nBenchmarking", name)
	for _, size := range []int{20, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192} {
		r := testing.Benchmark(func(b *testing.B) {
			f(b, size)
		})
		fmt.Println(size, r)
	}
}

//------------------------------------------------------------------------------------------
// benchmark copy byte slice

func runCopyByteSliceBenchmarks() {
	runBytesBenchmarksFunc(benchmarkCopyByteSlice, "CopyByteSlice")
}

func benchmarkCopyByteSlice(b *testing.B, size int) {
	var v = RandBytes(size)
	var h = make([]byte, size)
	for i := 0; i < b.N; i++ {
		copy(h, v)
	}
}

//------------------------------------------------------------------------------------------
// benchmark write byte slice to buffer

func runWriteByteSliceBenchmarks() {
	runBytesBenchmarksFunc(benchmarkWriteByteSlice, "WriteByteSlice")
}

func benchmarkWriteByteSlice(b *testing.B, size int) {
	var v = RandBytes(size)
	var w = new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		w.Reset() // takes about 8 ns
		w.Write(v)
	}
}

//------------------------------------------------------------------------------------------
// benchmark BasicCodec.Encode on byte slice

func runBasicCodecEncodeBenchmarks() {
	runBytesBenchmarksFunc(benchmarkBasicCodecEncode, "BasicCodecEncode")
}

func benchmarkBasicCodecEncode(b *testing.B, size int) {
	var v = RandBytes(size)
	var w, n, err = new(bytes.Buffer), new(int64), new(error)
	for i := 0; i < b.N; i++ {
		w.Reset() // takes about 8 ns
		wire.BasicCodec.Encode(v, w, n, err)
	}
}

//------------------------------------------------------------------------------------------
// benchmark RandBytes

func runRandBytesBenchmarks() {
	runBytesBenchmarksFunc(benchmarkRandBytes, "RandBytes")
}

func benchmarkRandBytes(b *testing.B, size int) {
	var v []byte
	for i := 0; i < b.N; i++ {
		v = RandBytes(size)
	}
	result = v
}

//------------------------------------------------------------------------------------------
// benchmark ripemd160

func runRipemdHashBenchmarks() {
	runBytesBenchmarksFunc(benchmarkRipemd160, "Ripemd160")
}

var hashResult []byte

func benchmarkRipemd160(b *testing.B, size int) {

	v := RandBytes(size)
	var hash []byte
	for i := 0; i < b.N; i++ {
		hasher := ripemd160.New()
		hasher.Write(v)
		hash = hasher.Sum(nil)
		// using binary instead of raw costs us about 1000 ns/op
		// hash = wire.BinaryRipemd160(v)
	}
	hashResult = hash
}

//------------------------------------------------------------------------------------------
// benchmark sha2

func runSha256HashBenchmarks() {
	runBytesBenchmarksFunc(benchmarkSha256, "Sha256")
}

func benchmarkSha256(b *testing.B, size int) {

	v := RandBytes(size)
	var hash []byte
	for i := 0; i < b.N; i++ {
		hasher := sha256.New()
		hasher.Write(v)
		hash = hasher.Sum(nil)
		// using binary instead of raw costs us about 1000 ns/op
		// hash = wire.BinarySha256(v)
	}
	hashResult = hash
}

//------------------------------------------------------------------------------------------
// benchmark sha3

func runSha3HashBenchmarks() {
	runBytesBenchmarksFunc(benchmarkSha3, "Sha3")
}

func benchmarkSha3(b *testing.B, size int) {

	v := RandBytes(size)
	var hash []byte
	for i := 0; i < b.N; i++ {
		hash = sha3.Sha3(v)
	}
	hashResult = hash
}

//------------------------------------------------------------------------------------------
// benchmark CompactHexDecode

func runCompactHexDecodeBenchmarks() {
	runBytesBenchmarksFunc(benchmarkCompactHexDecode, "CompactHexDecode")
}

func benchmarkCompactHexDecode(b *testing.B, size int) {
	//var v = RandBytes(size)
	for i := 0; i < b.N; i++ {
		//trie.CompactHexDecode(v)
	}
}
