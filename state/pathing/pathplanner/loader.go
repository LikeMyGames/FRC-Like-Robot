package pathplanner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func loadAllPathPlannerPaths() {
	entrys, err := os.ReadDir("./deploy/pathplanner/paths")
	if err != nil {
		panic(err)
	}

	sum := 0

	for _, v := range entrys {
		if !v.IsDir() && strings.Contains(v.Name(), ".path") {
			sum++
			name, _ := strings.CutSuffix(v.Name(), ".path")
			paths[name] = loadPath(name)
			fmt.Printf("%s\n\n", paths[name])
		}
	}

	fmt.Printf("Found %v Paths\n", sum)
}

func loadAllPathPlannerAutos() {
	entrys, err := os.ReadDir("./deploy/pathplanner/paths")
	if err != nil {
		panic(err)
	}

	sum := 0

	for _, v := range entrys {
		if !v.IsDir() && strings.Contains(v.Name(), ".auto") {
			sum++
			name, _ := strings.CutSuffix(v.Name(), ".auto")
			autos[name] = loadAuto(name)
			fmt.Printf("%s\n\n", autos[name])
		}
	}

	fmt.Printf("Found %v Autos\n", sum)
}

func loadPath(name string) *Path {
	bytes, err := os.ReadFile(fmt.Sprintf("./deploy/pathplanner/paths/%s.path", name))
	if err != nil {
		panic(err)
	}

	path := new(Path)
	json.Unmarshal(bytes, path)
	path.Name = name

	return path
}

func loadAuto(name string) *Auto {
	bytes, err := os.ReadFile(fmt.Sprintf("./deploy/pathplanner/autos/%s.auto", name))
	if err != nil {
		panic(err)
	}

	auto := new(Auto)
	json.Unmarshal(bytes, auto)
	auto.Name = name

	return auto
}
