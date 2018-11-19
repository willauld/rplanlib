package rplanlib

func Assert(a bool, str string) {
	if !a {
		panic(str)
	}
}
