package collision

import "fmt"

func assert(a bool) {
	if !a {
		fmt.Println("assert failed") //log("")
	}
}
