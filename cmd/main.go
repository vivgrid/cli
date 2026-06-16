package main

import (
	"fmt"
	"os"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/cobra"
	"github.com/vivgrid/cli/pkg"
)

func main() {
	tid, err := gonanoid.New(8)
	if err != nil {
		fmt.Println("generate target error:", err)
		os.Exit(1)
	}

	var configFile string
	if c, ok := os.LookupEnv("VIV_CONFIG_FILE"); ok {
		configFile = c
	} else if c, ok := os.LookupEnv("YC_CONFIG_FILE"); ok {
		// fall back to the legacy env var for backward compatibility
		configFile = c
	} else if _, err := os.Stat("./vivgrid.yml"); err == nil {
		configFile = "./vivgrid.yml"
	} else if _, err := os.Stat("./yc.yml"); err == nil {
		// fall back to the legacy config file for backward compatibility
		configFile = "./yc.yml"
	}

	rootCmd := &cobra.Command{
		Use:   "viv",
		Short: "Manage your globally deployed Serverless LLM Functions on vivgrid.com from the command line",
	}

	err = pkg.Execute(rootCmd, configFile, tid, "https://hosting.vivgrid.com")
	if err != nil {
		fmt.Println("cmd error:", err)
		os.Exit(1)
	}
}
