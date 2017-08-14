/*
	rf.go
	Copyright 2017 Mac Radigan
	All Rights Reserved
*/
// Package comments should use /* */ style and describe the purpose and usage
//  of the package.

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	// glog provies much friendlier, more flexible logging than "log",
	//  which in turn is a lot better than fmt.Print*.
	"github.com/golang/glog"
)

// +1 Constants do use UPPER_SNAKE_CASE, like the convention for C macros.
// Type annotations are usually optional on numerical constants.
const (
	INPUT_BUFFER_SIZE = 2048
	MAX_PAYLOAD_SIZE  = 1024
)

// The initialization was unnecessary. Go variables without explicit
// initializers are always initialized with the zero value of their type.
var frameId uint64 // atomic frame counter // Use camel case for variable names rather than underscores.

// See http://godoc.org/io/ioutil#WriteFile

// Do 3 slashes have some special meaning I'm unaware of? Normally comments get 2 slashes.
/// data framer, chunks a channel into MAX_PAYLOAD_SIZE frames
func framer(out *os.File, ch chan []byte) {
	// There shouldn't be a blank line at the start of a block.
	// Comments should be sentences where possible, and normally begin with a
	//  capital letter and end with some punctuation, eg a period.
	// Choose variable names that are concise but self-explanatory.
	// Avoid redundant comments.
	pad := make([]byte, MAX_PAYLOAD_SIZE, MAX_PAYLOAD_SIZE)
	// i, j, and k are all idiomatic names for integer indices, in that order of
	//  preference. k often means "key". You can view a slice index as a key, but
	//  it may be misleading to readers.
	// You don't need the ", _"; you can just use the single value form.
	for i := range pad {
		// NB 48 ASCII Zero, placeholder for 0:NUL
		// But why?
		pad[i] = byte(48)
	}

	/// read from channel, chunk into MAX_PAYLOAD_SIZE frames
	// +1 Ranging over a channel is idiomatic Go.
	for data := range ch {
		// Use descriptive variable names.
		// Yikes, that's a lot of nested conversions and calls.
		numFrames := len(data) / MAX_PAYLOAD_SIZE
		if len(data)%MAX_PAYLOAD_SIZE > 0 {
			numFrames++
		}

		for frameIdx := 0; frameIdx < numFrames; frameIdx++ {
			// Go for simplicity, including simpler expressions.
			start := frameIdx * MAX_PAYLOAD_SIZE
			end := (frameIdx + 1) * MAX_PAYLOAD_SIZE
			if len(data) < end {
				end = len(data)
			}

			// Name that frame! :)
			// Note that this is *not* an allocation. Well, not of len(frame) bytes anyway.
			// It's actually a reference into a subrange of the underlying array.
			frame := data[start:end]

			// The fmt package is standard and shorter.
			out.WriteString(fmt.Sprintf("Frame %v : %v bytes\n", frame_id, len(frame)))

			out.Write(frame)

			padLength := MAX_PAYLOAD_SIZE - len(frame)
			// Slices starting from index 0 or ending at the end of the original slice
			// or array do not require that the obvious bound be specified.
			out.Write(pad[:padLength])
			out.WriteString("\n")
			atomic.AddUint64(&frame_id, 1)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := os.Stdout
	ch := make(chan []byte)
	// Since you have a 0-capacity channel, you only need to allocate buf once.
	buf := make([]byte, INPUT_BUFFER_SIZE)

	go framer(out, ch)

	/// read stdin to buffer, send to downstream processing
out:
	for {
		buf = buf[:0]
		n, err := in.Read(buf)
		// Handle errors as close as possible to their origin.
		// Use switch instead of if/else if chains.
		switch err {
		case io.EOF:
			break
		case nil:
			buf = buf[:n]
			if n > 0 {
				ch <- buf
			}
		default:
			// DON'T PANIC
			// Do realize this will crash with a stacktrace on any error returned.
			// The stack trace part might not be so useful since it's a standard
			//  library function
			glog.Fatal(err)
		}
	}

	os.Exit(0)
}

// *EOF*
