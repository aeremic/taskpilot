package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func parseStartDeamonCommand(input string) (int, error) {
	process := exec.Command(input)
	err := process.Start()
	if err != nil {
		return 0, err
	}

	return process.Process.Pid, nil
}

func parseStopDeamonCommand(input string) error {
	pid, err := strconv.Atoi(input)
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	defer process.Release()

	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		return err
	}

	return nil
}

// func parseStartCommand(parsedInput []string) {
// 	if len(parsedInput) < 3 {
// 		log.Printf("Invalid start command\n")
// 	}

// 	path := parsedInput[2]

// 	// send command to deamon
// }

func main() {
	// while
	// read io
	// parse with err check
	// switch with all of the commands
	// sends commands to deamon

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
		case "start-deamon":
			if len(parsedInput) < 3 {
				log.Printf("Invalid start-deamon command\n")
				break
			}

			pid, err := parseStartDeamonCommand(parsedInput[2])
			if err != nil {
				log.Print(err)
			}

			log.Printf("%s", "Started deamon with "+strconv.Itoa(pid)+" pid.\n")
			break
		case "stop-deamon":
			err := parseStopDeamonCommand(parsedInput[2])
			if err != nil {
				log.Print(err)
			}
		case "start":
			log.Printf("%s", "Command "+command+" unsupported.\n")
			break
			// parseStartCommand(parsedInput)
			// break
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
