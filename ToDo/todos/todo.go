package todos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const FileName = "./tasks.json"
const DefaultImportance = 5

type Task struct {
	Name          string
	Type          string
	Importance    int
	Done          bool
	WhenCreated   time.Time
	WhenCompleted time.Time
}

func (task Task) getAsString(index int) *[]string {

	//create array with number of fields from Task
	values := make([]string, reflect.TypeOf(Task{}).NumField()+1)

	v := reflect.ValueOf(task)
	reflectType := v.Type()

	values[0] = strconv.Itoa(index)

	for i := 0; i < v.NumField(); i++ {

		tmpString := fmt.Sprint(v.Field(i).Interface())

		switch reflectType.Field(i).Name {
		case "WhenCompleted":
			if tmpString == "0000-01-01 01:00:00 +0000 UTC" {
				values[i+1] = ""
			} else {
				values[i+1] = tmpString[:19]
			}

		case "WhenCreated":
			values[i+1] = tmpString[:19]

		case "Done":
			if tmpString == "false" {
				values[i+1] = "❌"
			} else {
				values[i+1] = "✔️"
			}

		default:
			values[i+1] = tmpString
		}
	}
	return &values
}

func AddTask(title *string, taskType *string, importance *int) {
	toDoList := LoadFromFile()

	newTask := Task{
		Name:          *title,
		Type:          *taskType,
		Done:          false,
		Importance:    *importance,
		WhenCreated:   time.Now(),
		WhenCompleted: time.Date(0, time.Month(1), 1, 1, 0, 0, 0, time.UTC),
	}

	toDoList = append(toDoList, newTask)
	saveTasks(&toDoList)
}

