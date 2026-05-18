package main

import (
	"fmt"
	"os"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/cobra"
	"github.com/vivgrid/yc/pkg"
)

func main() {
	tid, err := gonanoid.New(8)
	if err != nil {
		fmt.Println("generate target error:", err)
		os.Exit(1)
	}

	var configFile string
	if c, ok := os.LookupEnv("YC_CONFIG_FILE"); ok {
		configFile = c
	} else if _, err := os.Stat("./yc.yml"); err == nil {
		configFile = "./yc.yml"
	}

	rootCmd := &cobra.Command{
		Use:   "yc",
		Short: "Manage your globally deployed Serverless LLM Functions on vivgrid.com from the command line",
	}

	err = pkg.Execute(rootCmd, configFile, tid, "http://127.0.0.1:9040")
	if err != nil {
		fmt.Println("cmd error:", err)
		os.Exit(1)
	}
}
