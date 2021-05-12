package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/KubaiDoLove/gophercises_cli/db"
	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Marks a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("Failed to parse the argument: %s", arg)
			}

			ids = append(ids, id)
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			os.Exit(1)
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number: ", id)
				continue
			}
			task := tasks[id-1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
				continue
			}
			fmt.Printf("Marked \"%d\" as completed.\n", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
