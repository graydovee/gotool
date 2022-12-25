package main

import (
	"fmt"
	"github.com/grydovee/gotool/cmd/app"
	"os"
)

func main() {
	if err := app.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
