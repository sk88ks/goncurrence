package goncurrency

import (
	"errors"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("Given worker number", t, func() {
		num := 2

		Convey("When create new worker manager", func() {
			w := New(num)

			Convey("Then worker manager should be created", func() {
				So(w, ShouldNotBeNil)
				So(w.workerNum, ShouldEqual, num)

			})
		})
	})
}

func TestWorkerManager_Release(t *testing.T) {
	Convey("Given worker manager", t, func() {
		num := 2
		w := New(num)

		Convey("When release workers", func() {
			w.Release()

			Convey("Then worker open status should be false", func() {
				So(w.open, ShouldBeFalse)

			})
		})
	})
}

func TestWorkerManager_Add(t *testing.T) {
	Convey("Given worker manager and process funcs", t, func() {
		w := New(1)
		fs := []ProcessFunc{
			func() (interface{}, error) {
				return nil, nil
			},
			func() (interface{}, error) {
				return nil, nil
			},
		}

		Convey("When add funcs", func() {
			w.Add(fs...)

			Convey("Then count should be increased", func() {
				So(w.addCount, ShouldEqual, len(fs))

			})
		})
	})
}

func TestWorkerManager_Iter(t *testing.T) {
	Convey("Given worker manager and added process funcs", t, func() {
		w := New(1)
		fs := []ProcessFunc{
			func() (interface{}, error) {
				return nil, nil
			},
			func() (interface{}, error) {
				return nil, nil
			},
		}
		w.Add(fs...)

		Convey("When get iterator", func() {
			iter := w.Iter()

			Convey("Then count should be increased", func() {
				So(iter, ShouldNotBeNil)
				So(iter.wm, ShouldNotBeNil)

			})
		})
	})
}

func TestProcessIterator_Next(t *testing.T) {
	Convey("Given iterator having processes", t, func() {
		w := New(1)
		fs := []ProcessFunc{
			func() (interface{}, error) {
				return nil, nil
			},
			func() (interface{}, error) {
				return nil, nil
			},
		}
		w.Add(fs...)
		iter := w.Iter()

		Convey("When iterate next process result", func() {
			var count int
			for iter.Next() {
				count++
			}

			Convey("Then count should be increased", func() {
				So(count, ShouldEqual, len(fs))

			})
		})

		Convey("When iterate next process result after close", func() {
			var count int
			for iter.Next() {
				count++
				if count == 1 {
					w.Release()
				}
			}

			Convey("Then count should be increased", func() {
				So(count, ShouldEqual, 1)
				So(iter.wm.open, ShouldBeFalse)

			})
		})
	})
}

func TestProcessIterator_Result(t *testing.T) {
	Convey("Given iterator having processes", t, func() {
		w := New(1)
		type res struct {
			Name string
			Age  int
		}
		fs := []ProcessFunc{
			func() (interface{}, error) {
				return res{
					"test001",
					1,
				}, nil
			},
			func() (interface{}, error) {
				return &res{
					"test002",
					2,
				}, nil
			},
			func() (interface{}, error) {
				return nil, errors.New("Test error")
			},
		}
		w.Add(fs...)
		iter := w.Iter()

		Convey("When get result by iterated value", func() {
			var count int
			var errs []error
			var resulsts []res
			var r res
			for iter.Next() {
				count++
				err := iter.Result(&r)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				resulsts = append(resulsts, r)
			}

			Convey("Then result and errors should be set correctly", func() {
				So(count, ShouldEqual, len(fs))
				So(len(errs), ShouldEqual, 1)
				So(errs[0].Error(), ShouldEqual, "Test error")
				So(len(resulsts), ShouldEqual, 2)
				for i := range resulsts {
					So(resulsts[i].Name, ShouldEqual, "test00"+strconv.Itoa(resulsts[i].Age))
				}

			})
		})

		Convey("When get result with un pointer dst", func() {
			var count int
			var resulsts []res
			var r res
			for iter.Next() {
				count++
				iter.Result(r)
				if r.Name != "" {
					resulsts = append(resulsts, r)
				}
			}

			Convey("Then count should be increased", func() {
				So(count, ShouldEqual, len(fs))
				So(len(resulsts), ShouldEqual, 0)

			})
		})
	})

}
