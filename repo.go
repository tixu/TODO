package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

var currentId int64
var todos Todos
var database *bolt.DB
var tofo = []byte("tofo")

// Give us some seed data
func init() {
	var err error
	database, err = bolt.Open("bolt.db", 0644, nil)
	if err != nil {
		panic(err)
	}
	RepoCreateTodo(Todo{Name: "Write presentation"})
	RepoCreateTodo(Todo{Name: "Host meetup"})
}

func RepoFindTodo(id int64) Todo {
	var task Todo
	err := database.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(tofo)
		if bucket == nil {
			panic(fmt.Errorf("Bucket %q not found!", tofo))
		}

		var key []byte = make([]byte, 10)
		binary.PutVarint(key, id)
		val := bucket.Get(key)

		err := json.Unmarshal(val, &task)
		if err != nil {
			panic(err)
		}
		return nil

	})

	if err != nil {
		return Todo{}
	}
	return task
}

func RepoCreateTodo(t Todo) (Todo, error) {
	currentId += 1
	t.Id = currentId
	todos = append(todos, t)

	dbErr := database.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(tofo)
		if err != nil {
			return err
		}
		dbt, marshallErr := json.Marshal(t)
		if marshallErr != nil {
			return marshallErr
		}
		var key []byte = make([]byte, 10)
		binary.PutVarint(key, t.Id)
		err = bucket.Put(key, dbt)
		if err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		panic(dbErr)
	}

	return t, nil
}

func RepoDestroyTodo(id int64) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
