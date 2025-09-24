package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type (
	Command struct {
		Action           func()
		ShortDescription string
		LongDescription  string
	}

	Settings struct {
		Name         string
		Version      string
		RobotIP      string
		TeamNum      uint8
		Architecture string
	}
)

func main() {
	args := os.Args[1:]

	cmds := map[string]Command{
		"create": {
			Action: func() {
				if len(args) == 1 {
					fmt.Println("The create command takes an argument as the name of the project")
					os.Exit(0)
				}
				NewProject(args[1])
			},
			ShortDescription: "A command that take the next provided argument and creates a new gobot project with that name",
			LongDescription:  "",
		},
		"build": {
			Action:           func() {},
			ShortDescription: "A command that compiles the code into an exe and sends it to the robot to be executed",
			LongDescription:  "",
		},
		"compile": {
			Action:           func() {},
			ShortDescription: "A command that compiles the gobot project into a linux .exe file",
			LongDescription:  "",
		},
		"send": {
			Action:           func() {},
			ShortDescription: "A command that sends a compiled linux .exe file to the robot",
			LongDescription:  "",
		},
	}

	if args[0] == "help" {
		if len(args) > 1 {
			for _, v := range args[1:] {
				cmd := cmds[v]
				if cmd.LongDescription == "" {
					fmt.Printf("%s\t %s\n", v, cmd.ShortDescription)
				} else {
					fmt.Printf("%s\t %s\n", v, cmd.LongDescription)
				}
			}
		} else {
			for i, v := range cmds {
				fmt.Printf("%s\t %s\n", i, v.ShortDescription)
			}
		}
	} else {
		for i, v := range cmds {
			if i == args[0] {
				v.Action()
				return
			}
		}
	}
	fmt.Println("The command that you are trying to use does not exist, try using <help> to learn what commands area available")
}

func NewProject(name string) {
	settings := Settings{
		Name:         name,
		TeamNum:      1,
		Version:      "0.0.0",
		Architecture: "Stated",
	}

	// ./{project-name}/
	os.Mkdir(settings.Name, os.ModeDir)

	// ./project.json file
	file, _ := os.Create(fmt.Sprintf("./%s/project.json", settings.Name))
	data, _ := json.MarshalIndent(settings, "", "\t")
	file.Write(data)

	// main.go file
	file, _ = os.Create(fmt.Sprintf("./%s/main.go", settings.Name))
	resp, err := http.Get("https://raw.githubusercontent.com/LikeMyGames/FRC-Like-Robot/refs/heads/main/main.go_template.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprint("Error: Received non-OK status code:", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprint("Error reading response body:", err))
	}

	file.WriteString(string(body))
	file.Close()

	// constants.go
	os.Mkdir(fmt.Sprintf("./%s/constants", settings.Name), os.ModeDir)
	file, _ = os.Create(fmt.Sprintf("./%s/constants/constants.go", settings.Name))
	resp, err = http.Get("https://raw.githubusercontent.com/LikeMyGames/FRC-Like-Robot/refs/heads/main/constants.go_template.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprint("Error: Received non-OK status code:", resp.StatusCode))
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprint("Error reading response body:", err))
	}

	file.WriteString(string(body))
	file.Close()

	err = os.Chdir(fmt.Sprintf("./%s", settings.Name))
	if err != nil {
		panic(err)
	}

	// go.mod
	if err = exec.Command("go", "mod", "init", settings.Name).Run(); err != nil {
		panic(fmt.Sprint("Could not create go.mod file:", err.Error()))
	}

	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/robot").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/constants").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/conn").Run()

	cmd := exec.Command("go", "mod", "tidy")
	fmt.Println(cmd.Path)
	fmt.Println(cmd.Args)
	fmt.Println(os.Getwd())

	// adding dependencies to src/go.mod
	out, err := cmd.CombinedOutput()
	// if  err != nil {
	// 	// panic(fmt.Sprint("Coult not create dependency in go.mod file:", err))
	// }
	fmt.Println(string(out), err)
}
