package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// Initialising flags
	csvFilenamePtr := flag.String("csv", "problems.csv", "A CSV file containing questions and answer in `question,answer` format.")
	timeDurationPtr := flag.Int("limit", 30, "Time limit in which users must answer all questions.")
	mustShufflePtr := flag.Bool("shuffle", false, "Specifies whether the questions should be presented in a pseudo-random order")
	flag.Parse()

	// Parsing flag arguments
	csvFile, err := os.Open(*csvFilenamePtr)
	if err != nil {
		exitWithErrMsg(fmt.Sprintf("CSV file not found: %s\n", *csvFilenamePtr))
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		exitWithErrMsg("Error in parsing CSV file.\n")
	}

	problemList := parseCsvLines(lines)
	if *mustShufflePtr {
		rand.Shuffle(len(problemList), func(i, j int) {
			problemList[i], problemList[j] = problemList[j], problemList[i]
		})
	}

	// Main question-answer problem loop
	correctNum := 0
	timer := time.NewTimer(time.Duration(*timeDurationPtr) * time.Second)
	for ind, problem := range problemList {
		fmt.Printf("Problem #%d: %s\n", ind+1, problem.question)

		// Goroutine waits for user response on separate thread without blocking main loop
		responseCh := make(chan string)
		go func() {
			var response string
			fmt.Scanf("%s\n", &response)
			responseCh <- response
		}()

		// Select statemet waits for either timer or response channel and proceeds accordingly
		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			fmt.Printf("\nYou got %d out %d questions correct.\n", correctNum, len(problemList))
			return
		case response := <-responseCh:
			if response == problem.answer {
				correctNum++
				fmt.Println("Correct!")
			} else {
				fmt.Println("Wrong!")
			}
		}
	}
	fmt.Printf("\nYou got %d out %d questions correct.\n", correctNum, len(problemList))
}

type problem struct {
	question string
	answer   string
}

// Parses 2D slice of string into a 1D slice of problem
// lines must only have two columns: question and answer
func parseCsvLines(lines [][]string) []problem {
	problemList := make([]problem, len(lines))
	for ind, qaPair := range lines {
		problemList[ind] = problem{
			question: qaPair[0],
			answer:   strings.TrimSpace(qaPair[1]),
		}
	}
	return problemList
}

// Prints error message and exits with exit code 1
func exitWithErrMsg(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
