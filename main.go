package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilenamePtr := flag.String("csv", "problems.csv", "A CSV file containing questions and answer in `question,answer` format.")
	timeDurationPtr := flag.Int("limit", 30, "Time limit in which users must answer all questions.")
	flag.Parse()

	csvFile, err := os.Open(*csvFilenamePtr)
	if err != nil {
		exitWithErrMsg(fmt.Sprintf("CSV file not found: %s\n", *csvFilenamePtr))
	}

	lines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		exitWithErrMsg("Error in parsing CSV file.\n")
	}

	problemList := parseCsvLines(lines)
	correctNum := 0
	timer := time.NewTimer(time.Duration(*timeDurationPtr) * time.Second)
	for ind, problem := range problemList {
		fmt.Printf("Problem #%d: %s\n", ind+1, problem.question)

		responseCh := make(chan string)
		go func() {
			var response string
			fmt.Scanf("%s\n", &response)
			responseCh <- response
		}()

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

func exitWithErrMsg(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
