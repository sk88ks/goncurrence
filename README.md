# goncurrency

[![Build Status](https://travis-ci.org/sk88ks/goncurrency.svg?branch=master)](https://travis-ci.org/sk88ks/goncurrency)
 
[![Coverage Status](https://coveralls.io/repos/sk88ks/goncurrency/badge.svg?branch=master&service=github)](https://coveralls.io/github/sk88ks/goncurrency?branch=master)

goncurrency enables you to implement concurrnt processes more simple.

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

	workerNum := runtime.NumCPU()
	w := goncurrency.New(workerNum)

	// DefaultProcess implementing ProcessHandler interface
	processes := []goncurrency.DefaultProcess{
		goncurrency.DefaultProcess{
			Func: func() (interface{}, error) {
				return "result 0", nil
			},
		},
		goncurrency.DefaultProcess{
			Func: func() (interface{}, error) {
				return "result 1", nil
			},
		},
		goncurrency.DefaultProcess{
			Func: func() (interface{}, error) {
				return "result 2", nil
			},
		},
	}

	for i := range processes {
		w.Add(&processes[i])
	}

	err := w.Run()
	if err != nil {
		panic(err)
	}

	for i := range processes {
		result := processes[i].Result.(string)
		fmt.Printf("Result %d: %s\n", i, result)
	}     

    /* Results
      Result 0: result 0
      Result 1: result 1
      Result 2: result 2
    */

```
