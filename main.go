package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type item struct {
	Domain   string `json:"domain"`
	Tag      string `json:"tag"`
	Name     string `json:"name"`
	Classify string `json:"classify"`
}

var (
	tPath   = kingpin.Flag("tPath", "txt file path").Default("edu.txt").ExistingFile()
	jPath   = kingpin.Flag("jPath", "json file path").Default("edu.json").ExistingFile()
	compile = kingpin.Flag("compile", "json to txt").Bool()
)

func main() {
	kingpin.Parse()
	if *compile {
		txt2json()
	}
	f, err := os.Open(*jPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var items []item
	json.Unmarshal(b, &items)
	for _, v := range items {
		fmt.Printf("Tag:%s	Name:%s	Domain:%s \n", v.Tag, v.Name, v.Domain)
	}
	fmt.Printf("本库共收录%d所高校\n", len(items))
}

func txt2json() {
	f, err := os.Open(*tPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var items []item
	var str = strings.Split(string(b), "\n")
	var tempMap = map[string]struct{}{}
	for _, s := range str {
		l := strings.Split(s, "\t")
		if len(l) > 2 {
			tempItem := item{
				Tag:      strings.TrimPrefix(l[0], "Tag:"),
				Name:     strings.TrimPrefix(l[1], "Name:"),
				Domain:   strings.TrimPrefix(l[2], "Domain:"),
				Classify: "university",
			}
			if _, ok := tempMap[tempItem.Domain]; !ok {
				tempMap[tempItem.Domain] = struct{}{}
			} else {
				fmt.Printf("Repeat Item %s\n", s)
			}
			items = append(items, tempItem)
		} else {
			fmt.Printf("Invalid Format %s\n", s)
		}
	}
	sj, _ := json.Marshal(&items)
	ioutil.WriteFile(*jPath, []byte(sj), 0644)
	fmt.Printf("本库共收录%d所高校\n", len(items))
	os.Exit(0)
}
