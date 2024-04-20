package main

import (
	"fmt"
	"os"
	"strconv"
)

func writePodScoreToFile(podName string, score int, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("team: TeamName=%s, score=%v\n", podName, score))
	if err != nil {
		fmt.Printf("error writing to file: %v", err)
		return
	}
}

func main() {
	podName, exists := os.LookupEnv("POD_NAME")
	if !exists {
		fmt.Println("POD_NAME environment variable not set")
		os.Exit(1)
	}

	scoreStr, exists := os.LookupEnv("SCORE")
	if !exists {
		fmt.Println("SCORE environment variable not set")
		os.Exit(1)
	}

	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		fmt.Printf("error converting SCORE to integer: %v", err)
		os.Exit(1)
	}

	filename, exists := os.LookupEnv("FILENAME")
	if !exists {
		fmt.Println("FILENAME environment variable not set")
		os.Exit(1)
	}

	writePodScoreToFile(podName, score, filename)
}
