package main

import (
	"aliyun-sls-exporter/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
