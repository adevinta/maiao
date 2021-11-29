package main

import (
	"fmt"
	"os"

	"github.com/adevinta/maiao/pkg/cmd"
)

func main() {
	if err := cmd.NewCommand().Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
