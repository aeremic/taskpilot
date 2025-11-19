package main

import (
	"encoding/json"
	"log"
	"os"
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

func main() {
	for {
		log.Print("Deamon started")
	}
}
