package grouped_functions

import "time"

type (
	Group []func()
)

func NewParrallelGroup(functions func()) {

}

func (g *Group) RunParallel() {
	finished := make([]bool, len(*g))
	for i, f := range *g {
		go func() {
			finished[i] = false
			f()
			finished[i] = true
		}()
	}
	checkCompletion(&finished)
}

func (g *Group) RunSeries() {
	for _, f := range *g {
		f()
	}
}

func checkCompletion(finished *[]bool) bool {
	for _, v := range *finished {
		if !v {
			time.Sleep(time.Millisecond * 10)
			return checkCompletion(finished)
		}
	}
	return true
}
