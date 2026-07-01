package main

import (
	"fmt"
	"os"
	"rustdesk-api-server-pro/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
