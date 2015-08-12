package main

import ()

// fill tree with this many itmes
var numX = len(x)

//var x = []int{0, 10, 100, 1000, 10000, 100000, 1000000, 10000000}
//var x = []int{100, 1000, 10000, 25000, 50000, 100000, 500000, 1000000}
var x = []int{1000, 10000, 100000, 1000000} //, 10, 100, 1000, 10000} //, 16384, 25000, 32768, 50000, 65536, 75000, 100000, 250000, 500000}

func X(i int) int {
	// no check!
	return x[i]
}

// keys of this length
var numY = 1 // + whatever distributions
// var y = []int{20, 32, 64, 128, 256, 512, 1024}
var y = []int{32}

func Y(i int) int {
	if i < 7 {
		return y[i]
	} else {
		// TODO: return distributions
	}
	return -1
}

// values of this size
var numZ = 1 // + whatever distributions
// var z = []int{20, 32, 64, 128, 256, 512, 1024}
var z = []int{256}

func Z(i int) int {
	if i < 7 {
		return z[i]
	} else {
		// TODO: return distributions
	}
	return -1
}

var numC = 1

//var c = []int{0, 1000, 10000, 100000, 1000000}
var c = []int{0, 1000, 1000000}

func C(i int) int {
	return c[i]
}

var numW = 5

var w = []int{1, 10, 100, 1000, 10000}

func W(i int) int {
	return w[i]
}
