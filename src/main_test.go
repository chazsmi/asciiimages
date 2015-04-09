package main

import "testing"

func BenchmarkProcess(b *testing.B) {
	process("/Users/charlie/Downloads/220px-Tennis_Racket_and_Balls.jpg")
}