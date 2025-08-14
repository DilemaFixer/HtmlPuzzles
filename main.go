package main

import (
	"fmt"
	"time"

	render "github.com/DilemaFixer/HtmlPuzzles/reflection"
)

type One struct {
	Dummy int
	T     Two
}

type Two struct {
	Dummy string
	T     Three
}

type Three struct {
	Dummy bool
	F     Four
}

type Four struct {
	Dummy float64
	Four  int
}

func main() {
	one := One{
		T: Two{
			T: Three{
				F: Four{
					Four: 42,
				},
			},
		},
	}
	path := []string{"T", "T", "F", "Four"}
	path2 := []string{"T", "T", "F", "Dummy"}

	start := time.Now()
	offset1, ptrs1, _, err1 := render.FindOffsetForField(one, path)
	duration1 := time.Since(start)

	start = time.Now()
	offset2, ptrs2, _, err2 := render.FindOffsetForField(one, path2)
	duration2 := time.Since(start)

	fmt.Printf("%v (%d, %v, %v)\n", duration1, offset1, ptrs1, err1)
	fmt.Printf("%v (%d, %v, %v)\n", duration2, offset2, ptrs2, err2)
	fmt.Printf("Speed bust: %.2fx\n", float64(duration1)/float64(duration2))
}
