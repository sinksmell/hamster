package hamster

import (
	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// indexDoc defines one index model
type indexDoc struct {
	Keys    bson.D
	Options *options.IndexOptions
}

type indexDocBuilder builder.Builder

var (
	// IndexDocBuilder is a singleton builder for indexDoc
	IndexDocBuilder = builder.Register(indexDocBuilder{}, indexDoc{}).(indexDocBuilder)
)

func (i indexDocBuilder) Doc() indexDoc {
	return builder.GetStruct(i).(indexDoc)
}

func (i indexDoc) ToModel() mongo.IndexModel {
	return mongo.IndexModel{Keys: i.Keys, Options: i.Options}
}

func (i indexDocBuilder) Key(field string, order OrderClause) indexDocBuilder {
	return builder.Append(i, "Keys", bson.E{Key: field, Value: order}).(indexDocBuilder)
}

func (i indexDocBuilder) Asc(field ...string) indexDocBuilder {
	for _, f := range field {
		i = i.Key(f, SortAsc)
	}
	return i
}

func (i indexDocBuilder) Desc(field ...string) indexDocBuilder {
	for _, f := range field {
		i = i.Key(f, SortDesc)
	}
	return i
}

func (i indexDocBuilder) Unique() indexDocBuilder {
	opt := options.Index()
	opt.SetUnique(true)
	return builder.Set(i, "Options", opt).(indexDocBuilder)
}

func (i indexDocBuilder) Name(name string) indexDocBuilder {
	idx := builder.GetStruct(i).(indexDoc)
	opt := idx.Options
	if opt == nil {
		opt = options.Index()
	}
	opt.SetName(name)
	return builder.Set(i, "Options", opt).(indexDocBuilder)
}
