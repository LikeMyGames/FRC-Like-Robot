package pathplanner

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadAllPathPlannerPaths() {
	entrys, err := os.ReadDir("./deploy/pathplanner/paths")
	if err != nil {
		panic(err)
	}

	for _, v := range entrys {
		if !v.IsDir() {
			paths[v.Name()] = loadPath(v.Name())
			fmt.Printf("%s\n\n", paths[v.Name()])
		}
	}
}

func loadAllPathPlannerAutos() {
	entrys, err := os.ReadDir("./deploy/pathplanner/paths")
	if err != nil {
		panic(err)
	}

	for _, v := range entrys {
		if !v.IsDir() {
			autos[v.Name()] = loadAuto(v.Name())
			fmt.Printf("%s\n\n", autos[v.Name()])
		}
	}
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
