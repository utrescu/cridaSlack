package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "0.0.1"
)

var (

	// Debug shows info
	Debug bool

	// Message to send
	Message string
	// Channel to send message
	Channel string
	// User who sends the message
	User string
	// Icon is the associated icon
	Icon string
)

// RootCmd is the default command
var RootCmd = &cobra.Command{
	Use:     "cridaSlack",
	Short:   "Sends a message to Slack",
	Long:    "This application sends messages to Slack via Webhook",
	Run:     func(cmd *cobra.Command, args []string) {},
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Stop execution if help or version are requested
	var helpWanted = RootCmd.Flags().Lookup("help")
	if helpWanted.Changed {
		os.Exit(0)
	}

	if RootCmd.Flags().Lookup("version").Changed {
		os.Exit(0)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Show debug info")

	RootCmd.PersistentFlags().StringVarP(&Message, "message", "m", "Test Message", "Message to send")
	RootCmd.PersistentFlags().StringVarP(&Channel, "channel", "c", "#general", "Channel to send the message")
	RootCmd.PersistentFlags().StringVarP(&User, "user", "u", "", "User who send the message")
	RootCmd.PersistentFlags().StringVarP(&Icon, "icon", "i", "", "Database user password")

	viper.BindPFlag("message", RootCmd.Flags().Lookup("message"))
	viper.BindPFlag("channel", RootCmd.Flags().Lookup("channel"))
	viper.BindPFlag("user", RootCmd.Flags().Lookup("user"))
	viper.BindPFlag("icon", RootCmd.Flags().Lookup("icon"))
}

// initConfig reads ENV variables if set. (CRIDA_*)
func initConfig() {
	viper.SetEnvPrefix("crida")
	viper.AutomaticEnv()
}
