package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/exoscale/egoscale"

	"github.com/go-ini/ini"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	exoConfigFileName = "exoscale"
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

	if viper.ConfigFileUsed() != "" {
		println("Good day! exo is already configured with accounts:")
		listAccounts()
		if askQuestion("Do you wish to add another account?") {
			if err := addNewAccount(); err != nil {
				log.Fatal(err)
			}
		}
		return
	}
	csPath, ok := isCloudstackINIFileExist()
	if ok {
		resp, ok, err := askCloudstackINIMigration(csPath)
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			if err := createAccount(false); err != nil {
				log.Fatal(err)
			}
			return
		}

		cfgPath, err := createConfigFile(exoConfigFileName)
		if err != nil {
			log.Fatal(err)
		}
		if err := importCloudstackINI(resp, csPath, cfgPath); err != nil {
			log.Fatal(err)
		}
		if !askQuestion("Do you wish to add another account?") {
			return
		}
		if err := addNewAccount(); err != nil {
			log.Fatal(err)
		}
		return
	}
	println("Hi happy Exoscalian, some configuration is required to use exo")
	println(`
We now need some very important informations, find them there.
https://portal.exoscale.com/account/profile/api
`)
	createAccount(false)

}

func addNewAccount() error {
	newAccount, err := getAccount()
	if err != nil {
		return err
	}
	isDefault := false
	if askQuestion("Make " + newAccount.Name + " your current profile?") {
		isDefault = true
	}
	return addAccount(viper.ConfigFileUsed(), &config{DefaultAccount: newAccount.Name, Accounts: []account{*newAccount}}, isDefault)
}

