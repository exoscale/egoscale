package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate config file for this cli",
}

func configCmdRun(cmd *cobra.Command, args []string) {
	isPrint, err := configCmd.Flags().GetBool("print")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)

	if viper.ConfigFileUsed() != "" {
		println("Good day! exo is already configured with accounts:")
		for _, acc := range allAccount.Accounts {
			print("- ", acc.Name)
			if acc.Name == allAccount.DefaultAccount {
				print(" [current]")
			}
			println("")
		}
		resp, err := readInput(reader, "Do you wish to add another account?", "Yn")
		if err != nil {
			log.Fatal(err)
		}
		if validateResponse(resp) {
			getAccount()
		}
	} else {
		println("Hi happy Exoscalian, some configuration is required to use exo")
		resp, err := readInput(reader, "Do you have an exoscale account already?", "Yn")
		if err != nil {
			log.Fatal(err)
		}
		if !validateResponse(resp) {
			println("explain how to create one's account")
			return
		}
		println(`
We now need some very important informations, find them there.
https://portal.exoscale.com/account/profile/api
`)
		account, err := getAccount()
		if err != nil {
			log.Fatal(err)
		}
		generateConfigFile(isPrint, account)
	}
}

func getAccount() (*account, error) {
	reader := bufio.NewReader(os.Stdin)

	account := &account{}

	apiURL, err := readInput(reader, "Compute API Endpoint", "https://api.exoscale.ch/compute")
	if err != nil {
		return nil, err
	}
	account.Endpoint = apiURL

	apiKey, err := readInput(reader, "API Key", "")
	if err != nil {
		return nil, err
	}
	account.Key = apiKey

	secretKey, err := readInput(reader, "Secret Key", "")
	if err != nil {
		return nil, err
	}
	account.Secret = secretKey

	return account, nil
}

func generateConfigFile(isPrint bool, newAccount *account) (string, error) {
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	filepath := path.Join(configFolder, "exoscale.toml")

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		return "", fmt.Errorf("File %q already exist", filepath)
	}

	viper.SetConfigType("toml")

	viper.SetConfigFile(filepath)

	viper.Set("defaultAccount", newAccount.Name)

	//config := &config{DefaultAccount: newAccount.Name, Accounts: []account{*newAccount}}

	if err := viper.WriteConfig(); err != nil {
		return "", err
	}

	return filepath, nil
}

func readInput(reader *bufio.Reader, text, def string) (string, error) {
	if def == "" {
		fmt.Printf("[+] %s [%s]: ", text, "none")
	} else {
		fmt.Printf("[+] %s [%s]: ", text, def)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input != "" {
		return input, nil
	}
	return def, nil
}

func validateResponse(response string) bool {
	return (strings.ToLower(response) == "y" || strings.ToLower(response) == "yes")
}

func init() {

	configCmd.Run = configCmdRun
	configCmd.Flags().BoolP("print", "p", false, "Print configuration")
	RootCmd.AddCommand(configCmd)
}
