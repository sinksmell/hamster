package examples_test

import (
	"context"
	"fmt"
	"time"

	"github.com/sinksmell/hamster"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExampleFilterDocBuilder() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}

	opt := options.Find()
	filter := hamster.FilterDocBuilder.Gt("year", 2000).Type("imdb.rating", bsontype.Double.String()).Doc()
	sort := hamster.SortDocBuilder.OrderBy("year", hamster.SortDesc).OrderDescBy("imdb.rating").Doc()
	opt.SetSort(sort)
	opt.SetLimit(10)
	cursor, err := client.Database("sample_mflix").Collection("movies").Find(ctx, filter, opt)
	if err != nil {
		fmt.Println(err)
		return
	}

	type movie struct {
		Title string `bson:"title"`
		Year  int    `bson:"year"`
	}
	var data []movie
	defer cursor.Close(context.TODO())

	cursor.All(ctx, &data)
	fmt.Println(len(data))
	fmt.Println(data[0].Title)

	// Output:
	// 10
	// A Brave Heart: The Lizzie Velasquez Story
}
