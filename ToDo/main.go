package main

import (
	"ToDo/todos"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	//'add' subcommand
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	//inputs for 'add' command
	addName := addCmd.String("name", "", "Task name")
	addType := addCmd.String("type", "work", "Task type")
	addImportance := addCmd.Int("impo", todos.DefaultImportance, "Task importance from 0 to 10")

	//'print' subcommand
	printToDo := flag.Bool("print", false, "Print ToDo list")

	//'del' delete subcommand
	delTask := flag.Int("del", 0, "Specify which task to delete")

	//'done' mark as done subcommand
	doneTask := flag.Int("done", 0, "Specify which task to mark as complete")

	if len(os.Args) < 2 {
		printExpected()
		os.Exit(1)
	}

	flag.Parse()

	switch {
	case *printToDo:
		HandlePrint()

	case *delTask > 0:
		HandleDel(delTask)

	case *doneTask > 0:
		HandleDone(doneTask)

	default:
		HandleArgs(addCmd, addName, addType, addImportance)
	}
}

func HandleArgs(addCmd *flag.FlagSet, title *string, taskType *string, importance *int) {

	addCmd.Parse(os.Args[2:])

	switch os.Args[1] {
	case "add":
		HandleAdd(addCmd, title, taskType, importance)

	default:
		printExpected()
	}
}

func HandleAdd(addCmd *flag.FlagSet, title *string, taskType *string, importance *int) {

	if strings.TrimSpace(*title) == "" {
		fmt.Println("title can not be empty")
		addCmd.PrintDefaults()
		os.Exit(1)
	} else if len(strings.TrimSpace(*taskType)) > 10 {
		fmt.Println("task type can not be longer than 10")
		addCmd.PrintDefaults()
		os.Exit(1)
	} else if (*importance) > 10 || (*importance) < 0 {
		fmt.Println("task importance has to be between 0 and 10")
		addCmd.PrintDefaults()
		os.Exit(1)
	} else {
		todos.AddTask(title, taskType, importance)
	}
}

func HandlePrint() {
	todos.PrintToDo()
}

func HandleDel(taskToDelete *int) {
	todos.Delete(*taskToDelete)
}

func HandleDone(handleDone *int) {
	todos.MarkAsComplete(*handleDone)
}

func HandleDefault() {
	fmt.Print("Default!")
}

func printExpected() {
	fmt.Println("expected 'add', 'print', 'del' or 'done' subcommands")
	fmt.Printf("instead got: %s\n", strings.Join(os.Args[1:], " "))
	os.Exit(1)
}
