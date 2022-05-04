package hamster

import (
	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// sortDoc is a MQL sort document
type sortDoc struct {
	Sorts bson.D
}

// sortDocBuilder is a builder for sortDoc
type sortDocBuilder builder.Builder

// OrderClause is a sort order
type OrderClause int32

const (
	// SortAsc is ascending sort order
	SortAsc OrderClause = 1
	// SortDesc is descending sort order
	SortDesc OrderClause = -1
)

var (
	// SortDocBuilder is a singleton builder for sortDoc
	SortDocBuilder = builder.Register(sortDocBuilder{}, sortDoc{}).(sortDocBuilder)
)

// Doc returns the sortDoc instance
func (s sortDocBuilder) Doc() sortDoc {
	return builder.GetStruct(s).(sortDoc)
}

// ToD convert the sortDoc to a bson.D
func (sd sortDoc) ToD() bson.D {
	return sd.Sorts
}

// ToM convert sortDoc to a bson.M
func (sd sortDoc) ToM() bson.M {
	return sd.ToD().Map()
}

// MarshalBSON marshals sortDoc to BSON
func (sd sortDoc) MarshalBSON() ([]byte, error) {
	return bson.Marshal(sd.ToD())
}

// UnmarshalBSON unmarshals BSON to sortDoc
func (sd *sortDoc) UnmarshalBSON(data []byte) error {
	return bson.Unmarshal(data, &sd.Sorts)
}

func (s sortDocBuilder) OrderBy(fieldName string, order OrderClause) sortDocBuilder {
	return builder.Append(s, "Sorts", primitive.E{Key: fieldName, Value: order}).(sortDocBuilder)
}

func (s sortDocBuilder) OrderAscBy(fieldName ...string) sortDocBuilder {
	for _, field := range fieldName {
		s = s.OrderBy(field, SortAsc)
	}
	return s
}

func (s sortDocBuilder) OrderDescBy(fieldName ...string) sortDocBuilder {
	for _, field := range fieldName {
		s = s.OrderBy(field, SortDesc)
	}
	return s
}
