// rf.go
// Copyright 2017 Mac Radigan
// All Rights Reserved

  package main

  import (
    "bufio"
    "io"
    "math"
    "os"
    "strconv"
    "sync/atomic"
  )

  var frame_id uint64 = 0 // atomic frame counter

  /// "echo" function, writes a byte channel to a file
  func display(out *os.File, ch chan []byte) {
    
    for char := range ch {
      out.Write(char)
    }
    
  }

  /// data framer, chunks a channel into PAYLOAD_SIZE frames
  func framer(out *os.File, ch chan []byte) {
    
    const PAYLOAD_SIZE int = 1024                    // length of data frame payload
    pad := make([]byte, PAYLOAD_SIZE, PAYLOAD_SIZE)  // padding source
    for k, _ := range pad {
      pad[k] = byte(48) // NB 48 ASCII Zero, placeholder for 0:NUL
    }
    
    /// read from channel, chunk into PAYLOAD_SIZE frames
    for data := range ch {
      /// required number of frames
      n_frames := int(math.Max( 1.0, math.Ceil(float64(len(data))/float64(PAYLOAD_SIZE)) )) 
      for n := 0; n < n_frames; n++ { // each frame
        /// starting / ending index of frame ( k_0 : k_1 )
        k_0 := n * PAYLOAD_SIZE
        k_1 := int(math.Min( float64(k_0+PAYLOAD_SIZE), float64(len(data)) ))
        /// placeholder for frame header
        out.WriteString("Frame " + strconv.FormatUint(frame_id, 10) + " : " + strconv.Itoa(k_1-k_0) + " bytes\n")
        /// write payload data
        out.Write(data[k_0:k_1]) // payload data
        n_pad := int(math.Mod(float64(k_1-k_0), float64(PAYLOAD_SIZE)) )
        out.Write(pad[0:n_pad])  // residual frame padding
        out.WriteString("\n")
        atomic.AddUint64(&frame_id, 1)
      }
    
    }
  }

  func main() {
    
    const BUF_SIZE int = 2048          // input buffer length
    in  := bufio.NewReader(os.Stdin)
    out := os.Stdout
    ch  := make(chan []byte)
  //go display(out, ch)
    go framer(out, ch)
    
    /// read stdin to buffer, send to downstream processing
    for {
      buf := make([]byte, BUF_SIZE)
      n, err := in.Read(buf)
      buf = buf[:n]
      if n > 0 { ch <- buf }
      if err == io.EOF {
        break
      } else if err != nil {
        panic(err)
      }
    }
    
    os.Exit(0)
    
  }

// *EOF*
