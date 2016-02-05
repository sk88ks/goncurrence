# goncurrency

[![Build Status](https://travis-ci.org/sk88ks/goncurrency.svg?branch=master)](https://travis-ci.org/sk88ks/goncurrency)
[![Coverage Status](https://coveralls.io/repos/sk88ks/goncurrency/badge.svg?branch=master&service=github)](https://coveralls.io/github/sk88ks/goncurrency?branch=master)

goncurrency enables you to implement concurrnt processes more easily.

Current API Documents:

Installation
----

```
go get github.com/sk88ks/goncurrency
```
 
Quick start
----

To create a new client and concurrntly execute processes

```go
import(
  "fmt"
  "runtime"
  "github.com/sk88ks/goncurrency"
)

	w := goncurrency.New(1)

	// DefaultProcess implementing ProcessHandler interface
	processes := []goncurrency.ProcessFunc{
		func() (interface{}, error) {
			return "result 0", nil
		},
		func(interface{}, error) {
			return "result 1", nil
		},
		func() (interface{}, error) {
			return "result 2", nil
		},
	}

	w.Add(processes...)

    iter := w.Iter()

    var count int
    var result string
    for iter.Next() {
        err := iter.Result(&result)
        if err != nil {
            panic(err)
        }
		fmt.Printf("Result %d: %s\n", count, result)
        count++
    }

    /* Results unordered
      Result: result 0
      Result: result 1
      Result: result 2
    */

```
