package cmd

import (
	"github.com/sejo412/ya-boo/internal/app"
	"github.com/sejo412/ya-boo/internal/db"
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
		var storage app.Storage
		storage = db.NewPostgres()
		a := app.NewApp(cfg, &storage)
		if err := a.Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntP("port", "p", defaultPort, "http server port")
}
