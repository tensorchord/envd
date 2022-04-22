package main

import "github.com/tensorchord/MIDI/pkg/progress"

func main() {
	l := progress.Current(false)
	l.WithPrefix("😀 [MIDI]").Printf("?")
}
