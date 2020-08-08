package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Student struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Married bool   `json:"married"`
}

func (s *Student) run() {
	fmt.Print("%s在跑。。。", s.Name)
}

func (s *Student) wang() {
	fmt.Print("%s在叫 汪汪汪。。。", s.Name)
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Println("connect to es success")
	p1 := Student{Name: "rion", Age: 90, Married: false}
	// 链式操作
	put1, err := client.Index().
		Index("student").
		Type("go").
		BodyJson(p1).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed user %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}
