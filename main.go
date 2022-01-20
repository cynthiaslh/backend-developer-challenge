package main

import (
	"time"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"encoding/csv"
	"encoding/json"

)

type Item struct {
	Name      	string  `json:"name"`
	ID        	string  `json:"id"`
	Tracked   	bool    `json:"tracked"`
	Cost 		float64 `json:"cost"`
	CreatedAt	string  `json:"created_at"`
	UpdatedAt	string  `json:"updated_at"`
}

var items = []Item{
	{Name: "History Book", ID: "1", Tracked: true, Cost: 20.50, CreatedAt: "2021-08-24T14:01:47-04:00", UpdatedAt: "2021-08-24T14:01:47-04:00"},
	{Name: "Winter boots", ID: "2", Tracked: true, Cost: 102.20, CreatedAt: "2021-08-29T14:01:47-04:00", UpdatedAt: "2021-09-24T14:01:47-04:00"},
}

func main() {
	r := gin.Default()

	r.POST("/item", createItem)
	r.PATCH("/item/:id", editItem)
	r.DELETE("/item/:id", deleteItem)
	r.GET("/item", viewItems)
	r.GET("/item/csv", exportToCSV)

	r.Run() 
}

func createItem(g *gin.Context) {
	var newItem Item

	if err := g.BindJSON(&newItem); err != nil {
		return
	}
	newItem.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
	newItem.UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
	
	items = append(items, newItem)
	g.IndentedJSON(http.StatusCreated, newItem)
}

func getItemByID(g *gin.Context) {
	id := g.Param("id")

	for _, i := range items {
		if i.ID == id {
			g.IndentedJSON(http.StatusOK, i)
			return
		}
	}
	g.IndentedJSON(http.StatusNotFound, gin.H{"message": "item not found"})
}


func editItem(g *gin.Context) {
	var newItem Item
	id := g.Param("id")

	if err := g.BindJSON(&newItem); err != nil {
		return
	}

	for index, item := range items {
		if item.ID == id {
			items[index].Name = newItem.Name
			items[index].Tracked = newItem.Tracked
			items[index].Cost = newItem.Cost
			items[index].UpdatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z07:00")
			g.IndentedJSON(http.StatusOK, items[index])
			return
		}
	}
	g.IndentedJSON(http.StatusNotFound, gin.H{"message": "Item does not exist"})
}

func deleteItem(g *gin.Context) {
	id := g.Param("id")
	for index, i := range items {
		if i.ID == id {
			items = append(items[:index], items[index+1:]...)
		}
	}
}

func viewItems(g *gin.Context) {
	g.IndentedJSON(http.StatusOK, items)
}

func exportToCSV(g *gin.Context) {
	file, err := os.Create("items.csv")

	if err != nil {
		g.IndentedJSON(http.StatusNotFound, gin.H{"message": "Failed creating the CSV file"})
	}

    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, value := range items {
		item, err := json.Marshal(value)
		if err != nil {
			return
		}
		item_str := []string{string(item)}
        err = writer.Write(item_str)
		if err != nil {
            return 
        }
    }
	//downloaded the csv file
	g.FileAttachment("items.csv", "items.csv")
}


