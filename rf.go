/*
	rf.go
	Copyright 2017 Mac Radigan
	All Rights Reserved
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/golang/glog"
)

const (
	INPUT_BUFFER_SIZE = 2048
	MAX_PAYLOAD_SIZE  = 1024
)

var (
	atomicFrameId uint64
	pad           []byte
)

func sampler(chIn chan []byte, chOut chan []byte) {
	const SIZE = MAX_PAYLOAD_SIZE
	var arr [SIZE]byte
	var index = 0
	for data := range chIn {
		for byteIdx := 0; byteIdx < len(data); byteIdx++ {
			for bitIdx := 0; bitIdx < 8; bitIdx++ {
				if(data[byteIdx] & 1<<uint8(bitIdx) == 0) {
					arr[index] = byte(48)
				} else {
					arr[index] = byte(49)
				}
				index++
				if(index == SIZE) {
					chOut <- arr[:index]
					index = 0
				}
			}
		}
	}
	chOut <- arr[:index]
}

func framer(chIn chan []byte, chOut chan []byte) {
	for data := range chIn {
		start, end := 0, MAX_PAYLOAD_SIZE
		numFrames := len(data) / MAX_PAYLOAD_SIZE
		if len(data)%MAX_PAYLOAD_SIZE > 0 {
			numFrames++
		}
		for frameIdx := 0; frameIdx < numFrames; frameIdx++ {
			if len(data) < end {
				end = len(data)
			}
			frame := data[start:end]
			chOut <- []byte(fmt.Sprintf("Frame %v : %v bytes\n", atomicFrameId, len(frame)))
			chOut <- frame
			padLength := MAX_PAYLOAD_SIZE - len(frame)
			chOut <- pad[:padLength]
			chOut <- []byte("\n")
			atomic.AddUint64(&atomicFrameId, 1)
			start = end
			end += MAX_PAYLOAD_SIZE
		}
	}
}

func printer(out *os.File, chIn chan []byte) {
	for data := range chIn {
		out.Write(data)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := os.Stdout
	chIn := make(chan []byte)
	chFramer := make(chan []byte)
	chPrinter := make(chan []byte)
	var arr [INPUT_BUFFER_SIZE]byte

	// More idiomatic to populate an array. Specifying the capacity reserves the
	//  memory upfront, so it's literally 1 allocation rather than
	//  "amortized O(1)".
	pad = make([]byte, 0, MAX_PAYLOAD_SIZE)
	for i := 0; i < MAX_PAYLOAD_SIZE; i++ {
		// NB 48 ASCII Zero, placeholder for 0:NUL
		pad = append(pad, byte(48))
	}

	go sampler(chIn, chFramer)
	go framer(chFramer, chPrinter)
	go printer(out, chPrinter)

loop:
	for {
		buf := arr[:]
		n, err := in.Read(buf)
		// Handle errors as close as possible to their origin.
		// Use switch instead of if/else if chains.
		switch err {
		case io.EOF:
			break loop
		case nil:
			buf = buf[:n]
			if n > 0 {
				chIn <- buf
			}
		default:
			glog.Fatal(err)
		}
	}

	fmt.Println("Done.")

	os.Exit(0)
}
