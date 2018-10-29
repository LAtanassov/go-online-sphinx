package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

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

func newCommand() *cobra.Command {

	var rootCmd = cobra.Command{
		Use:   "oscli",
		Short: "Online SPHINX CLI",
		Long:  `Online SPHINX CLI is a new password mananger inspired by SPHINX`,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialzes Online SPHINX",
		Long:  `Initialzes Online SPHINX`,
		Run:   initRun,
	}

	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Registers to Online SPHINX",
		Long:  `Registers to Online SPHINX`,
		Run:   registerRun,
	}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to Online SPHINX",
		Long:  `Login to Online SPHINX`,
		Run:   loginRun,
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(loginCmd)

	return &rootCmd

}

func getConfig() *viper.Viper {
	config := viper.New()

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config.AddConfigPath(home)
	config.SetConfigName(".osctl")
	config.ReadInConfig()

	return config
}

func initRun(cmd *cobra.Command, args []string) {

	config := getConfig()
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

func registerRun(cmd *cobra.Command, args []string) {
	// TODO: load from env. or file
	cfg := client.Configuration{}
	repo := client.NewInMemoryUserRepository()
	user, err := client.NewUser("username", 8)
	if err != nil {
		return
	}
	client.New(http.DefaultClient, cfg, repo).Register(user)
}

func loginRun(cmd *cobra.Command, args []string) {
	// TODO: load from env. or file
	cfg := client.Configuration{}
	repo := client.NewInMemoryUserRepository()
	client.New(http.DefaultClient, cfg, repo).Login("username", "password")
}
