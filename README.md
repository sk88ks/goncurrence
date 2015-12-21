# goncurrence

[![Build Status](https://travis-ci.org/sk88ks/goncurrence.svg?branch=master)](https://travis-ci.org/sk88ks/goncurrence)
 
[![Coverage Status](https://coveralls.io/repos/sk88ks/goncurrence/badge.svg?branch=master&service=github)](https://coveralls.io/github/sk88ks/goncurrence?branch=master)

Goncurrence enables you to implement concurrnt processes more simple.

Current API Documents:

Installation
----

```
go get github.com/sk88ks/goncurrence
```
 
Quick start
----

To create a new client and concurrntly execute processes

```go
import(
  "fmt"
  "runtime"
  "github.com/sk88ks/goncurrence"
)

	workerNum := runtime.NumCPU()
	w := goncurrence.New(workerNum)

	// DefaultProcess implementing ProcessHandler interface
	processes := []goncurrence.DefaultProcess{
		goncurrence.DefaultProcess{
			Func: func() (interface{}, error) {
				return "result 0", nil
			},
		},
		goncurrence.DefaultProcess{
			Func: func() (interface{}, error) {
				return "result 1", nil
			},
		},
		goncurrence.DefaultProcess{
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
