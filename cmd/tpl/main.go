package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use: "tpl",
}

var (
	logger  *zap.Logger
	slogger *zap.SugaredLogger
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "explain what is being done")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		panic(err)
	}
}

func initConfig() {
	if viper.GetBool("verbose") {
		logger, _ = zap.NewDevelopment(zap.IncreaseLevel(zap.DebugLevel))
	} else {
		logger, _ = zap.NewDevelopment(zap.IncreaseLevel(zap.WarnLevel))
	}
	slogger = logger.Sugar()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
