package main

import (
	dt "DecisionMaker/decisiontree"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPrintTree(t *testing.T) {
	file, _ := ioutil.ReadFile("treedb.json")

	data := dt.TreeNodes{}
	err := json.Unmarshal([]byte(file), &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)

}

func TestAskQuestion1(t *testing.T) {
	m := make(map[string]interface{}, 4)
	m["feeling sick ?"] = "yes"
	m["Are you in a strong pain?"] = "yes"
	m["do you have any of the following conditions? if yes select one of them"] = "no, i dont have any of the above conditions"
	m["How old are you?"] = "23"
	ResolveTree(m)
}

func TestAskQuestion2(t *testing.T) {
	m := make(map[string]interface{}, 4)
	m["feeling sick ?"] = "yes"
	m["Are you in a strong pain?"] = "no"

	ResolveTree(m)
}