func saveTasks(tasks *[]Task) {

	videoBytes, err := json.Marshal(tasks)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(FileName, videoBytes, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadFromFile() (toDoList []Task) {

	fileBytes, err := ioutil.ReadFile(FileName)
	if err != nil {
		panic(err)
	}

	if len(fileBytes) == 0 {
		panic(err)
	}

	err = json.Unmarshal(fileBytes, &toDoList)

	if err != nil {
		panic(err)
	}

	return toDoList
}

func Delete(index int) {
	toDoList := LoadFromFile()

	if index <= 0 || index > len(toDoList) {

		fmt.Printf("index should be between 1 and %d\n", len(toDoList))
	} else {
		slice := toDoList
		toDoList = append(slice[:index-1], slice[index:]...)

		saveTasks(&toDoList)
		fmt.Printf("Task number %d has been deleted!\n", index)
	}
}

func MarkAsComplete(index int) {
	toDoList := LoadFromFile()
	index = index - 1

	if index < 0 || index > len(toDoList) {
		fmt.Printf("index should be between 1 and %d\n", len(toDoList))
	} else {
		(&toDoList[index]).Done = true
		(&toDoList[index]).WhenCompleted = time.Now()
		saveTasks(&toDoList)
		fmt.Printf("Task number %d has been marked complete!\n", index+1)
	}
}

func PrintToDo() {
	tasks := LoadFromFile()

	titles := []string{"ID", "Title", "Type", "Importance", "Done", "Created ", "Completed "}
	maxLengths := []int{4, 36, 14, 4, 12, 34, 34}
	centering := []int{0, 2, 3, 4, 5, 6}

	content := make([][]string, len(tasks))

	for index, task := range tasks {
		content[index] = *task.getAsString(index + 1)
	}

	PrintTable(&titles, &maxLengths, &centering, &content)
}

func PrintTable(titles *[]string, maxLengths *[]int, centering *[]int, content *[][]string) {
	// titles := []string{"ID", "Title", "Importance", "Done", "Created", "Completed"}
	// maxLengths := []int{4, 30, 4, 4, 24, 24}
	// centering := []int{0, 2}

	// content := [][]string{
	// 	{"1", "test1", "10", "yes", "2022-08-02 00:00:00 AM", "2022-08-02 00:00:00 AM"},
	// 	{"100", "test2 awdwd awd awdawd wdawdaw dawd wd awd w", "2", "no", "2022-08-02 00:00:00 AM", "2022-08-02 00:00:00 AM"},
	// }

	if len(*titles) != len(*maxLengths) || len(*titles) != len((*content)[0]) {
		fmt.Println("Invalid length of data")
		return
	}

	var calculatedLengths = make([]int, len(*titles))
	var top = make([]string, len(*titles))

	//Calculating max length of columns
	for index, value := range *titles {

		//we want to have spaces in table
		if len(value)+4 >= (*maxLengths)[index] {
			calculatedLengths[index] = len(value) + 4
		} else {
			calculatedLengths[index] = (*maxLengths)[index]
		}

		//if number is not even we need to add one
		if len((*titles)[index])%2 != 0 {
			calculatedLengths[index] += 1
		}

		top[index] = strings.Repeat("─", calculatedLengths[index])
	}

	//Building top
	//┌──────────┬─────────┬──────────────┬────────┬───────────┬─────────────┐
	for index, value := range top {
		if index == 0 {
			fmt.Printf("┌%s", value)
		} else if index == len(top)-1 {
			fmt.Printf("┬%s┐", value)
		} else {
			fmt.Printf("┬%s", value)
		}
	}
	fmt.Printf("\n")

	//Printing Titles values
	//│  ID  │  Title  │  Importance  │  Done  │  Created  │  Completed  │
	for index, value := range *titles {
		numOfSpaces := (calculatedLengths[index] - len((*titles)[index])) / 2
		spaces := strings.Repeat(" ", numOfSpaces)

		if index == 0 {
			fmt.Printf("│%s%s%s", spaces, value, spaces)
		} else if index == len(top)-1 {
			fmt.Printf("│%s%s%s│", spaces, value, spaces)
		} else {
			fmt.Printf("│%s%s%s", spaces, value, spaces)
		}

	}
	fmt.Printf("\n")

	//Printing Content
	for i, row := range *content {

		//Top of each cell
		//├─────┼─────┼────┼─────┼────┼────┤
		for index, value := range top {
			if index == 0 {
				fmt.Printf("├%s", value)
			} else if index == len(top)-1 {
				fmt.Printf("┼%s┤", value)
			} else {
				fmt.Printf("┼%s", value)
			}
		}
		fmt.Printf("\n")

		//Print each row
		//│  1   │test1 │10 │yes  │2022-08-02 00:00:00 AM   │2022-08-02 00:00:00 AM   │
		for j, value := range row {

			var minLen = 0
			valueToPrint := ""

			valueLen := 0
			if value == "✔️" {
				valueLen = 1
			} else if value == "❌" {
				valueLen = 2
			} else {
				valueLen = len(value)
			}

			//if some content is too long we add dots at the end
			//│ 100  │test2 awdwd awd awdawd wd...   │2 │no │2022-08-02 00:00:00 AM  │2022-08-02 00:00:00 AM  │
			if valueLen >= calculatedLengths[j] {
				minLen = calculatedLengths[j] - 3
				valueToPrint = value[:minLen-3] + "..."
			} else {
				minLen = valueLen
				valueToPrint = value[:minLen]
			}

			if value == "✔️" || value == "❌" {
				valueToPrint = value
			}

			spaces := strings.Repeat(" ", calculatedLengths[j]-minLen)

			//If centering
			if sliceContains(centering, j) {

				numOfSpaces := (calculatedLengths[j] - valueLen) / 2

				spaces = strings.Repeat(" ", numOfSpaces)

				if valueLen%2 != 0 {
					fmt.Printf("│%s%s%s ", spaces, valueToPrint, spaces)
				} else {
					fmt.Printf("│%s%s%s", spaces, valueToPrint, spaces)
				}

				if j == len(top)-1 {
					fmt.Printf("│")
				}

			} else {

				if j == len(top)-1 {
					fmt.Printf("│%s%s│", valueToPrint, spaces)
				} else {
					fmt.Printf("│%s%s", valueToPrint, spaces)
				}
			}

		}
		fmt.Printf("\n")

		//Printing bottom
		//└────┴─────┴─────┴────┴────┴───┘
		if i == len(*content)-1 {
			for index, value := range top {
				if index == 0 {
					fmt.Printf("└%s", value)
				} else if index == len(top)-1 {
					fmt.Printf("┴%s┘", value)
				} else {
					fmt.Printf("┴%s", value)
				}
			}
			fmt.Printf("\n")
		}

	}

	// fmt.Println("┌──────────┬─────────┬──────────────┬────────┬───────────┬─────────────┐")
	// fmt.Println("│ ID       │  Title  │  Importance  │  Done  │  Created  │  Completed  │")
	// fmt.Println("├──────────┼─────────┼──────────────┼────────┼───────────┼─────────────┤")

	/*
		┌──────┬───────────────────────────────┬──────────────┬────────┬─────────────────────────┬─────────────────────────┐
		│  ID  │             Title             │  Importance  │  Done  │         Created         │        Completed        │
		├──────┼───────────────────────────────┼──────────────┼────────┼─────────────────────────┼─────────────────────────┤
		│  1   │test1                          │10            │yes     │2022-08-02 00:00:00 AM   │2022-08-02 00:00:00 AM   │
		├──────┼───────────────────────────────┼──────────────┼────────┼─────────────────────────┼─────────────────────────┤
		│ 100  │test2 awdwd awd awdawd wd...   │2             │no      │2022-08-02 00:00:00 AM   │2022-08-02 00:00:00 AM   │
		└──────┴───────────────────────────────┴──────────────┴────────┴─────────────────────────┴─────────────────────────┘
	*/

}

func sliceContains(slice *[]int, number int) bool {
	//function that checks if slice contains number

	for _, value := range *slice {
		if value == number {
			return true
		}
	}
	return false
}
