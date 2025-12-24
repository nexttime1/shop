package main

import (
	"context"
	"github.com/olivere/elastic/v7"
)

func main() {
	url := "http://192.168.163.50:9200"

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	q := elastic.NewMatchQuery("address", "street")
	client.Search().Query(q).Do(context.Background())

}
