package cmd

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/sedind/flow/app/flow/watcher"

	"github.com/sedind/flow/app/flow/config"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "runs watch process for your project",
	Run: func(c *cobra.Command, args []string) {
		ctx := context.Background()
		RunWithContext(ctx, ProjectFile)
	},
}

// RunWithContext run watcher process
func RunWithContext(ctx context.Context, cfgFile string) {
	c := &config.Configuration{}
	err := c.Load(cfgFile)
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	wg := sync.WaitGroup{}

	defer wg.Wait()

	for name, wc := range c.Watcher {
		wg.Add(1)
		go func(name string, cfg config.WatcherConfig) {
			err := startWatcher(ctx, name, cfg)
			if err != nil {
				log.Fatalln(err)
			}
		}(name, wc)
	}
}

func startWatcher(ctx context.Context, name string, c config.WatcherConfig) error {
	//fmt.Println(c.Name)
	w := watcher.NewWithContext(ctx, &c)
	return w.Start()
}
