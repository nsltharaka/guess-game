package main

import "flag"

func main() {

	limit := flag.Uint("limit", 10, "a number more than 0 which specifies the upper limit of guessing number")
	flag.Parse()

	game := newGame()
	game.Start(*limit)

}
