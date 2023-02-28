package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/mauricioabreu/golings/golings/exercises"
	"github.com/mauricioabreu/golings/golings/ui"
	"github.com/spf13/cobra"
)

func WatchCmd(infoFile string) *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Run a single exercise",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Create watcher")
			reader := bufio.NewReader(os.Stdin)
			update := make(chan bool)

			for {
				go WatchEvents(update)

				for <-update {
					exercise, err := exercises.NextPending(infoFile)
					result, err := exercise.Run()

					if err != nil {
						color.Cyan("Failed to compile the exercise %s\n\n", result.Exercise.Path)
						color.White("Check the output below: \n\n")
						color.Red(result.Err)
						color.Red(result.Out)
						color.Yellow("If you feel stuck, ask a hint by executing `golings hint %s`", result.Exercise.Name)
					} else {
						color.Green("Congratulations!\n\n")
						color.Green("Here is the output of your program:\n\n")
						color.Cyan(result.Out)
						if result.Exercise.State() == exercises.Pending {
							color.White("Remove the 'I AM NOT DONE' from the file to keep going\n")
							return fmt.Errorf("exercise is still pending")
						}
					}
				}

				color.Yellow("$")
				cmdString, err := reader.ReadString('\n')

				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}

				cmdStr := strings.TrimSuffix(cmdString, "\n")
				log.Println("CMDr")
				switch cmdStr {
				case "list":
					log.Println("List command", cmdString)
					exs, err := exercises.List(infoFile)
					if err != nil {
						color.Red(err.Error())
						os.Exit(1)
					}
					ui.PrintList(os.Stdout, exs)
				default:
					fmt.Errorf("ERROR :/")
				}
			}
		},
	}
}

func WatchEvents(updateF chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	path, _ := os.Getwd()
	file_path := fmt.Sprintf("%s/exercises", path)

	err = filepath.WalkDir(file_path, func(path_dir string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if d.IsDir() {
			err = watcher.Add(path_dir)

			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error in file path:", err.Error())
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					updateF <- true
				}
			}
		}
	}()

	// Block main goroutine forever.
	<-make(chan struct{})
}
