package hamster

import (
	"github.com/lann/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// filterDoc is a MQL filter document
type filterDoc struct {
	Filters bson.D
}

// filterDocBuilder is a builder for filterDoc
type filterDocBuilder builder.Builder

var (
	// FilterDocBuilder is a builder for MQL filter document
	// https://www.mongodb.com/docs/drivers/java/sync/v4.6/fundamentals/builders/filters/#std-label-filters-builders
	// All filter operators category are:
	// Comparison
	// Logical
	// Arrays
	// Elements
	// Evaluation
	// Bitwise
	// Geospatial
	FilterDocBuilder = builder.Register(filterDocBuilder{}, filterDoc{}).(filterDocBuilder)
)

// ToD convert filterDoc to a bson.D filter document
func (f filterDoc) ToD() bson.D {
	return f.Filters
}

// ToM convert filterDoc to a bson.M filter document
func (f filterDoc) ToM() bson.M {
	return f.ToD().Map()
}

// MarshalBSON marshals filterDoc to BSON
func (f filterDoc) MarshalBSON() ([]byte, error) {
	return bson.Marshal(f.ToD())
}

// UnmarshalBSON unmarshals BSON to filterDoc
func (f *filterDoc) UnmarshalBSON(data []byte) error {
	return bson.Unmarshal(data, &f.Filters)
}

func (f filterDocBuilder) Eq(fieldName string, value interface{}) filterDocBuilder {
	return builder.Append(f, "Filters", bson.E{Key: fieldName, Value: value}).(filterDocBuilder)
}

func (f filterDocBuilder) Gt(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$gt", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GtE(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$gte", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Lt(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$lt", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) LtE(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$lte", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Ne(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$ne", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) In(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$in", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Nin(fieldName string, value interface{}) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$nin", Value: value}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Empty() filterDocBuilder {
	return builder.Set(f, "Filters", bson.D{}).(filterDocBuilder)
}

func (f filterDocBuilder) Doc() filterDoc {
	return builder.GetStruct(f).(filterDoc)
}

func (f filterDocBuilder) And(filters ...filterDoc) filterDocBuilder {
	v := make([]bson.D, 0, len(filters))
	for _, filter := range filters {
		v = append(v, filter.ToD())
	}
	e := bson.E{Key: "$and", Value: v}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Or(filters ...filterDoc) filterDocBuilder {
	v := make([]bson.D, 0, len(filters))
	for _, filter := range filters {
		v = append(v, filter.ToD())
	}
	e := bson.E{Key: "$or", Value: v}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Not(filter filterDoc) filterDocBuilder {
	e := bson.E{Key: "$not", Value: filter.ToD()}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Nor(filters ...filterDoc) filterDocBuilder {
	v := make([]bson.D, 0, len(filters))
	for _, filter := range filters {
		v = append(v, filter.ToD())
	}
	e := bson.E{Key: "$nor", Value: v}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) All(fieldName string, values []interface{}) filterDocBuilder {
	array := make(bson.A, len(values))
	copy(array, values)
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$all", Value: array}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) ElemMatch(fieldName string, filter bson.D) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$elemMatch", Value: filter}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Size(fieldName string, size int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$size", Value: size}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Exists(fieldName string) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$exists", Value: true}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Type(fieldName string, bsonType ...string) filterDocBuilder {
	if len(bsonType) == 1 {
		return f.typeOne(fieldName, bsonType[0])
	}
	return f.typeMany(fieldName, bsonType)
}

func (f filterDocBuilder) typeOne(fieldName string, bsonType string) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$type", Value: bsonType}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) typeMany(fieldName string, bsonTypes []string) filterDocBuilder {
	arr := make(bson.A, 0, len(bsonTypes))
	for _, tp := range bsonTypes {
		arr = append(arr, tp)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$type", Value: arr}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Mod(fieldName string, divisor, remainder int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$mod", Value: bson.A{divisor, remainder}}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Regex(fieldName string, pattern string, options string) filterDocBuilder {
	if options == "" {
		e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$regex", Value: pattern}}}
		return builder.Append(f, "Filters", e).(filterDocBuilder)
	}

	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$regex", Value: pattern},
		bson.E{Key: "$options", Value: options}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

type FilterDocTextSearchOptions struct {
	Language           *string `bson:"$language"`
	CaseSensitive      *bool   `bson:"$caseSensitive"`
	DiacriticSensitive *bool   `bson:"$diacriticSensitive"`
}

func (f filterDocBuilder) Text(search string, opt *FilterDocTextSearchOptions) filterDocBuilder {
	d := bson.D{bson.E{Key: "$search", Value: search}}
	if opt != nil {
		if opt.Language != nil {
			d = append(d, bson.E{Key: "$language", Value: *opt.Language})
		}
		if opt.CaseSensitive != nil {
			d = append(d, bson.E{Key: "$caseSensitive", Value: *opt.CaseSensitive})
		}
		if opt.DiacriticSensitive != nil {
			d = append(d, bson.E{Key: "$diacriticSensitive", Value: *opt.DiacriticSensitive})
		}
	}
	e := bson.E{Key: "$text", Value: d}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Where(js primitive.JavaScript) filterDocBuilder {
	e := bson.E{Key: "$where", Value: js}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAllClearWithMask(fieldName string, bitmask int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAllClear", Value: bitmask}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAllClearWithBitPosition(fieldName string, position []int64) filterDocBuilder {
	arr := make(bson.A, 0, len(position))
	for _, p := range position {
		arr = append(arr, p)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAllClear", Value: arr}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAllSetWithMask(fieldName string, bitmask int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAllSet", Value: bitmask}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAllSetWithBitPosition(fieldName string, position []int64) filterDocBuilder {
	arr := make(bson.A, 0, len(position))
	for _, p := range position {
		arr = append(arr, p)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAllSet", Value: arr}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAnyClearWithMask(fieldName string, bitmask int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAnyClear", Value: bitmask}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAnyClearWithBitPosition(fieldName string, position []int64) filterDocBuilder {
	arr := make(bson.A, 0, len(position))
	for _, p := range position {
		arr = append(arr, p)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAnyClear", Value: arr}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAnySetWithMask(fieldName string, bitmask int64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAnySet", Value: bitmask}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) BitsAnySetWithBitPosition(fieldName string, position []int64) filterDocBuilder {
	arr := make(bson.A, 0, len(position))
	for _, p := range position {
		arr = append(arr, p)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$bitsAnySet", Value: arr}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoWithin(fieldName string, geoWithinDoc bson.D) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoWithin", Value: geoWithinDoc}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoWithinBox(fieldName string, lowerLeftX, lowerLeftY, upperRightX, upperRightY float64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoWithin",
		Value: bson.D{bson.E{Key: "$box", Value: bson.A{
			bson.A{lowerLeftX, lowerLeftY},
			bson.A{upperRightX, upperRightY},
		}}}}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoWithinPolygon(fieldName string, points []bson.A) filterDocBuilder {
	arr := make(bson.A, 0, len(points))
	for _, p := range points {
		arr = append(arr, p)
	}
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoWithin",
		Value: bson.D{bson.E{Key: "$polygon", Value: arr}}}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoWithCenter(fieldName string, x, y, radius float64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoWithin",
		Value: bson.D{bson.E{Key: "$center", Value: bson.A{bson.A{x, y}, radius}}}}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoWithCenterSphere(fieldName string, x, y, radius float64) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoWithin",
		Value: bson.D{bson.E{Key: "$centerSphere", Value: bson.A{bson.A{x, y}, radius}}}}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) GeoIntersects(fieldName string, geoIntersectsDoc bson.D) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$geoIntersects", Value: geoIntersectsDoc}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}

func (f filterDocBuilder) Near(fieldName string, nearDoc bson.D) filterDocBuilder {
	e := bson.E{Key: fieldName, Value: bson.D{bson.E{Key: "$near", Value: nearDoc}}}
	return builder.Append(f, "Filters", e).(filterDocBuilder)
}
