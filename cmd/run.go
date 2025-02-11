package cmd

import (
	"github.com/sejo412/ya-boo/internal/app"
	"github.com/sejo412/ya-boo/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run ya-boo",
	Long:  "Run ya-boo\n",
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		if err := v.BindPFlags(cmd.Flags()); err != nil {
			panic(err)
		}
		cfg := config.NewConfig()
		if err := cfg.Load(v); err != nil {
			panic(err)
		}
		a := app.NewApp(cfg)
		if err := a.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("config", "c", "", "config file path")
	runCmd.Flags().IntP("port", "p", defaultPort, "http server port")
	runCmd.Flags().StringP("dsn", "d", "", "database connection string")
}
