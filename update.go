package hamster

import (
	"time"

	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
)

// updateDoc is a MQL update document
type updateDoc struct {
	Updates bson.D
}

// updateDocBuilder is a builder for updateDoc
type updateDocBuilder builder.Builder

var (
	// UpdateDocBuilder is a singleton builder for updateDoc
	UpdateDocBuilder = builder.Register(updateDocBuilder{}, updateDoc{}).(updateDocBuilder)
)

// Doc returns the updateDoc instance
func (u updateDocBuilder) Doc() updateDoc {
	return builder.GetStruct(u).(updateDoc)
}

// ToD convert updateDoc to bson.D
func (u updateDoc) ToD() bson.D {
	return u.Updates
}

// ToM convert updateDoc to bson.M
func (u updateDoc) ToM() bson.M {
	return u.ToD().Map()
}

// MarshalBSON marshals updateDoc to BSON
func (u updateDoc) MarshalBSON() ([]byte, error) {
	return bson.Marshal(u.ToD())
}

// UnmarshalBSON unmarshals BSON to updateDoc
func (u *updateDoc) UnmarshalBSON(data []byte) error {
	return bson.Unmarshal(data, &u.Updates)
}

func (u updateDocBuilder) appendOperator(operator, field string, value interface{}) updateDocBuilder {
	e := bson.E{Key: operator, Value: bson.D{{Key: field, Value: value}}}
	return builder.Append(u, "Updates", e).(updateDocBuilder)
}

func (u updateDocBuilder) Set(field string, value interface{}) updateDocBuilder {
	return u.appendOperator("$set", field, value)
}

func (u updateDocBuilder) Unset(field string) updateDocBuilder {
	return u.appendOperator("$unset", field, "")
}

func (u updateDocBuilder) Inc(field string, value interface{}) updateDocBuilder {
	return u.appendOperator("$inc", field, value)
}

func (u updateDocBuilder) Mul(field string, value interface{}) updateDocBuilder {
	return u.appendOperator("$mul", field, value)
}

func (u updateDocBuilder) Min(field string, value interface{}) updateDocBuilder {
	return u.appendOperator("$min", field, value)
}

func (u updateDocBuilder) Max(field string, value interface{}) updateDocBuilder {
	return u.appendOperator("$max", field, value)
}

func (u updateDocBuilder) Rename(oldField, newField string) updateDocBuilder {
	return u.appendOperator("$rename", oldField, newField)
}

func (u updateDocBuilder) CurrentDate(field string) updateDocBuilder {
	return u.appendOperator("$currentDate", field, true)
}

func (u updateDocBuilder) CurrentTimestamp(field string) updateDocBuilder {
	return u.appendOperator("$currentDate", field, bson.D{{Key: "$type", Value: "timestamp"}})
}

func (u updateDocBuilder) SetDate(field string, value time.Time) updateDocBuilder {
	return u.Set(field, value)
}
