package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

type Question struct {
	Text   string
	Answer string
}

var counter int64

func main() {
	var timeLimit time.Duration
	var random bool
	fmt.Println("Hello, 'Enter' to start the quiz.")

	timeLimitPtr := flag.Int("timeLimit", 30, "time limit for quiz in seconds")
	randomPtr := flag.Bool("random", false, "randomize questions order")

	flag.Parse()
	timeLimit = time.Duration(*timeLimitPtr) * time.Second
	random = *randomPtr

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	startQuiz(scanner, timeLimit, random)
}

func startQuiz(scanner *bufio.Scanner, timeLimit time.Duration, random bool) {
	questionDb, err := initializeQuestionDb(random)
	if err != nil {
		fmt.Printf("Can't open questions file: %v\n", err)
		return
	}

	done := make(chan bool)
	timer := time.NewTimer(timeLimit)
	go func() {
		gameExcution(scanner, questionDb)
		done <- true
	}()
	select {
	case <-done:
	case <-timer.C:
		fmt.Println("Time's up!")
	}
	fmt.Println("Correct answer:", counter)
}

func initializeQuestionDb(random bool) ([]Question, error) {
	file, err := os.Open("problems.csv")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	questions := make([]Question, 0, len(records))
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		questions = append(questions, Question{Text: record[0], Answer: record[1]})
	}

	if random {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	return questions, nil
}

func gameExcution(scanner *bufio.Scanner, questionDb []Question) {
	for _, question := range questionDb {
		fmt.Println("Question:", question.Text)
		scanner.Scan()
		userAnswer := scanner.Text()

		if compareAnswers(userAnswer, question.Answer) {
			atomic.AddInt64(&counter, 1)

		}
	}
}

func compareAnswers(userAnswer, correctAnswer string) bool {
	userAnswer = strings.TrimSpace(userAnswer)
	userAnswer = strings.ReplaceAll(userAnswer, " ", "")
	return strings.EqualFold(userAnswer, correctAnswer) || strings.EqualFold(userAnswer, correctAnswer+"?") || strings.EqualFold(userAnswer, correctAnswer+".")
}
