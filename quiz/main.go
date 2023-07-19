package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func parseCli() (string, int) {
	var csvFile string
	var timeout int
	flag.StringVar(&csvFile, "f", "problems.csv", "CSV file to parse")
	flag.IntVar(&timeout, "t", 30, "timeout to finish quiz in seconds")
	flag.Parse()
	return csvFile, timeout
}

func readCsvFile(csvFile string) [][]string {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to read input file: %s\n%s", csvFile, err))
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to parse CSV file: %s\n%s", csvFile, err))
	}
	if cols := len(records[0]); cols != 2 {
		log.Fatal(fmt.Sprintf("CSV file should have 2 fields not %d", cols))
	}
	return records
}

func runGame(over chan bool, valid *int, records [][]string) {
	*valid = 0
	reader := bufio.NewReader(os.Stdin)
	for _, fields := range records {
		fmt.Println(fields[0])
		buff, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading from console\n", err)
		}
		buff = buff[:len(buff)-1] // remove last trailing "\n"
		if fields[1] == buff {
			*valid++
		}
	}
	over <- true
}

func execute(csvFile string, timeout int) {
	records := readCsvFile(csvFile)
	valid := 0                 // valid answer counter
	over := make(chan bool, 1) // set to true when player has answered every questions
	fmt.Printf("You have to answer every question correctly in less than %d seconds to win this game\nPress enter when you are ready", timeout)
	fmt.Scanln()
	go runGame(over, &valid, records)
	select {
	case <-over:
		// player answered every question
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Printf("Game timed out after %d seconds\n", timeout)
	}
	fmt.Printf("Result: %d/%d\n", valid, len(records))
}

func main() {
	csvFile, timeout := parseCli()
	execute(csvFile, timeout)
}
