package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"time"
)

type Status int

const (
	todo Status = iota
	inprogress
	done
)

const StorageName string = "tasks.json"

const helpMessage string = `Available commands:
	add [task name] -> create the new task
	update [task id] [new task name] -> update task name by id
	delete [task id] -> delete task by id
	mark-in-progress [task id] -> change the task status to "in progress"
	mark-in-done [task id] -> change the task status to "done"
	*by default, tasks have the status "todo"*
	list -> display tasks
	list [todo|in-progress|done] -> display tasks with filtered by status
`

type Task struct {
	Id          int
	Description string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func findTaskIndex(ts []Task, id int) (int, bool) {
	n, found := slices.BinarySearchFunc(ts, Task{Id: id}, func(a, b Task) int {
		return cmp.Compare(a.Id, b.Id)
	})
	return n, found
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})
	return tasks, nil
}

func saveTasks(path string, ts []Task) {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	check(err)
	defer f.Close()
	b, err := json.Marshal(ts)
	check(err)
	f.Write(b)
}

func getStringArg(args []string, pos int, message string) string {
	var arg string
	if len(args) <= pos && message != "" {
		fmt.Print(message)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		arg = scanner.Text()
		if arg == "q" {
			os.Exit(0)
		}
		return arg
	}
	arg = os.Args[pos]
	return arg
}

func getIntArg(args []string, pos int, message string, repeat bool) int {
	var id int
	if len(args) > pos {
		var err error
		id, err = strconv.Atoi(args[pos])
		if err == nil {
			return id
		}
	}

	if len(args) <= pos && message != "" {
		fmt.Print(message)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			input := scanner.Text()
			if input == "q" {
				os.Exit(0)
			}
			tmpId, err := strconv.Atoi(input)
			if err == nil {
				id = tmpId
				break
			}
			if !repeat {
				break
			}
			fmt.Print("Invalid id, please enter again['q' for exit]:")
		}
	}

	return id
}

func exitIfNotExists(ts []Task) {
	if len(ts) == 0 {
		fmt.Println("No tasks, nothing to do")
		os.Exit(0)
	}
}

func updateTaskDescription(ts []Task, id int, desc string) error {
	if index, ok := findTaskIndex(ts, id); ok {
		ts[index].Description = desc
		ts[index].UpdatedAt = time.Now()
		return nil
	}
	return errors.New("Task index not found")
}

func updateTaskStatus(ts []Task, id int, status Status) error {
	if index, ok := findTaskIndex(ts, id); ok {
		ts[index].Status = status
		ts[index].UpdatedAt = time.Now()
		return nil
	}
	return errors.New("Task index not found")
}

func deleteTaskById(ts []Task, id int) ([]Task, error) {
	index, ok := findTaskIndex(ts, id)
	if !ok {
		return nil, errors.New("Task index not found")
	}

	// O(N)
	ts = slices.Delete(ts, index, index+1)
	// O(1)
	// ts[index], ts[len(ts)-1] = ts[len(ts)-1], ts[index]
	// ts = ts[:len(ts)-1]
	return ts, nil
}

func main() {
	statuses := map[Status]string{todo: "todo", inprogress: "in-progress", done: "done"}
	strStatuses := map[string]Status{"todo": todo, "in-progress": inprogress, "done": done}

	tasks, err := getTasks(StorageName)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if len(os.Args) == 1 {
		fmt.Print(helpMessage)
		os.Exit(0)
	}
	action, args := os.Args[1], os.Args[2:]

	switch action {
	case "add":
		name := getStringArg(os.Args, 2, "Enter the task name:")
		maxId := 1
		if len(tasks) > 0 {
			maxId = tasks[len(tasks)-1].Id + 1
		}

		task := Task{maxId, name, todo, time.Now(), time.Now()}
		tasks = append(tasks, task)
		saveTasks(StorageName, tasks)

	case "update":
		exitIfNotExists(tasks)
		id := getIntArg(os.Args, 2, "Enter task id ['q' for exit]:", true)
		description := getStringArg(os.Args, 3, "Enter task description ['q' for exit]")
		err := updateTaskDescription(tasks, id, description)
		check(err)
		saveTasks(StorageName, tasks)

	case "mark-in-progress", "mark-done":
		exitIfNotExists(tasks)
		id := getIntArg(os.Args, 2, "Enter task id ['q' for exit]:", true)
		status := inprogress
		if action == "mark-done" {
			status = done
		}
		err := updateTaskStatus(tasks, id, status)
		check(err)
		saveTasks(StorageName, tasks)

	case "delete":
		exitIfNotExists(tasks)
		id := getIntArg(os.Args, 2, "Enter task id ['q' for exit]:", true)
		tasks, err = deleteTaskById(tasks, id)
		check(err)
		saveTasks(StorageName, tasks)

	case "list":
		var status Status
		if len(args) != 0 {
			status = strStatuses[args[0]]
		}
		for _, task := range tasks {
			if status != 0 && task.Status != status {
				continue
			}
			fmt.Println(task.Id, task.Description, statuses[task.Status], task.CreatedAt.Format(time.DateTime), task.UpdatedAt.Format(time.DateTime))
		}
	case "help":
		fmt.Print(helpMessage)
	default:
		fmt.Println("The command doesn't exists")
		fmt.Print(helpMessage)
		os.Exit(1)
	}
}
