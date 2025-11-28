package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

const TestStorageName string = "tasks_test.json"

func checkFatal(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Fatal:%s", err)
	}

}

// func TestCreateStorage(t *testing.T) {
// 	defer os.Remove(TestStorageName)
// 	_, err := getTasks(TestStorageName)
// 	checkFatal(err, t)

// 	if _, err := os.Stat(TestStorageName); errors.Is(err, os.ErrNotExist) {
// 		t.Errorf("The storage was not created")
// 	}
// 	os.Remove(TestStorageName)
// }

func TestGetTasks(t *testing.T) {
	defer os.Remove(TestStorageName)
	size := 5
	createTestStorage(t, size)
	tasks, err := getTasks(TestStorageName)
	checkFatal(err, t)

	if len(tasks) != size {
		t.Errorf("Result was incorrent, got: %d, want: %d", size, len(tasks))
	}
}

func TestSaveTask(t *testing.T) {
	size := 5
	createTestStorage(t, size)
	defer os.Remove(TestStorageName)
	tasks, err := getTasks(TestStorageName)
	checkFatal(err, t)
	task := Task{tasks[len(tasks)-1].Id + 1, "NewTask", todo, time.Now(), time.Now()}
	tasks = append(tasks, task)
	saveTasks(TestStorageName, tasks)

	tasks, err = getTasks(TestStorageName)
	checkFatal(err, t)
	addedTask := tasks[len(tasks)-1]

	if task.Id != addedTask.Id || task.Description != addedTask.Description {
		t.Errorf("Task was not added")
	}
}

func TestUpdateTaskDescription(t *testing.T) {
	size := 5
	createTestStorage(t, size)
	defer os.Remove(TestStorageName)
	tasks, err := getTasks(TestStorageName)
	checkFatal(err, t)
	id := 3
	desc := "Updated"
	err = updateTaskDescription(tasks, id, desc)
	checkFatal(err, t)
	saveTasks(TestStorageName, tasks)
	tasks, err = getTasks(TestStorageName)
	checkFatal(err, t)
	if tasks[id].Description != desc {
		t.Errorf("Description was not updated")
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	size := 5
	createTestStorage(t, size)
	defer os.Remove(TestStorageName)
	tasks, err := getTasks(TestStorageName)
	checkFatal(err, t)
	id := 3
	status := done
	err = updateTaskStatus(tasks, id, status)
	checkFatal(err, t)
	saveTasks(TestStorageName, tasks)
	tasks, err = getTasks(TestStorageName)
	checkFatal(err, t)
	if tasks[id].Status != status {
		t.Errorf("Description was not updated")
	}

}

func TestDeleteTaskById(t *testing.T) {
	size := 5
	createTestStorage(t, size)
	defer os.Remove(TestStorageName)
	tasks, err := getTasks(TestStorageName)
	checkFatal(err, t)
	id := 3
	tasks, err = deleteTaskById(tasks, id)
	checkFatal(err, t)
	saveTasks(TestStorageName, tasks)
	tasks, err = getTasks(TestStorageName)
	checkFatal(err, t)
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == id {
			t.Errorf("Task was not deleted")
		}
	}

	if len(tasks) != size-1 {
		t.Errorf("Deleted unnecessary tasks")
	}
}

// func AddTest(t *testing.T) {

// }

func createTestStorage(t *testing.T, size int) []Task {
	f, err := os.OpenFile("tasks_test.json", os.O_CREATE|os.O_RDWR, 0666)
	checkFatal(err, t)
	defer f.Close()
	tasks := make([]Task, size)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = Task{i, "Task", todo, time.Now(), time.Now()}
	}
	b, err := json.Marshal(tasks)
	checkFatal(err, t)
	f.Write(b)
	return tasks
}
