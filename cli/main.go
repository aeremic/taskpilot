package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type processDefinition struct {
	Name      string   `json:"name"`      // "name": "my-api",
	Cmd       string   `json:"cmd"`       // "cmd": "./my-api",
	Args      []string `json:"args"`      // "args": ["--port", "8080"],
	Cwd       string   `json:"cwd"`       // "cwd": "/home/user/projects/my-api",
	Instances int      `json:"instances"` // "instances": 2,
}

func get(path string) (*processDefinition, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := file
	decoder := json.NewDecoder(reader)

	var p processDefinition
	decoderErr := decoder.Decode(&p)
	if decoderErr != nil {
		return nil, decoderErr
	}

	return &p, nil
}

func parseStartCommand(parsedInput []string) {
	if len(parsedInput) < 3 {
		log.Printf("Invalid start command\n")
	}

	path := parsedInput[2]
	pd, err := get(path)
	if err != nil {
		log.Print(err)
	}

	// print(pd.Name)
}

func main() {
	// while
	// read io
	// parse with err check
	// switch with all of the commands

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		parsedInput := strings.Fields(strings.TrimSpace(input))

		if len(parsedInput) < 2 || strings.ToLower(parsedInput[0]) != "taskpilot" {
			log.Printf("Unknown command\n")
			continue
		}

		command := strings.ToLower(parsedInput[1])
		switch command {
		case "start":
			parseStartCommand(parsedInput)
			break
		case "stop":
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
		case "restart":
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
		case "list":
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
		case "logs":
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
		default:
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