func createAccount(askDefault bool) error {
	newAccount, err := getAccount()
	if err != nil {
		return err
	}
	filePath, err := createConfigFile(exoConfigFileName)
	if err != nil {
		return err
	}

	if askDefault {
		if !askQuestion("Make " + newAccount.Name + " your current profile?") {
			askDefault = false
		}
	}

	return addAccount(filePath, &config{DefaultAccount: newAccount.Name, Accounts: []account{*newAccount}}, askDefault)
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

	accountResp, err := checkCredentials(account)
	if err != nil {
		fmt.Printf("Account %q: unable to verify user credentials\n", account.Name)
	}

	for err != nil {
		account, err = getAccount()
		if err == nil {
			return account, nil
		}
	}

	account.Account = accountResp.Name

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

func addAccount(filePath string, newAccounts *config, isDefault bool) error {

	accountsSize := 0
	currentAccounts := []account{}
	if allAccount != nil {
		accountsSize = len(allAccount.Accounts)
		currentAccounts = allAccount.Accounts
	}

	newAccountsSize := 0

	if newAccounts != nil {
		newAccountsSize = len(newAccounts.Accounts)
	}

	accounts := make([]map[string]string, accountsSize+newAccountsSize)

	conf := &config{}

	for i, acc := range currentAccounts {

		accounts[i] = map[string]string{}

		accounts[i]["name"] = acc.Name
		accounts[i]["endpoint"] = acc.Endpoint
		accounts[i]["key"] = acc.Key
		accounts[i]["secret"] = acc.Secret
		accounts[i]["defaultZone"] = acc.DefaultZone
		accounts[i]["account"] = acc.Account

		conf.Accounts = append(conf.Accounts, acc)
	}

	if newAccounts != nil {

		for i, acc := range newAccounts.Accounts {

			accounts[accountsSize+i] = map[string]string{}

			accounts[accountsSize+i]["name"] = acc.Name
			accounts[accountsSize+i]["endpoint"] = acc.Endpoint
			accounts[accountsSize+i]["key"] = acc.Key
			accounts[accountsSize+i]["secret"] = acc.Secret
			accounts[accountsSize+i]["defaultZone"] = acc.DefaultZone
			accounts[accountsSize+i]["account"] = acc.Account
			conf.Accounts = append(conf.Accounts, acc)
		}
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(filePath)

	if isDefault {
		viper.Set("defaultAccount", newAccounts.DefaultAccount)
	}

	viper.Set("accounts", accounts)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	conf.DefaultAccount = viper.Get("defaultAccount").(string)
	allAccount = conf

	return nil

}

func isCloudstackINIFileExist() (string, bool) {

	envConfigPath := os.Getenv("CLOUDSTACK_CONFIG")

	usr, _ := user.Current()

	localConfig, _ := filepath.Abs("cloudstack.ini")
	inis := []string{
		localConfig,
		filepath.Join(usr.HomeDir, ".cloudstack.ini"),
		filepath.Join(configFolder, "cloudstack.ini"),
		envConfigPath,
	}

	cfgPath := ""

	for _, i := range inis {
		if _, err := os.Stat(i); err != nil {
			continue
		}
		cfgPath = i
		break
	}

	if cfgPath == "" {
		return "", false
	}
	return cfgPath, true
}

func askCloudstackINIMigration(csFilePath string) (string, bool, error) {

	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, csFilePath)
	if err != nil {
		return "", false, err
	}

	if len(cfg.Sections()) <= 0 {
		return "", false, nil
	}

	println("We've found the following configurations:")
	for i, acc := range cfg.Sections() {
		if i == 0 {
			continue
		}
		fmt.Printf("- [%s] %s\n", acc.Name(), acc.Key("key").String())
	}

	reader := bufio.NewReader(os.Stdin)

	resp, err := readInput(reader, "Do you wish to import them automagically?", "All, some, none")
	if err != nil {
		return "", false, err
	}

	resp = strings.ToLower(resp)

	return resp, (resp == "all" || resp == "some"), nil
}

func importCloudstackINI(option, csPath, cfgPath string) error {
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, csPath)
	if err != nil {
		return err
	}

	config := &config{}

	for i, acc := range cfg.Sections() {
		if i == 0 {
			continue
		}

		if option == "some" {
			if !askQuestion(fmt.Sprintf("Importing %s %s?", acc.Name(), acc.Key("key").String())) {
				continue
			}
		}

		reader := bufio.NewReader(os.Stdin)

		defaultZone, err := readInput(reader, fmt.Sprintf("%s default zone", acc.Name()), "ch-dk-2")
		if err != nil {
			return err
		}

		isDefault := false
		if askQuestion("Make " + acc.Name() + " your current profile?") {
			isDefault = true
		}

		account := account{
			Name:        acc.Name(),
			Endpoint:    acc.Key("endpoint").String(),
			Key:         acc.Key("key").String(),
			Secret:      acc.Key("secret").String(),
			DefaultZone: defaultZone,
		}

		accountResp, err := checkCredentials(&account)
		if err != nil {
			fmt.Printf("Account %q: unable to verify user credentials\n", acc.Name())
			if !askQuestion("Do you want to keep this account?") {
				continue
			}
		}

		account.Account = accountResp.Name

		config.Accounts = append(config.Accounts, account)

		if i == 1 || isDefault {
			config.DefaultAccount = acc.Name()
		}
	}

	addAccount(cfgPath, config, true)

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

func askQuestion(text string) bool {

	reader := bufio.NewReader(os.Stdin)

	resp, err := readInput(reader, text, "Yn")
	if err != nil {
		log.Fatal(err)
	}

	return (strings.ToLower(resp) == "y" || strings.ToLower(resp) == "yes")
}

func checkCredentials(account *account) (*egoscale.Account, error) {
	cs := egoscale.NewClient(account.Endpoint, account.Key, account.Secret)

	resp, err := cs.Request(&egoscale.ListAccounts{})
	if err != nil {
		return nil, err
	}

	accountsResp := resp.(*egoscale.ListAccountsResponse)

	if accountsResp.Count == 1 {
		return &accountsResp.Account[0], nil
	}

	return nil, fmt.Errorf("more than one account found")
}

func listAccounts() {
	for _, acc := range allAccount.Accounts {
		print("- ", acc.Name)
		if acc.Name == allAccount.DefaultAccount {
			print(" [current]")
		}
		println("")
	}
}

func init() {

	configCmd.Run = configCmdRun
	configCmd.Flags().BoolP("print", "p", false, "Print configuration")
	RootCmd.AddCommand(configCmd)
}
