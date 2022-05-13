package hamster

import (
	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
)

// projectDoc is a MQL project document
type projectDoc struct {
	Projects bson.D
}

// projectDocBuilder is a builder for projectDoc
type projectDocBuilder builder.Builder

var (
	// projectDocBuilder is a builder for MQL project document
	ProjectDocBuilder = builder.Register(projectDocBuilder{}, projectDoc{}).(projectDocBuilder)
)

// ToD convert projectDoc to a bson.D project document
func (p projectDoc) ToD() bson.D {
	book := map[string]interface{}{}
	for _, e := range p.Projects {
		book[e.Key] = e.Value
	}

	d := bson.D{}
	included := map[string]bool{}
	for i := len(p.Projects) - 1; i >= 0; i-- {
		if !included[p.Projects[i].Key] {
			d = append(d, bson.E{
				Key:   p.Projects[i].Key,
				Value: book[p.Projects[i].Key],
			})
			included[p.Projects[i].Key] = true
		}
	}

	// reverse d
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}

	return d
}

// ToM convert projectDoc to a bson.M project document
func (p projectDoc) ToM() bson.M {
	return p.ToD().Map()
}

// MarshalBSON marshals projectDoc to BSON
func (p projectDoc) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.ToD())
}

// UnmarshalBSON unmarshals BSON to projectDoc
func (p *projectDoc) UnmarshalBSON(data []byte) error {
	return bson.Unmarshal(data, &p.Projects)
}

func (p projectDocBuilder) Doc() projectDoc {
	return builder.GetStruct(p).(projectDoc)
}

// Creates a projection that excludes the _id field.  This suppresses the automatic inclusion of _id that is the default, even when
func (p projectDocBuilder) ExcludeId() projectDocBuilder {
	return builder.Append(p, "Projects", bson.E{"_id", int32(0)}).(projectDocBuilder)
}

func (p projectDocBuilder) exclude(field string) projectDocBuilder {
	return builder.Append(p, "Projects", bson.E{field, int32(0)}).(projectDocBuilder)
}

// Creates a projection that excludes all of the given fields.
func (p projectDocBuilder) Exclude(fields ...string) projectDocBuilder {
	for _, field := range fields {
		p = p.exclude(field)
	}
	return p
}

// Creates a projection that includes all of the given fields.
func (p projectDocBuilder) Include(fields ...string) projectDocBuilder {
	for _, field := range fields {
		p = p.include(field)
	}
	return p
}

func (p projectDocBuilder) include(filed string) projectDocBuilder {
	return builder.Append(p, "Projects", bson.E{filed, int32(1)}).(projectDocBuilder)
}

// Creates a projection that includes for the given field only the first element of the array value of that field that matches the given
func (p projectDocBuilder) ElemMatch(field string, filter bson.D) projectDocBuilder {
	e := bson.E{Key: field, Value: bson.D{{"$elemMatch", filter}}}
	return builder.Append(p, "Projects", e).(projectDocBuilder)
}

// Creates a projection to the given field name of a slice of the array value of that field.
func (p projectDocBuilder) Slice(field string, limit int64) projectDocBuilder {
	e := bson.E{Key: field, Value: bson.D{{"$slice", limit}}}
	return builder.Append(p, "Projects", e).(projectDocBuilder)
}

// Creates a projection to the given field name of a slice of the array value of that field.
func (p projectDocBuilder) SliceWithSkip(field string, skip, limit int64) projectDocBuilder {
	e := bson.E{Key: field, Value: bson.D{{"$slice", bson.A{skip, limit}}}}
	return builder.Append(p, "Projects", e).(projectDocBuilder)
}

func (p projectDocBuilder) Field(field string, value bson.D) projectDocBuilder {
	return builder.Append(p, "Projects", bson.E{field, value}).(projectDocBuilder)
}

// Creates a $meta projection to the given field name for the given meta field name.
func (p projectDocBuilder) Meta(field, metaFieldName string) projectDocBuilder {
	return builder.Append(p, "Projects", bson.E{field, bson.D{{"$meta", metaFieldName}}}).(projectDocBuilder)
}

// Creates a projection to the given field name of the textScore, for use with text queries.
func (p projectDocBuilder) MetaTextScore(field string) projectDocBuilder {
	return p.Meta(field, "textScore")
}
