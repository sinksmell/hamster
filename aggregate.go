package hamster

import (
	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
)

// aggregateDoc is a MQL aggregate pipeline
type aggregateDoc struct {
	Pipeline bson.A
}

// aggregateDocBuilder is a builder for aggregateDoc
type aggregateDocBuilder builder.Builder

var (
	// AggregateDocBuilder is a singleton builder for aggregateDoc
	AggregateDocBuilder = builder.Register(aggregateDocBuilder{}, aggregateDoc{}).(aggregateDocBuilder)
)

func (a aggregateDocBuilder) Doc() aggregateDoc {
	return builder.GetStruct(a).(aggregateDoc)
}

func (a aggregateDoc) ToA() bson.A {
	return a.Pipeline
}

func (a aggregateDoc) MarshalBSON() ([]byte, error) {
	return bson.Marshal(bson.D{{Key: "pipeline", Value: a.Pipeline}})
}

func (a *aggregateDoc) UnmarshalBSON(data []byte) error {
	var raw struct {
		Pipeline bson.A `bson:"pipeline"`
	}
	if err := bson.Unmarshal(data, &raw); err != nil {
		return err
	}
	a.Pipeline = raw.Pipeline
	return nil
}

func (a aggregateDocBuilder) stage(stage string, value interface{}) aggregateDocBuilder {
	return builder.Append(a, "Pipeline", bson.D{{Key: stage, Value: value}}).(aggregateDocBuilder)
}

func (a aggregateDocBuilder) Match(filter bson.D) aggregateDocBuilder {
	return a.stage("$match", filter)
}

func (a aggregateDocBuilder) Project(project bson.D) aggregateDocBuilder {
	return a.stage("$project", project)
}

func (a aggregateDocBuilder) Group(group bson.D) aggregateDocBuilder {
	return a.stage("$group", group)
}

func (a aggregateDocBuilder) Sort(sort bson.D) aggregateDocBuilder {
	return a.stage("$sort", sort)
}

func (a aggregateDocBuilder) Limit(limit int64) aggregateDocBuilder {
	return a.stage("$limit", limit)
}

func (a aggregateDocBuilder) Skip(skip int64) aggregateDocBuilder {
	return a.stage("$skip", skip)
}

func (a aggregateDocBuilder) Unwind(path string) aggregateDocBuilder {
	return a.stage("$unwind", path)
}

func (a aggregateDocBuilder) AddStage(stage bson.D) aggregateDocBuilder {
	return builder.Append(a, "Pipeline", stage).(aggregateDocBuilder)
}
