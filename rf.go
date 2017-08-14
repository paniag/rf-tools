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

func framer(out *os.File, ch chan []byte) {
	for data := range ch {
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

			out.WriteString(fmt.Sprintf("Frame %v : %v bytes\n", atomicFrameId, len(frame)))

			out.Write(frame)

			padLength := MAX_PAYLOAD_SIZE - len(frame)
			out.Write(pad[:padLength])
			out.WriteString("\n")
			atomic.AddUint64(&atomicFrameId, 1)

			start = end
			end += MAX_PAYLOAD_SIZE
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := os.Stdout
	ch := make(chan []byte)
	var arr [INPUT_BUFFER_SIZE]byte

	// More idiomatic to populate an array. Specifying the capacity reserves the
	//  memory upfront, so it's literally 1 allocation rather than
	//  "amortized O(1)".
	pad = make([]byte, 0, MAX_PAYLOAD_SIZE)
	for i := 0; i < MAX_PAYLOAD_SIZE; i++ {
		// NB 48 ASCII Zero, placeholder for 0:NUL
		pad = append(pad, byte(48))
	}

	go framer(out, ch)

loop:
	for {
		buf := arr[:]
		fmt.Printf("len(buf) = %v\n", len(buf))
		fmt.Printf("cap(buf) = %v\n", cap(buf))
		n, err := in.Read(buf)
		fmt.Printf("%v, %#v = in.Read(buf)\n", n, err)
		// Handle errors as close as possible to their origin.
		// Use switch instead of if/else if chains.
		switch err {
		case io.EOF:
			break loop
		case nil:
			buf = buf[:n]
			if n > 0 {
				ch <- buf
			}
		default:
			glog.Fatal(err)
		}
	}

	fmt.Println("Done.")

	os.Exit(0)
}
