package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBuildFilterDoc(t *testing.T) {
	doc := FilterDocBuilder.Eq("a", 1).Eq("b", 2).Doc()
	require.Equal(t, len(doc.ToD()), 2)
	require.Contains(t, doc.ToD(), bson.E{Key: "a", Value: 1})
	require.Contains(t, doc.ToD(), bson.E{Key: "b", Value: 2})
}

func TestFilterDocComparison(t *testing.T) {
	doc := FilterDocBuilder.Eq("a", 1).
		Gt("b", 2).
		GtE("c", 3).
		Lt("d", 4).
		LtE("e", 5).
		Ne("f", 6).
		In("g", []int{7, 8}).
		Nin("h", []int{9, 10}).
		Doc()

	std := bson.D{
		bson.E{Key: "a", Value: 1},
		bson.E{Key: "b", Value: bson.D{{"$gt", 2}}},
		bson.E{Key: "c", Value: bson.D{{"$gte", 3}}},
		bson.E{Key: "d", Value: bson.D{{"$lt", 4}}},
		bson.E{Key: "e", Value: bson.D{{"$lte", 5}}},
		bson.E{Key: "f", Value: bson.D{{"$ne", 6}}},
		bson.E{Key: "g", Value: bson.D{{"$in", []int{7, 8}}}},
		bson.E{Key: "h", Value: bson.D{{"$nin", []int{9, 10}}}},
	}

	require.ElementsMatch(t, doc.ToD(), std)
}

func TestFilterDocMarshal(t *testing.T) {
	doc := FilterDocBuilder.Eq("a", 1).
		Gt("b", 2).
		GtE("c", 3).
		Lt("d", 4).
		LtE("e", 5).
		Ne("f", 6).
		In("g", []int{7, 8}).
		Nin("h", []int{9, 10}).
		Doc()

	std := bson.D{
		bson.E{Key: "a", Value: 1},
		bson.E{Key: "b", Value: bson.D{{"$gt", 2}}},
		bson.E{Key: "c", Value: bson.D{{"$gte", 3}}},
		bson.E{Key: "d", Value: bson.D{{"$lt", 4}}},
		bson.E{Key: "e", Value: bson.D{{"$lte", 5}}},
		bson.E{Key: "f", Value: bson.D{{"$ne", 6}}},
		bson.E{Key: "g", Value: bson.D{{"$in", []int{7, 8}}}},
		bson.E{Key: "h", Value: bson.D{{"$nin", []int{9, 10}}}},
	}

	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := bson.Marshal(std)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, data, data2)
}

func TestFilterDocUnmarshal(t *testing.T) {
	doc := FilterDocBuilder.Eq("a", int32(1)).
		Gt("b", int32(2)).
		GtE("c", int32(3)).
		Lt("d", int32(4)).
		LtE("e", int32(5)).
		Ne("f", int32(6)).
		In("g", bson.A{int32(7), int32(8)}).
		Nin("h", bson.A{int32(9), int32(10)}).
		Doc()

	data, err := bson.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	doc2 := filterDoc{}
	bson.Unmarshal(data, &doc2)

	require.EqualValues(t, doc.ToD(), doc2.ToD())
}

func TestEmptyFilterDoc(t *testing.T) {
	doc := FilterDocBuilder.Empty().Doc()
	require.Equal(t, len(doc.ToD()), 0)
	require.Equal(t, len(doc.ToM()), 0)
	require.Equal(t, doc.ToD(), bson.D{})
}

func TestFilterDocLogic(t *testing.T) {
	// $and
	andDoc := FilterDocBuilder.And(
		FilterDocBuilder.Eq("a", 1).Doc(),
		FilterDocBuilder.Gt("b", 2).Doc(),
	).Doc()

	andBson := bson.D{
		{
			Key: "$and",
			Value: []bson.D{{{Key: "a", Value: 1}},
				{{Key: "b", Value: bson.D{{Key: "$gt", Value: 2}}}}},
		},
	}

	require.ElementsMatch(t, andDoc.ToD(), andBson)

	// $or
	orDoc := FilterDocBuilder.Or(
		FilterDocBuilder.Eq("a", 1).Doc(),
		FilterDocBuilder.Gt("b", 2).Doc(),
	).Doc()

	orBson := bson.D{
		{
			Key: "$or",
			Value: []bson.D{{{Key: "a", Value: 1}},
				{{Key: "b", Value: bson.D{{Key: "$gt", Value: 2}}}}},
		},
	}

	require.ElementsMatch(t, orDoc.ToD(), orBson)

	// $not
	notDoc := FilterDocBuilder.Not(
		FilterDocBuilder.Eq("a", 1).Doc(),
	).Doc()

	notBson := bson.D{
		{
			Key:   "$not",
			Value: bson.D{{Key: "a", Value: 1}},
		},
	}

	require.ElementsMatch(t, notDoc.ToD(), notBson)

	// $nor
	norDoc := FilterDocBuilder.Nor(
		FilterDocBuilder.Eq("a", 1).Doc(),
		FilterDocBuilder.Gt("b", 2).Doc(),
	).Doc()

	norBson := bson.D{
		{
			Key: "$nor",
			Value: []bson.D{{{Key: "a", Value: 1}},
				{{Key: "b", Value: bson.D{{Key: "$gt", Value: 2}}}}},
		},
	}

	require.ElementsMatch(t, norDoc.ToD(), norBson)
}

