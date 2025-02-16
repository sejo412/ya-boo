package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sejo412/ya-boo/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Update database scheme",
	Long:  "Update database scheme\n",
	Run: func(cmd *cobra.Command, args []string) {
		v := viper.New()
		if err := v.BindPFlags(cmd.Flags()); err != nil {
			panic(err)
		}
		cfg := config.NewConfig()
		if err := cfg.Load(v); err != nil {
			panic(err)
		}
		m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), cfg.Dsn)
		if err != nil {
			panic(err)
		}
		if err := m.Up(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				panic(err)
			} else {
				log.Println("no change")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
