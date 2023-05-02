package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

func getSubscriberIDs() ([]int64, error) {
	file, err := os.ReadFile(DATA_FILE)
	if err != nil {
		return nil, err
	}

	var ids []int64
	if err = json.Unmarshal(file, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func addNewSubscriber(id int64) error {
	ids, err := getSubscriberIDs()
	if err != nil {
		return err
	}

	ids = append(ids, id)
	file, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	log.Println("saving id: " + strconv.FormatInt(id, 10))

	return os.WriteFile(DATA_FILE, file, 0644)
}
