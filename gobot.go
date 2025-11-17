package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type (
	Command struct {
		Action           func()
		ShortDescription string
		LongDescription  string
	}

	Settings struct {
		Name         string `json:"Name"`
		Version      string `json:"Version"`
		RobotIP      string `json:"RobotIP"`
		TeamNum      uint8  `json:"TeamNum"`
		Architecture string `json:"Architecture"`
		EntranceFile string `json:"EntranceFile"`
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
		"deploy": {
			Action: func() {
				CompileProject()
				TransferExeToRobot()
			},
			ShortDescription: "A command that compiles the code into an exe and sends it to the robot to be executed",
			LongDescription:  "",
		},
		"compile": {
			Action: func() {
				CompileProject()
			},
			ShortDescription: "A command that compiles the gobot project into a linux .exe file",
			LongDescription:  "",
		},
		"send": {
			Action: func() {
				TransferExeToRobot()
			},
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
		EntranceFile: "main.go",
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

	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state@0.0.0").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/robot@0.0.0").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/constants@0.0.0").Run()
	exec.Command("go", "get", "github.com/LikeMyGames/FRC-Like-Robot/state/conn@0.0.0").Run()

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

func CompileProject() {
	fileData, err := os.ReadFile("project.json")
	if err != nil {
		fmt.Println("Could not find files necessary for Gobot project compilation. Make sure you are in a valid Gobot project directory.")
		return
	}
	data := Settings{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		fmt.Println("Could not read contents of project.json file correctly. Make sure you are in a valid Gobot project directory.")
		return
	}

	curGOOS := runtime.GOOS
	curGOARCH := runtime.GOARCH
	err = exec.Command("go", "env", "-w", "GOOS=linux", "GOARCH=arm64").Run()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("go", "build", "-o", "build/bin.exe", data.EntranceFile)
	// fmt.Println(cmd.Args)
	// err = cmd.Run()
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}
	err = exec.Command("go", "env", "-w", fmt.Sprintf("GOOS=%s", curGOOS), fmt.Sprintf("GOARCH=%s", curGOARCH)).Run()
	if err != nil {
		panic(err)
	}
}

func TransferExeToRobot() {
	fileData, err := os.ReadFile("project.json")
	if err != nil {
		fmt.Println("Could not find files necessary for Gobot project compilation. Make sure you are in a valid Gobot project directory.")
		return
	}
	data := Settings{}
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		fmt.Println("Could not read contents of project.json file correctly. Make sure you are in a valid Gobot project directory.")
		return
	}

	wd, err := os.Getwd()
	buildPath := filepath.Join(wd, "build", "bin.exe")
	if err != nil {
		fmt.Println("Could not get working directory of command execution")
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:8080", data.RobotIP)) // Replace localhost with server IP
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	sourceFile, err := os.Open(buildPath) // Replace with your file
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer sourceFile.Close()

	bytesSent, err := io.Copy(conn, sourceFile)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Printf("Sent %d bytes from %s\n", bytesSent, buildPath)
}
