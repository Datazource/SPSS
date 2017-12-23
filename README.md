# SPSS

Secure password for Golang.

This package provides secure and concurrent access to store password in memory for Windows and Linux.

```go
import "github.com/themester/SPSS"

func main() {
  shdw := &Shadow{}

  // do not reset my password
  err := shdw.Read(0)
  if err != nil {
    panic(err)
  }
  println("My password is:", *shdw.Get())
}

```
