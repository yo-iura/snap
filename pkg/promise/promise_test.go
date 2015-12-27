package promise

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPromise(t *testing.T) {
	Convey("IsComplete()", t, func() {
		Convey("it should return the completion status", func() {
			p := NewPromise()
			So(p.IsComplete(), ShouldBeFalse)
			p.Complete([]error{})
			So(p.IsComplete(), ShouldBeTrue)
		})
	})
	Convey("Complete()", t, func() {
		Convey("it should unblock any waiting goroutines", func() {
			p := NewPromise()

			numWaiters := 3
			var wg sync.WaitGroup
			wg.Add(numWaiters)

			for i := 0; i < numWaiters; i++ {
				go func() {
					Convey("all waiting goroutines should see empty errors", t, func() {
						errors := p.Await()
						So(errors, ShouldBeEmpty)
						wg.Done()
					})
				}()
			}

			p.Complete([]error{})
			wg.Wait()
		})
	})
	Convey("AndThen()", t, func() {
		Convey("it should defer the supplied closure until after completion", func() {
			p := NewPromise()

			funcRan := false
			c := make(chan struct{})

			p.AndThen(func(errors []error) {
				funcRan = true
				close(c)
			})

			// The callback should not have been executed yet.
			So(funcRan, ShouldBeFalse)

			// Trigger callback execution by completing the queued job.
			p.Complete([]error{})

			// Wait for the deferred function to be executed.
			<-c
			So(funcRan, ShouldBeTrue)
		})
	})
}

func TestRendezVous(t *testing.T) {
	Convey("IsComplete()", t, func() {
		Convey("it should return the completion status", func() {
			r := NewRendezVous()
			So(r.IsComplete(), ShouldBeFalse)
			go r.A()
			r.B()
			So(r.IsComplete(), ShouldBeTrue)
		})
	})
	Convey("A() and B()", t, func() {
		Convey("it should synchronize concurrent processes", func() {
			r1, r2 := NewRendezVous(), NewRendezVous()
			evidence := false

			go func() {
				r1.A()
				evidence = true
				r2.A()
			}()

			r1.B()
			r2.B()
			So(evidence, ShouldBeTrue)
		})
	})
}