func TestFilterDocArrays(t *testing.T) {
	// $all
	// { tags: { $all: [ "ssl" , "security" ] } }
	tags := []interface{}{"ssl", "security"}
	allDoc := FilterDocBuilder.All(
		"tags", tags,
	).Doc()
	allBson := bson.D{
		{
			Key: "tags",
			Value: bson.D{{
				Key:   "$all",
				Value: bson.A{"ssl", "security"},
			}},
		},
	}

	require.ElementsMatch(t, allDoc.ToD(), allBson)

	// $elemMatch
	//  { results: { $elemMatch: { $gte: 80, $lt: 85 } } }
	elemMatchDoc := FilterDocBuilder.ElemMatch(
		"results",
		bson.D{{
			Key:   "$gte",
			Value: 80,
		},
			{
				Key:   "$lt",
				Value: 85,
			},
		},
	).Doc()
	elemMatchBson := bson.D{
		{
			Key: "results",
			Value: bson.D{{
				Key:   "$elemMatch",
				Value: bson.D{{Key: "$gte", Value: 80}, {Key: "$lt", Value: 85}},
			}},
		},
	}

	require.ElementsMatch(t, elemMatchDoc.ToD(), elemMatchBson)

	// $size
	// { tags: { $size: 3 } }
	sizeDoc := FilterDocBuilder.Size("tags", 3).Doc()
	sizeBson := bson.D{{
		Key:   "tags",
		Value: bson.D{{Key: "$size", Value: int64(3)}},
	}}
	require.ElementsMatch(t, sizeDoc.ToD(), sizeBson)
}

func TestFilterDocElements(t *testing.T) {
	// $exists
	// { a: { $exists: true }
	existsDoc := FilterDocBuilder.Exists("a").Doc()
	existsBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$exists", Value: true}},
	}}
	require.ElementsMatch(t, existsDoc.ToD(), existsBson)

	// $type
	// { "zipCode" : { $type : "string" } }
	//  {"data" : { $type : ["array","string"] } }
	typeOneDoc := FilterDocBuilder.Type("zipCode", "string").Doc()
	typeOneBson := bson.D{{
		Key:   "zipCode",
		Value: bson.D{{Key: "$type", Value: "string"}},
	}}
	require.ElementsMatch(t, typeOneDoc.ToD(), typeOneBson)

	typeTwoDoc := FilterDocBuilder.Type("data", "array", "string").Doc()
	typeTwoBson := bson.D{{
		Key:   "data",
		Value: bson.D{{Key: "$type", Value: bson.A{"array", "string"}}},
	}}
	require.ElementsMatch(t, typeTwoDoc.ToD(), typeTwoBson)
}

func TestFilterDocEvaluation(t *testing.T) {
	// $mod
	// { qty: { $mod: [ 4, 0 ] } }
	modDoc := FilterDocBuilder.Mod("qty", 4, 0).Doc()
	modBson := bson.D{
		{
			Key:   "qty",
			Value: bson.D{{Key: "$mod", Value: bson.A{int64(4), int64(0)}}},
		},
	}
	require.ElementsMatch(t, modDoc.ToD(), modBson)

	// $regex
	// { name: { $regex: 'acme.*corp', $options: "si" } }
	regexDoc := FilterDocBuilder.Regex("name", "acme.*corp", "si").Doc()
	regexBson := bson.D{
		{
			Key:   "name",
			Value: bson.D{{Key: "$regex", Value: "acme.*corp"}, {Key: "$options", Value: "si"}},
		},
	}
	require.ElementsMatch(t, regexDoc.ToD(), regexBson)

	// $text
	// { $text: { $search: "coffee" } }
	// { $text: { $search: "leches -cafés", $diacriticSensitive: true } }
	textNoOptionDoc := FilterDocBuilder.Text("coffee", nil).Doc()
	textNoOptionBson := bson.D{{
		Key:   "$text",
		Value: bson.D{{Key: "$search", Value: "coffee"}},
	}}
	require.ElementsMatch(t, textNoOptionDoc.ToD(), textNoOptionBson)

	opt := &FilterDocTextSearchOptions{}
	diacriticSensitive := true
	opt.DiacriticSensitive = &diacriticSensitive
	textWithOption := FilterDocBuilder.Text("leches -cafés", opt).Doc()
	textWithOptionBson := bson.D{{
		Key:   "$text",
		Value: bson.D{{Key: "$search", Value: "leches -cafés"}, {Key: "$diacriticSensitive", Value: true}},
	}}
	require.ElementsMatch(t, textWithOption.ToD(), textWithOptionBson)

	// $where
	// { $where: function() {
	//    return (hex_md5(this.name) == "9b53e667f30cd329dca1ec9e6a83e994")
	// } }
	whereDoc := FilterDocBuilder.Where(`
	function() {
	    return (hex_md5(this.name) == "9b53e667f30cd329dca1ec9e6a83e994")
	 }`).Doc()
	whereBson := bson.D{{
		Key: "$where",
		Value: primitive.JavaScript(`
	function() {
	    return (hex_md5(this.name) == "9b53e667f30cd329dca1ec9e6a83e994")
	 }`),
	}}
	require.ElementsMatch(t, whereDoc.ToD(), whereBson)
}

