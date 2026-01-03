package main

import (
	"common"
	"database/sql"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type processDefinition struct {
	Pid       int
	Name      string   `json:"name"`      // "name": "my-api",
	Cmd       string   `json:"cmd"`       // "cmd": "./my-api",
	Args      []string `json:"args"`      // "args": ["--port", "8080"],
	Cwd       string   `json:"cwd"`       // "cwd": "/home/user/projects/my-api",
	Instances int      `json:"instances"` // "instances": 2,
}

func write(conn net.Conn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Fatal(err)
	}
}

func startProcess(pd processDefinition) (*exec.Cmd, error) {
	process := exec.Command(pd.Cmd, pd.Args...)
	err := process.Start()
	if err != nil {
		return nil, err
	}
	return process, nil
}

func getSocketConnection() (net.Listener, error) {
	socketPath := "/tmp/echo.sock"
	os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func getDbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./storage.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db, nil
}

func reRunProcesses(db *sql.DB) error {
	var processesFromDb []processDefinition
	err := db.QueryRow("SELECT * FROM PROCESS;").Scan(&processesFromDb)
	if err != nil {
		log.Fatal("Error on database read: ", err)
	}

	for _, pd := range processesFromDb {
		processFromSys, err := os.FindProcess(pd.Pid)
		if err == nil {
			_, err = startProcess(pd)
			if err != nil {
				return err
			}
			continue
		}

		err = processFromSys.Signal(syscall.Signal(0))
		if err == nil {
			_, err = startProcess(pd)
			if err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func handleConnection(conn net.Conn, db *sql.DB) {
	for {
		var err error

		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			log.Print(err)
		}

		msg := strings.Fields(strings.TrimSpace(string(buf[0:n])))
		if len(msg) < 1 {
			write(conn, "Invalid command")
		}

		command := msg[0]

		switch command {
		case "start":
			if len(msg) < 2 {
				write(conn, "Invalid start command")
			}

			pdPath := msg[1]

			pd, err := common.GetAndDecodeJsonFile[processDefinition](pdPath)
			process := exec.Command(pd.Cmd, pd.Args...)
			process, err = startProcess(*pd)
			if err != nil {
				write(conn, err.Error())
				break
			}

			q := "INSERT INTO process(pid, name, cmd, args, cwd, instances) values(?, ?, ?, ?, ?, ?)"
			_, err = db.Exec(q, process.Process.Pid, pd.Name, pd.Cmd, pd.Args, pd.Cwd, pd.Instances)
			if err != nil {
				write(conn, err.Error())
				break
			}

			write(conn, "Process "+pd.Name+" started.")
			break
		case "stop":
			// TODO: Find running processes in the table and get their pids. Then stop them with sigkill. Here it might happen that process is already killed/not existant. Return msg then
			// TODO: There can be multiple processes by name since multiple processes can be runned so it needs to stop them all
			if len(msg) < 2 {
				write(conn, "Invalid start command")
			}

			var processesFromDb []processDefinition
			err := db.QueryRow("SELECT * FROM PROCESS WHERE Name = ?;", msg[1]).Scan(&processesFromDb)
			if err != nil {
				log.Fatal("Error on database read: ", err)
			}

			for _, pd := range processesFromDb {
				process, err := os.FindProcess(pd.Pid)
				if err != nil {
					write(conn, err.Error())
					continue
				}

				defer process.Release()

				err = process.Signal(syscall.SIGKILL)
				if err != nil {
					write(conn, err.Error())
					continue
				}
			}

			write(conn, "Stopping "+msg[1]+" process is done.")
			break
		case "restart":
			// TODO: Find running processes in the table. Sigkill then run them again them if they exist.
			// TODO: Here it might happen that restarting process is already killed/not existant. Only start then but return msg
		case "list":
			// TODO: Query all running processes from the table. Consider joining db state with running state since db state might be out of sync. Look into syncing states
			break
		default:
			write(conn, "Unsupported "+command+" command")
		}
	}
}

func main() {
	l, err := getSocketConnection()
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	db, err := getDbConnection()
	if err != nil {
		log.Fatal("Database error: ", err)
	}

	err = reRunProcesses(db)
	if err != nil {
		log.Fatal("Unable to restart running process: ", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}

		go handleConnection(conn, db)
	}
}
