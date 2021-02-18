package main

import (
	"fmt"
	"os"
	"errors"
)

func main() {
  argsSize := len(os.Args)
  if argsSize > 1 {
    fmt.Print("Hello ")
    i := 1
    for i <= argsSize - 1 {
      fmt.Print(os.Args[i])
      fmt.Print(" ")
      i = i + 1
    }
    fmt.Print(", Welcome to the Jungle \n")
  } else {
    err := errors.New("Error, Missing Name")
		if err != nil {
		  fmt.Println(err)
	  }
  }
}
