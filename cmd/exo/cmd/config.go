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
	// isPrint, err := configCmd.Flags().GetBool("print")
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
			account, err := getAccount()
			if err != nil {
				log.Fatal(err)
			}
			isDefault := false
			resp, err := readInput(reader, "Make "+account.Name+" your current profile?", "Yn")
			if err != nil {
				log.Fatal(err)
			}
			if validateResponse(resp) {
				isDefault = true
			}
			if err := addAccount(viper.ConfigFileUsed(), account, isDefault); err != nil {
				log.Fatal(err)
			}
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
		filePath, err := createConfigFile("exoscale")
		if err != nil {
			log.Fatal(err)
		}
		addAccount(filePath, account, true)
	}
}

func getAccount() (*account, error) {
	reader := bufio.NewReader(os.Stdin)

	account := &account{}

	name, err := readInput(reader, "Account name", "")
	if err != nil {
		return nil, err
	}

	for isAccountExist(name) {
		fmt.Printf("Account name %q already exist\n", name)
		name, err = readInput(reader, "Account name", "")
		if err != nil {
			return nil, err
		}
	}

	account.Name = name

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

	defaultZone, err := readInput(reader, "Default zone", "ch-dk-2")
	if err != nil {
		return nil, err
	}
	account.DefaultZone = defaultZone

	return account, nil
}

func isAccountExist(name string) bool {

	if allAccount == nil {
		return false
	}

	for _, acc := range allAccount.Accounts {
		if acc.Name == name {
			return true
		}
	}

	return false
}

func createConfigFile(fileName string) (string, error) {
	if _, err := os.Stat(configFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(configFolder, os.ModePerm); err != nil {
			return "", err
		}
	}

	filepath := path.Join(configFolder, fileName+".toml")

	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		return "", fmt.Errorf("File %q already exist", filepath)
	}
	return filepath, nil
}

func addAccount(filePath string, newAccount *account, isDefault bool) error {

	accountsSize := 0
	currentAccounts := []account{}
	if allAccount != nil {
		accountsSize = len(allAccount.Accounts)
		currentAccounts = allAccount.Accounts
	}

	accounts := make([]map[string]string, accountsSize+1)

	for i, acc := range currentAccounts {

		accounts[i] = map[string]string{}

		accounts[i]["name"] = acc.Name
		accounts[i]["endpoint"] = acc.Endpoint
		accounts[i]["key"] = acc.Key
		accounts[i]["secret"] = acc.Secret
		accounts[i]["defaultZone"] = acc.DefaultZone
	}

	accounts[accountsSize] = map[string]string{}

	accounts[accountsSize]["name"] = newAccount.Name
	accounts[accountsSize]["endpoint"] = newAccount.Endpoint
	accounts[accountsSize]["key"] = newAccount.Key
	accounts[accountsSize]["secret"] = newAccount.Secret
	accounts[accountsSize]["defaultZone"] = newAccount.DefaultZone

	viper.SetConfigType("toml")
	viper.SetConfigFile(filePath)

	if isDefault {
		viper.Set("defaultAccount", newAccount.Name)
	}

	viper.Set("accounts", accounts)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil

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
