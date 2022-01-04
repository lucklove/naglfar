package main

import (
	"fmt"

	"github.com/lucklove/naglfar/cmd"
)

func main() {
	if err := cmd.Command().Execute(); err != nil {
		fmt.Println(err)
	}
}
