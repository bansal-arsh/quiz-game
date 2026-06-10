package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFilenamePtr := flag.String("csv", "problems.csv", "A CSV file containing questions and answer in `question,answer` format.")
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
	for ind, problem := range problemList {
		fmt.Printf("Problem #%d: %s\n", ind+1, problem.question)
		var response string
		fmt.Scanf("%s\n", &response)
		if response == problem.answer {
			correctNum++
			fmt.Println("Correct!")
		} else {
			fmt.Println("Wrong!")
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
