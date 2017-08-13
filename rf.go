// rf.go
// Copyright 2017 Mac Radigan
// All Rights Reserved

  package main

  import (
    "os"
    "bufio"
  )

  func display(out *os.File, ch chan []byte) {
    for char := range ch {
      out.Write(char)
    }
  }

  func main() {
    const BUF_SIZE int = 1024
    in  := bufio.NewReader(os.Stdin)
    out := os.Stdout
    buf := make([]byte, BUF_SIZE)
    ch  := make(chan []byte)
    go display(out, ch)
    for {
      _, err := in.Read(buf)
      if err != nil {
        os.Exit(0) /* EOF encountered */
      //panic(err)
      }
      ch <- buf
    }
  }

// *EOF*
