# Hamster - fluent MQL generator for Go
## Hamster is in evolution.

Inspired by:
- [squirrel](https://github.com/Masterminds/squirrel)
- [mongodb-driver-java](https://github.com/mongodb/mongo-java-driver)

MQL Builder TODO:
- [x] Filter Builder
- [x] Sort Builder
- [ ] Projection Builder
- [ ] Update Builder
- [ ] Aggregate Builder
- [ ] Index Builder


## Usage

### Filter & Sort
[Official-Filters-Builder](https://www.mongodb.com/docs/drivers/java/sync/v4.6/fundamentals/builders/filters/#std-label-filters-builders)

```go
	import "github.com/sinksmell/hamster"

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

```


## License

Hamster is released under the
[Apache License]( http://www.apache.org/licenses/).
