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
// After the game show how many was guessed
func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	if answerLength < 1 || timeoutSeconds < 1 {
		log.Fatal(errors.New("options must be greater than or equal to 1"))
	}

	guessed := PlayGame()
	color.Green("\nGuessed: %v", guessed)
}

// Main game process - play rounds until user's answer
// is one of incorrect or timeout
func PlayGame() (guessed int) {
	for {
		isCorrect, timeout := PlayRound()
		switch {
		case timeout:
			color.Yellow("\ntimed out!")
			return
		case !isCorrect:
			color.Red("incorrect!")
			return
		default:
			guessed++
			color.Green("correct!\n\n")
		}
	}
}

// One round of game - return correctness of user answer flag
// and timeout flag if user didn't respond in time
func PlayRound() (isCorrect bool, timeout bool) {
	userInput, answer := PrepareGame()
	fmt.Println(answer)
	go GetInput(userInput)

	userAnswer, timeout := HandleUserAnswer(userInput)
	isCorrect = (userAnswer == answer)
	return
}

// Create variables used by game round
func PrepareGame() (chan string, string) {
	userInput := make(chan string)
	answer := GenRandomString()
	return userInput, answer
}

// Set user input placeholder and put user answer in answer chan
func GetInput(c chan<- string) {
	fmt.Print("Type this word: ")
	var word string
	fmt.Scanln(&word)
	c <- word
}

// Generate random string for game round
func GenRandomString() string {
	result := make([]byte, answerLength)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(result)
}

// Handle user input and return answer with empty timeout flag
// if user had time to answer else empty string with true timeout flag
func HandleUserAnswer(userInput <-chan string) (string, bool) {
	select {
	case userAnswer := <-userInput:
		return userAnswer, false
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		return "", true
	}
}
