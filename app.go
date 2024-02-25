package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type (
	userInput struct {
		value int
		err   error
	}

	gameStats struct {
		played int
		won    int
	}
)

type game struct {
	sigCh   chan os.Signal
	inputCh chan userInput
	stats   gameStats
}

func newGame() game {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	return game{
		sigCh:   c,
		inputCh: make(chan userInput),
		stats:   gameStats{},
	}
}

func (g *game) Start(limit uint) {

	for {
		randomInt := rand.Intn(int(limit))
		hint, err := getHint(randomInt)
		if err != nil {
			fmt.Println("error in getting the hint", err.Error())
			break
		}

		go g.getUserInput(hint)

		select {
		case <-g.sigCh:
			fmt.Printf("\nyou won %d out of %d rounds\n", g.stats.won, g.stats.played)
			fmt.Println("\nGood Game...!")
			os.Exit(0)

		case input := <-g.inputCh:
			if errors.Is(input.err, ErrInvalidInput) {
				fmt.Println("error in user input : ", input.err.Error())
				continue
			}

			if checkInput(input.value, randomInt) {
				g.stats.won++
			}
		}

		g.stats.played++
	}

}

func (g *game) getUserInput(hint string) {

	var value string
	fmt.Print("hint : ", hint, "\nguess the number :")
	fmt.Scan(&value)
	fmt.Println()

	input, err := strconv.Atoi(value)
	if err != nil {
		g.inputCh <- struct {
			value int
			err   error
		}{-1, err}

		return
	}

	g.inputCh <- struct {
		value int
		err   error
	}{input, nil}

}

func checkInput(userInput, randomInt int) bool {

	switch {
	case userInput < randomInt:
		fmt.Print("wrong answer, it was ", randomInt, "\n\n")
		return false

	case userInput > randomInt:
		fmt.Print("wrong answer, it was ", randomInt, "\n\n")
		return false

	default:
		fmt.Print("YES!! IT IS\n\n")
		return true
	}

}

func getHint(value int) (string, error) {

	url := "http://numbersapi.com/" + strconv.Itoa(value)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	source := string(bytes)

	pattern := regexp.MustCompile(`^\d+`)
	processedString := pattern.ReplaceAllString(source, "it")

	return processedString, nil

}
