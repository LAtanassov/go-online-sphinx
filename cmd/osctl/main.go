package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cli struct {
	clt *client.Client
}

func newCommand() *cobra.Command {

	tr := &http.Transport{
		IdleConnTimeout: 30 * time.Second,
	}

	c := cli{
		client.New(
			&http.Client{Transport: tr},
			getConfiguration(),
			client.NewSQLiteUserRepository(),
		),
	}

	var rootCmd = cobra.Command{
		Use:   "oscli",
		Short: "Online SPHINX CLI",
		Long:  `Online SPHINX CLI is a new password mananger inspired by SPHINX`,
	}

	var registerCmd = &cobra.Command{
		Use:   "register <username>",
		Short: "Registers to Online SPHINX",
		Long:  `Registers to Online SPHINX using your username.`,
		Run:   c.registerRun,
	}

	var loginCmd = &cobra.Command{
		Use:   "login <username> <password>",
		Short: "Login",
		Long: `
			Login to Online SPHINX using your username and password. 
			TODO: passwords should not be handled in CLI like that.`,
		Run: c.loginRun,
	}

	var addCmd = &cobra.Command{
		Use:   "add <domain>",
		Short: "Add new Domain",
		Long:  `Add new Domain`,
		Run:   c.addRun,
	}

	var getCmd = &cobra.Command{
		Use:   "get <domain>",
		Short: "Get Password of Domain",
		Long:  `Get Password of Domain`,
		Run:   c.getRun,
	}

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)

	return &rootCmd

}

func getConfiguration() client.Configuration {
	config := viper.New()

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config.AddConfigPath(home)
	config.SetConfigName(".osctl")
	config.ReadInConfig()

	return client.Configuration{}
}

func initRun(cmd *cobra.Command, args []string) {

	config := viper.New()
	config.Set("app.name", "online-sphinx")
	config.Set("app.version", "0.1.0")

	config.Set("client.id", "client")
	config.Set("client.protocol.bit", 8)

	k, err := rand.Int(rand.Reader, big.NewInt(256))
	if err != nil {
		log.Fatal(err)
	}
	config.Set("client.secret.k", k.Text(16))
	config.Set("client.protocol.hash", "sha256")

	q, err := rand.Prime(rand.Reader, 8)
	if err != nil {
		log.Fatal(err)
	}
	config.Set("client.protocol.q", q.Text(16))

	config.Set("server.address", ":8080")

	err = config.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func (c *cli) registerRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(-1)
	}

	err := c.clt.Register(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (c *cli) loginRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		cmd.Help()
		os.Exit(-1)
	}

	err := c.clt.Login(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (c *cli) addRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(-1)
	}

	err := c.clt.Add(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func (c *cli) getRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Help()
		os.Exit(-1)
	}

	pwd, err := c.clt.Get(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(pwd)
}
