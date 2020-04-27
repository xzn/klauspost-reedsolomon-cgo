package main

import "C"

import (
	reedsolomon "github.com/klauspost/reedsolomon"
)

type encoder_param struct {
	data_n   int
	parity_n int
}

var encoders = make(map[encoder_param]reedsolomon.Encoder)

func get(data_n int, parity_n int) (reedsolomon.Encoder, error) {
	param := encoder_param{data_n: data_n, parity_n: parity_n}
	enc := encoders[param]
	if enc != nil {
		return enc, nil
	}
	enc, err := reedsolomon.New(data_n, parity_n)
	if err != nil {
		return nil, err
	}
	encoders[param] = enc
	return enc, nil
}

//export Encode
func Encode(data [][]byte, data_n int) int {
	parity_n := len(data) - data_n
	enc, err := get(data_n, parity_n)
	if err != nil {
		return -1
	}
	err = enc.Encode(data)
	if err != nil {
		return -1
	}
	return 0
}

//export Decode
func Decode(data [][]byte, data_n int, nils []bool) int {
	total_n := len(data)
	parity_n := total_n - data_n
	enc, err := get(data_n, parity_n)
	if err != nil {
		return -1
	}
	d := make([][]byte, total_n)
	copy(d, data)
	for i := 0; i < total_n; i++ {
		if nils[i] {
			d[i] = nil
		}
	}
	err = enc.ReconstructData(d)
	if err != nil {
		return -1
	}
	for i := 0; i < total_n; i++ {
		if nils[i] {
			copy(data[i], d[i])
		}
	}
	return 0
}

func main() {
}
