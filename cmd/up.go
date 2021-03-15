package cmd

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/proura/drlm2t/cfg"
	testserver "github.com/proura/drlm2t/drlm2t-server"
	"github.com/proura/drlm2t/model"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Starts all networks and VM of a test",
	Long:  `Starts all networks and VM of a test`,
	Run: func(cmd *cobra.Command, args []string) {

		init, _ := cmd.Flags().GetBool("init")
		landmark, _ := cmd.Flags().GetBool("landmark")

		log.Println("==UP CALLED==")

		//Check for test to up
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t up example")
		} else {
			_, err := os.Stat(cfg.Config.Drlm2tPath + "/tests/" + args[0])
			if os.IsNotExist(err) {
				log.Fatalln("Specified infrastructure", args[0], "does not exists")
			}
		}

		// Carregem i inicialitzem totes les configuracions
		model.LoadInfrastructure(args[0])
		model.LoadRunningInfrastructure(args[0])
		model.InitInfrastructure(args[0])

		if init {
			log.Println("+ RESET: --init/-i enabled, resetting all tests")
			model.InitTesting(1)
		} else if landmark {
			log.Println("+ LANDMARK: --landmark/-l, reset all tests from first landmark found")
			model.InitTesting(2)
		} else {
			model.InitTesting(0)
		}

		model.SaveRunningIfrastructure()

		// Arranquem el servidor per enviar i rebre les configs
		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := testserver.RunServer("config", httpServerExitDone)

		// Creem les xarxes i els hosts
		model.Infrastructure.CreateNetworks()
		model.Infrastructure.CreateHosts()

		// Esperem rebre la confirmacio de tots els hosts
		for model.Infrastructure.CountConfigRun() < model.Infrastructure.CountLocalHosts() {
			time.Sleep(1 * time.Second)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println(err.Error())
		}
		httpServerExitDone.Wait()
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().BoolP("init", "i", false, "Restart all tests")
	upCmd.Flags().BoolP("landmark", "l", false, "Run from first landmark")
	//viper.BindPFlag("doInit", rootCmd.PersistentFlags().Lookup("init"))
	//viper.BindPFlag("doMark", rootCmd.PersistentFlags().Lookup("landmark"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
