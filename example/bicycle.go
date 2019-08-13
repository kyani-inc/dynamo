package main

import (
	"github.com/catmullet/dynamo"
	"log"
)

const tableName = "bicycle"

type Bicycle struct {
	Id       string `json:"id"`
	Gear     int    `json:"gear"`
	TireSize int    `json:"tireSize"`
}

func main() {

	bicycle := Bicycle{"123a", 12, 28}
	bicycles := []Bicycle{}

	// Save Item to Dynamo using table, the struct and ttl if desired
	err := dynamo.Put(tableName, bicycle, "")

	if err != nil {
		log.Fatal(err)
	}

	// Get single item and fill struct
	err = dynamo.GetItem(tableName, map[string]string{"id": "123a"}, &bicycle)

	if err != nil {
		log.Fatal(err)
	}

	// Get list of items with gear of 12
	err = dynamo.GetItemList(tableName, map[string]string{"gear": "12"}, &bicycles)

	if err != nil {
		log.Fatal(err)
	}

	// Delete item with id of 123a
	err = dynamo.Delete(tableName, "id", "123a")

	if err != nil {
		log.Fatal(err)
	}
}
