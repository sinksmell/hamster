# Hamster - fluent MQL generator for Go
## Hamster is in evolution.

Inspired by:
- [squirrel](https://github.com/Masterminds/squirrel)
- [mongodb-driver-java](https://github.com/mongodb/mongo-java-driver)

MQL Builder TODO:
- [x] Filter Builder
- [x] Sort Builder
- [x] Projection Builder
- [x] Update Builder
- [x] Aggregate Builder
- [x] Index Builder

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

### Projection

```go
projection := hamster.ProjectDocBuilder.Include("title", "year").ExcludeId().Doc()
opt := options.Find().SetProjection(projection)
```

### Update

```go
update := hamster.UpdateDocBuilder.
	Set("title", "Hamster 2").
	Inc("version", 1).
	CurrentDate("updated_at").
	Doc()

_, err = collection.UpdateOne(ctx, hamster.FilterDocBuilder.Eq("_id", id).Doc(), update)
```

### Aggregate

```go
pipeline := hamster.AggregateDocBuilder.
	Match(hamster.FilterDocBuilder.Gt("year", 2010).Doc().ToD()).
	Group(bson.D{{"_id", "$year"}, {"count", bson.D{{"$sum", 1}}}}).
	Sort(bson.D{{"_id", 1}}).
	Doc().ToA()

cursor, err := collection.Aggregate(ctx, pipeline)
```

### Index

```go
index := hamster.IndexDocBuilder.
	Asc("email").
	Desc("created_at").
	Unique().
	Name("email_created_idx").
	Doc().ToModel()

_, err = collection.Indexes().CreateOne(ctx, index)
```

## License

Hamster is released under the
[Apache License]( http://www.apache.org/licenses/).
