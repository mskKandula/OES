package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"
)

var count int = 0
var arr = []string{"Mobile", "computer", "flower"}

func counter() js.Func {
	countFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		var answers map[string][]string
		inputString := args[0].String()
		json.Unmarshal([]byte(inputString), &answers)
		for _, answer := range answers["ans"] {
			for _, an := range strings.Fields(answer) {
				if contain(an, arr) {
					count++
				}
			}
		}
		return count
	})
	return countFunc
}

func contain(s string, wordDict []string) bool {
	for _, v := range wordDict {
		if s == v {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("countWordsInAns", counter())
	<-make(chan bool)
}
