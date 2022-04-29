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
	// Provides helper flags when using CLI eg "./quiz --help"
	csvFilename := flag.String("csv", "problems.csv", "CSV file in format of 'question,answer'")
	flag.Parse()

	// Create time limit flag
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")

	// Open the CSV file
	file, err := os.Open(*csvFilename) // (*)- use value not pointer
	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file %s\n", *csvFilename)) // %s - Print full file name
	}
	_ = file

	// Create csv reader and parse
	r := csv.NewReader(file)
	lines, err := r.ReadAll() // read all lines in CSV
	if err != nil {
		exit("Failed to parse the CSV file provided")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0
	// <-timer.C // program stops until message from channel

	// Print problems
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() { // Annonymous function
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C: // execute when msg from timer is rec'd
			fmt.Printf("Your score is %d out of %d. \n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}

		}
	}
	fmt.Printf("Your score is %d out of %d. \n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],                    // col 1
			a: strings.TrimSpace(line[1]), // col 2 - remove white space incase answer has it
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
