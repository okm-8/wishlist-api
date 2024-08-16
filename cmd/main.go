package main

import (
	"api/cmd/migrations"
	"api/cmd/system"
	"api/cmd/user"
	"api/internal/model/log"
	systemContext "api/internal/system/context"
	"github.com/spf13/cobra"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:   "wishlist-api",
	Short: "wishlist-api is a command line tool for managing wishlist API",
	Long:  "wishlist-api is a command line tool for managing wishlist API",
}

func main() {
	ctx, err := systemContext.NewContext()

	if err != nil {
		panic(err)
	}

	var group sync.WaitGroup

	defer func() {
		r := recover()

		ctx.Cancel()
		group.Wait()

		if r != nil {
			panic(r)
		}
	}()

	rootCmd.AddCommand(migrations.Root())
	rootCmd.AddCommand(system.Root())
	rootCmd.AddCommand(user.Root())

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		ctx.Log(log.Error, err.Error())
	}
}
