package cmd

import (
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialzes Online SPHINX",
	Long:  `Initialzes Online SPHINX`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	C, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		log.Fatal(err)
	}
	initCmd.PersistentFlags().StringP("client.id", "c", C.Text(16), "defines the client id")
	viper.BindPFlag("client.id", rootCmd.PersistentFlags().Lookup("client.id"))

	k, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		log.Fatal(err)
	}
	initCmd.PersistentFlags().StringP("secret.k", "k", k.Text(16), "defines the ElGamal private key")
	viper.BindPFlag("secret.k", rootCmd.PersistentFlags().Lookup("secret.k"))

	initCmd.PersistentFlags().StringArrayP("server.address", "s", []string{":8080"}, "defines all online sphinx servers")
	viper.BindPFlag("server.address", rootCmd.PersistentFlags().Lookup("server.address"))

	initCmd.PersistentFlags().StringP("hash.algo", "a", "sha256", "defines hash algorithmn used within the protocol")
	viper.BindPFlag("hash.algo", rootCmd.PersistentFlags().Lookup("hash.algo"))

	b := initCmd.PersistentFlags().IntP("bit.length", "b", 32, "defines the bit length of cyclic group")
	viper.BindPFlag("bit.length", rootCmd.PersistentFlags().Lookup("bit.length"))

	q, err := rand.Prime(rand.Reader, *b)
	if err != nil {
		log.Fatal(err)
	}
	initCmd.PersistentFlags().StringP("prime.q", "q", q.Text(16), "specifies prime q used to calculate cyclic group")
	viper.BindPFlag("prime.q", rootCmd.PersistentFlags().Lookup("prime.q"))
}
