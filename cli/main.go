package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	common "common"
)

type CustomError struct {
	Message string
}

func (c CustomError) Error() string {
	return fmt.Sprintf("%s", c.Message)
}

type deamonConfig struct {
	Path string `json:"path"`
}

var conn net.Conn = nil

func dial() {
	if conn == nil {
		var err error

		conn, err = net.Dial("unix", "/tmp/echo.sock")
		if err != nil {
			log.Print(err)
		}
	}
}

func write(msg string) error {
	if conn != nil {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			return err
		}

		return nil
	}

	return CustomError{Message: "Unable to write. Connection not established."}
}

func read() (string, error) {
	if conn != nil {
		buf := make([]byte, 512)
		n, err := conn.Read(buf[:])
		if err != nil {
			return "", err
		}

		data := string(buf[0:n])
		return data, nil
	}

	return "", CustomError{Message: "Unable to read. Connection not established."}
}

func processStartDeamonCommand(input string) (int, error) {
	// TODO: prevent starting deamon multiple times by saving deamon state
	process := exec.Command(input)
	err := process.Start()
	if err != nil {
		return -1, err
	}

	return process.Process.Pid, nil
}

func processStopDeamonCommand(input string) error {
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

func processStartCommand(parsedInput []string) error {
	err := write(parsedInput[1] + " " + parsedInput[2])
	if err != nil {
		return err
	}

	data, err := read()
	if err != nil {
		return err
	}

	log.Print(data)

	return nil
}

func processStopCommand(parsedInput []string) error {
	err := write(parsedInput[1] + " " + parsedInput[2])
	if err != nil {
		return err
	}

	data, err := read()
	if err != nil {
		return err
	}

	log.Print(data)

	return nil
}

func main() {
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

			json, err := common.GetAndDecodeJsonFile[deamonConfig](parsedInput[2])
			if err != nil {
				log.Print(err)
				break
			}

			pid, err := processStartDeamonCommand(json.Path)
			if err != nil {
				log.Print(err)
				break
			}

			log.Printf("%s", "Started deamon with "+strconv.Itoa(pid)+" pid.\n")
			break
		case "stop-deamon":
			if len(parsedInput) < 3 {
				log.Printf("Invalid stop-deamon command\n")
				break
			}

			err := processStopDeamonCommand(parsedInput[2])
			if err != nil {
				log.Print(err)
				break
			}

			if conn != nil {
				conn.Close()
				conn = nil
			}

			log.Print("Stopped.")
			break
		case "start":
			if len(parsedInput) < 3 {
				log.Printf("Invalid start command\n")
			}

			dial()
			err := processStartCommand(parsedInput)
			if err != nil {
				log.Print(err)
				break
			}

			break
		case "stop":
			if len(parsedInput) < 3 {
				log.Printf("Invalid stop command\n")
			}

			dial()
			err := processStopCommand(parsedInput)
			if err != nil {
				log.Print(err)
				break
			}

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