func TestFilterDocBitwise(t *testing.T) {
	// $bitsAllClear
	//  { a: { $bitsAllClear: [ 1, 5 ] } }
	bitsAllClearPositionDoc := FilterDocBuilder.BitsAllClearWithBitPosition("a", []int64{1, 5}).Doc()
	bitsAllClearPositionBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAllClear", Value: bson.A{int64(1), int64(5)}}},
	}}
	require.ElementsMatch(t, bitsAllClearPositionDoc.ToD(), bitsAllClearPositionBson)

	// { a: { $bitsAllClear: 35 } }
	bitsAllClearMaskDoc := FilterDocBuilder.BitsAllClearWithMask("a", 35).Doc()
	bitsAllClearMaskBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAllClear", Value: int64(35)}},
	}}
	require.ElementsMatch(t, bitsAllClearMaskDoc.ToD(), bitsAllClearMaskBson)

	// $bitsAnyClear
	//  { a: { $bitsAnyClear: [ 1, 5 ] } }
	bitsAnyClearPositionDoc := FilterDocBuilder.BitsAnyClearWithBitPosition("a", []int64{1, 5}).Doc()
	bitsAnyClearPositionBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAnyClear", Value: bson.A{int64(1), int64(5)}}},
	}}
	require.ElementsMatch(t, bitsAnyClearPositionDoc.ToD(), bitsAnyClearPositionBson)

	// { a: { $bitsAnyClear: 35 } }
	bitsAnyClearMaskDoc := FilterDocBuilder.BitsAnyClearWithMask("a", 35).Doc()
	bitsAnyClearMaskBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAnyClear", Value: int64(35)}},
	}}
	require.ElementsMatch(t, bitsAnyClearMaskDoc.ToD(), bitsAnyClearMaskBson)

	// $bitsAllSet
	//  { a: { $bitsAllSet: [ 1, 5 ] } }
	bitsAllSetPositionDoc := FilterDocBuilder.BitsAllSetWithBitPosition("a", []int64{1, 5}).Doc()
	bitsAllSetPositionBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAllSet", Value: bson.A{int64(1), int64(5)}}},
	}}
	require.ElementsMatch(t, bitsAllSetPositionDoc.ToD(), bitsAllSetPositionBson)

	// { a: { $bitsAllSet: 50 } }
	bitsAllSetMaskDoc := FilterDocBuilder.BitsAllSetWithMask("a", 50).Doc()
	bitsAllSetMaskBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAllSet", Value: int64(50)}},
	}}
	require.ElementsMatch(t, bitsAllSetMaskDoc.ToD(), bitsAllSetMaskBson)

	// $bitsAnySet
	// { a: { $bitsAnySet: [ 1, 5 ] } }
	bitsAnySetPositionDoc := FilterDocBuilder.BitsAnySetWithBitPosition("a", []int64{1, 5}).Doc()
	bitsAnySetPositionBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAnySet", Value: bson.A{int64(1), int64(5)}}},
	}}
	require.ElementsMatch(t, bitsAnySetPositionDoc.ToD(), bitsAnySetPositionBson)

	//  { a: { $bitsAnySet: 35 } }
	bitsAnySetMaskDoc := FilterDocBuilder.BitsAnySetWithMask("a", 35).Doc()
	bitsAnySetMaskBson := bson.D{{
		Key:   "a",
		Value: bson.D{{Key: "$bitsAnySet", Value: int64(35)}},
	}}
	require.ElementsMatch(t, bitsAnySetMaskDoc.ToD(), bitsAnySetMaskBson)
}

