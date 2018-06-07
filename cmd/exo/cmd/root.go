package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	ttime "time"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/client"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var region string
var configFolder string
var configFilePath string
var cfgFilePath string

var cs *egoscale.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "exo",
	Short: "A simple CLI to use CloudStack using egoscale lib",
	//Long:  `A simple CLI to use CloudStack using egoscale lib`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", "", "Specify an alternate config file [env CLOUDSTACK_CONFIG]")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "cloudstack", "config ini file section name [env CLOUDSTACK_REGION]")

	cobra.OnInitialize(initConfig, buildClient)

}

func buildClient() {
	if cs != nil {
		return
	}

	var err error
	cs, err = client.BuildClient(configFilePath, region)
	if err != nil {
		log.Fatal(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

	filePrepender := func(filename string) string {
		now := ttime.Now().Format(ttime.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		url := "/commands/" + strings.ToLower(base) + "/"
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "/commands/" + strings.ToLower(base) + "/"
	}

	doc.GenMarkdownTreeCustom(rootCmd, "./docs", filePrepender, linkHandler)

	envs := map[string]string{
		"CLOUDSTACK_CONFIG": "config",
		"CLOUDSTACK_REGION": "region",
	}

	for env, flag := range envs {
		flag := rootCmd.Flags().Lookup(flag)
		if value := os.Getenv(env); value != "" {
			flag.Value.Set(value)
		}
	}

	envEndpoint := os.Getenv("CLOUDSTACK_ENDPOINT")
	envKey := os.Getenv("CLOUDSTACK_KEY")
	envSecret := os.Getenv("CLOUDSTACK_SECRET")

	if envEndpoint != "" && envKey != "" && envSecret != "" {
		cs = egoscale.NewClient(envEndpoint, envKey, envSecret)
		return
	}

	if cfgFilePath != "" {
		configFilePath = cfgFilePath
		return
	}

	usr, _ := user.Current()
	configFolder = path.Join(usr.HomeDir, ".exoscale")

	localConfig, _ := filepath.Abs("cloudstack.ini")
	inis := []string{
		localConfig,
		filepath.Join(usr.HomeDir, ".cloudstack.ini"),
		filepath.Join(configFolder, "cloudstack.ini"),
	}

	for _, i := range inis {
		if _, err := os.Stat(i); err != nil {
			continue
		}
		configFilePath = i
		break
	}

	if configFilePath == "" {
		path, err := generateConfigFile(false)
		if err != nil {
			log.Fatal(err)
		}
		configFilePath = path
	}

	if configFilePath == "" {
		log.Fatalf("Config file not found within: %s", strings.Join(inis, ", "))
	}
}
