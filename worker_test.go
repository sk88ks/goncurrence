package goncurrence

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

func TestIsUnordered(t *testing.T) {
	Convey("Given worker manager", t, func() {
		w := New(2)

		Convey("When isUnorderd is called", func() {
			w.IsUnordered()

			Convey("Then flag should be set", func() {
				So(w.isUnorderd, ShouldBeTrue)

			})
		})
	})
}

func TestRun(t *testing.T) {
	Convey("Given default processes", t, func() {
		f := func(name string, num int) func() (interface{}, error) {
			msg := name + strconv.Itoa(num)
			return func() (interface{}, error) {
				if num == 99 {
					return nil, errors.New("Invalid param")
				}
				return msg, nil
			}
		}

		Convey("When adding and run with all valid processes", func() {
			w := New(2)
			ps := make([]DefaultProcess, 3)
			for i := 0; i < 3; i++ {
				p := DefaultProcess{
					Func: f("Message", i),
				}
				ps[i] = p
				w.Add(&ps[i])
			}
			err := w.Run()

			Convey("Then results should be set", func() {
				So(err, ShouldBeNil)
				for i := range ps {
					msg := ps[i].Result.(string)
					So(msg, ShouldEqual, "Message"+strconv.Itoa(i))
				}
			})
		})

		Convey("When adding and run with invalid process", func() {
			w := New(2)
			ps := make([]DefaultProcess, 3)
			for i := 0; i < 3; i++ {
				num := i
				if i == 2 {
					num = 99
				}
				p := DefaultProcess{
					Func: f("Message", num),
				}
				ps[i] = p
				w.Add(&ps[i])
			}
			err := w.Run()

			Convey("Then results should be set", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Invalid param")
			})
		})

		Convey("When adding and run with invalid process as unordered", func() {
			w := New(2)
			w.IsUnordered()
			ps := make([]DefaultProcess, 3)
			for i := 0; i < 3; i++ {
				num := i
				if i == 2 {
					num = 99
				}
				p := DefaultProcess{
					Func: f("Message", num),
				}
				ps[i] = p
				w.Add(&ps[i])
			}
			err := w.Run()

			Convey("Then results should be set", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Invalid param")
				for i := 0; i < 2; i++ {
					msg := ps[i].Result.(string)
					So(msg, ShouldEqual, "Message"+strconv.Itoa(i))
				}
				So(ps[2].Result, ShouldBeNil)
			})
		})
	})
}
