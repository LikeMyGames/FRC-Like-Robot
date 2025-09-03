package main

import (
	"encoding/json"
	"fmt"
	"os"
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
			Action:           func() { NewProject(args[1]) },
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
					fmt.Printf("%s: %s\n", v, cmd.ShortDescription)
				} else {
					fmt.Printf("%s: %s\n", v, cmd.LongDescription)
				}
			}
		} else {
			for i, v := range cmds {
				fmt.Printf("%s: %s\n", i, v.ShortDescription)
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

	os.Mkdir(name, os.ModeDir)

	file, _ := os.Create("./temp/settings.json")
	data, _ := json.MarshalIndent(settings, "", "\t")
	file.Write(data)
}
