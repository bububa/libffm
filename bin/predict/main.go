package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bububa/libffm/core"
	//"github.com/bububa/libffm/models"
	"github.com/bububa/libffm/tool"
)

var (
	dataPath  string
	modelPath string
)

func init() {
	flag.StringVar(&dataPath, "data", "", "")
	flag.StringVar(&modelPath, "model", "", "")
	flag.Parse()
}

func main() {
	model, err := tool.LoadModel(modelPath)
	if err != nil {
		log.Fatalln(err)
	}
	instances, err := tool.LoadData(dataPath)
	if err != nil {
		log.Fatalln(err)
	}
	for _, nodes := range instances {
		predictValue := core.Predict(model, nodes)
		fmt.Println(predictValue)
	}
}