func TestFilterDocGeo(t *testing.T) {
	// $geoIntersects
	// { loc: { $geoIntersects: { $geometry: { type: "Point", coordinates: [ -73.97, 40.77 ] } } } }
	d := bson.D{{Key: "$geometry", Value: bson.D{{Key: "type", Value: "Point"},
		{Key: "coordinates", Value: bson.A{-73.97, 40.77}}}},
		{Key: "minDistance", Value: 0},
		{Key: "spherical", Value: true}}
	geoIntersectsDoc := FilterDocBuilder.GeoIntersects("loc", d).Doc()
	geoIntersectsBson := bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoIntersects", Value: d}},
	}}
	require.ElementsMatch(t, geoIntersectsDoc.ToD(), geoIntersectsBson)

	// $geoWithin
	// { loc: { $geoWithin: { $geometry: { type: "Polygon", coordinates: [ [ [ -73.9558, 40.8124 ], [ -73.9486, 40.7887 ], [ -73.9391, 40.7932 ], [ -73.9558, 40.8124 ] ] ] } } } }
	d = bson.D{{Key: "$geometry", Value: bson.D{{Key: "type", Value: "Polygon"},
		{Key: "coordinates", Value: bson.A{bson.A{bson.A{-73.9558, 40.8124},
			bson.A{-73.9486, 40.7887}, bson.A{-73.9391, 40.7932}, bson.A{-73.9558, 40.8124}}}}}}}
	geoWhithinDoc := FilterDocBuilder.GeoWithin("loc", d).Doc()
	geoWhithinBson := bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoWithin", Value: d}},
	}}
	require.ElementsMatch(t, geoWhithinDoc.ToD(), geoWhithinBson)

	// $near
	// { loc: { $near: { $geometry: { type: "Point", coordinates: [ -73.97, 40.77 ] } } } }
	d = bson.D{{Key: "$geometry", Value: bson.D{{Key: "type", Value: "Point"}}}}
	nearDoc := FilterDocBuilder.Near("loc", d).Doc()
	nearBson := bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$near", Value: d}},
	}}
	require.ElementsMatch(t, nearDoc.ToD(), nearBson)

	// $geoWithin
	// { loc: { $geoWithin: { $center: [ [ -73.97, 40.77 ], 5 ] } } }
	d = bson.D{{Key: "$center", Value: bson.A{bson.A{-73.97, 40.77}, float64(5)}}}
	geoWithinDoc := FilterDocBuilder.GeoWithCenter("loc", -73.97, 40.77, 5).Doc()
	geoWithinBson := bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoWithin", Value: d}},
	}}
	require.ElementsMatch(t, geoWithinDoc.ToD(), geoWithinBson)

	// $geoWithin
	// { loc: { $geoWithin: { $centerSphere: [ [ -73.97, 40.77 ], 5 ] } } }
	d = bson.D{{Key: "$centerSphere", Value: bson.A{bson.A{-73.97, 40.77}, float64(5)}}}
	geoWithinDoc = FilterDocBuilder.GeoWithCenterSphere("loc", -73.97, 40.77, 5).Doc()
	geoWithinBson = bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoWithin", Value: d}},
	}}
	require.ElementsMatch(t, geoWithinDoc.ToD(), geoWithinBson)

	// $geoWithin
	// { loc: { $geoWithin: { $box: [ [ -73.97, 40.77 ], [ -73.98, 40.78 ] ] } } }
	d = bson.D{{Key: "$box", Value: bson.A{bson.A{-73.97, 40.77}, bson.A{-73.98, 40.78}}}}
	geoWithinDoc = FilterDocBuilder.GeoWithinBox("loc", -73.97, 40.77, -73.98, 40.78).Doc()
	geoWithinBson = bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoWithin", Value: d}},
	}}
	require.ElementsMatch(t, geoWithinDoc.ToD(), geoWithinBson)

	// $geoWithin
	// { loc: { $geoWithin: { $polygon: [ [ -73.97, 40.77 ], [ -73.98, 40.78 ], [ -73.96, 40.78 ], [ -73.97, 40.77 ] ] } } }
	d = bson.D{{Key: "$polygon", Value: bson.A{bson.A{-73.97, 40.77}, bson.A{-73.98, 40.78}, bson.A{-73.96, 40.78}, bson.A{-73.97, 40.77}}}}
	geoWithinDoc = FilterDocBuilder.GeoWithinPolygon("loc",
		[]bson.A{bson.A{-73.97, 40.77}, bson.A{-73.98, 40.78}, bson.A{-73.96, 40.78}, bson.A{-73.97, 40.77}}).Doc()
	geoWithinBson = bson.D{{
		Key:   "loc",
		Value: bson.D{{Key: "$geoWithin", Value: d}},
	}}
	require.ElementsMatch(t, geoWithinDoc.ToD(), geoWithinBson)
}
