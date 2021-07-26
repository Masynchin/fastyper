package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var (
	answerLength   int
	timeoutSeconds int
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Add cmd options to set custom rules
func init() {
	flag.IntVar(&answerLength, "l", 4, "generated words length")
	flag.IntVar(&timeoutSeconds, "t", 10, "seconds per attempt")
}

// Set random seed and run game.
// After game show how many was guessed
func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	if answerLength < 1 || timeoutSeconds < 1 {
		log.Fatal(errors.New("options must be greater than or equal to 1"))
	}

	var guessed int
	guessed = PlayGame(guessed)
	color.Green("\nGuessed: %v", guessed)
}

// One round of game: correct answer - go next, incorrect - finish.
// Count guessed answers via recursion and return
func PlayGame(guessed int) int {
	userInput, answer := PrepareGame()
	fmt.Println(answer)
	go GetInput(userInput)

	userAnswer := HandleUserAnswer(userInput, answer)
	if userAnswer == nil {
		color.Yellow("\ntimed out!")
		return guessed
	} else if *userAnswer != answer {
		color.Red("incorrect!")
		return guessed
	} else {
		color.Green("correct!\n\n")
		guessed++
		return PlayGame(guessed)
	}
}

// Create variables used by game
func PrepareGame() (chan string, string) {
	userInput := make(chan string)
	answer := GenRandomString()
	return userInput, answer
}

// Set user input placeholder and put user answer in answer chan
func GetInput(c chan string) {
	fmt.Print("Type this word: ")
	var word string
	fmt.Scanln(&word)
	c <- word
}

// Generate random string for game
func GenRandomString() string {
	result := make([]byte, answerLength)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(result)
}

// Handle user input and return if not time out else return nil
func HandleUserAnswer(userInput <-chan string, answer string) *string {
	select {
	case userAnswer := <-userInput:
		return &userAnswer
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		return nil
	}
}
