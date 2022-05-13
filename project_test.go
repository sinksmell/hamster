package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestProjectDoc(t *testing.T) {

	t.Run("test-include", func(t *testing.T) {
		//         toBson(include('x')) == parse('{x : 1}')
		// toBson(include('x', 'y')) == parse('{x : 1, y : 1}')
		// toBson(include(['x', 'y'])) == parse('{x : 1, y : 1}')
		// toBson(include(['x', 'y', 'x'])) == parse('{y : 1, x : 1}')
		p := ProjectDocBuilder.Include("x").Doc().ToD()
		require.EqualValues(t, bson.D{{"x", int32(1)}}, p)

		p = ProjectDocBuilder.Include("x", "y").Doc().ToD()
		require.EqualValues(t, bson.D{{"x", int32(1)}, {"y", int32(1)}}, p)

		p = ProjectDocBuilder.Include("x", "y", "x").Doc().ToD()
		require.EqualValues(t, bson.D{{"y", int32(1)}, {"x", int32(1)}}, p)
	})

	t.Run("test-exclude", func(t *testing.T) {
		// toBson(exclude('x')) == parse('{x : 0}')
		// toBson(exclude('x', 'y')) == parse('{x : 0, y : 0}')
		// toBson(exclude(['x', 'y'])) == parse('{x : 0, y : 0}')

		p := ProjectDocBuilder.Exclude("x").Doc().ToD()
		require.EqualValues(t, bson.D{{"x", int32(0)}}, p)

		p = ProjectDocBuilder.Exclude("x", "y").Doc().ToD()
		require.EqualValues(t, bson.D{{"x", int32(0)}, {"y", int32(0)}}, p)
	})

	t.Run("test-exclude_id", func(t *testing.T) {
		// toBson(excludeId()) == parse('{_id : 0}')
		p := ProjectDocBuilder.ExcludeId().Doc().ToD()
		require.EqualValues(t, bson.D{{"_id", int32(0)}}, p)
	})

	t.Run("test-slice", func(t *testing.T) {
		//toBson(slice('x', 5)) == parse('{x : {$slice : 5}}')
		// toBson(slice('x', 5, 10)) == parse('{x : {$slice : [5, 10]}}')
		p := ProjectDocBuilder.Slice("x", 5).Doc().ToD()
		d1 := bson.D{bson.E{
			Key: "x",
			Value: bson.D{
				bson.E{
					Key:   "$slice",
					Value: int64(5),
				},
			},
		}}
		require.EqualValues(t, d1, p)

		p = ProjectDocBuilder.SliceWithSkip("x", 5, 10).Doc().ToD()
		d2 := bson.D{bson.E{
			Key: "x",
			Value: bson.D{
				bson.E{
					Key:   "$slice",
					Value: bson.A{int64(5), int64(10)},
				},
			},
		}}
		require.EqualValues(t, d2, p)
	})

	t.Run("test-elemMath", func(t *testing.T) {
		//  toBson(elemMatch('x', and(eq('y', 1), eq('z', 2)))) == parse('{x : {$elemMatch : {$and: [{y : 1}, {z : 2}]}}}')
		p := ProjectDocBuilder.ElemMatch("x", bson.D{
			bson.E{
				Key:   "$and",
				Value: bson.A{bson.D{bson.E{Key: "y", Value: int64(1)}, bson.E{Key: "z", Value: int64(2)}}},
			},
		}).Doc().ToD()
		d := bson.D{bson.E{
			Key: "x",
			Value: bson.D{
				bson.E{
					Key:   "$elemMatch",
					Value: bson.D{bson.E{Key: "$and", Value: bson.A{bson.D{bson.E{Key: "y", Value: int64(1)}, bson.E{Key: "z", Value: int64(2)}}}}}},
			},
		}}
		require.EqualValues(t, d, p)
	})

	t.Run("test-meta", func(t *testing.T) {
		// toBson(meta('x', 'textScore')) == parse('{x : {$meta : "textScore"}}')
		// toBson(meta('x', 'recordId')) == parse('{x : {$meta : "recordId"}}')
		// toBson(metaTextScore('x')) == parse('{x : {$meta : "textScore"}}')
		p := ProjectDocBuilder.Meta("x", "textScore").Doc().ToD()
		d := bson.D{bson.E{
			Key: "x",
			Value: bson.D{
				bson.E{
					Key:   "$meta",
					Value: "textScore",
				},
			},
		}}
		require.EqualValues(t, p, d)

		p = ProjectDocBuilder.MetaTextScore("x").Doc().ToD()
		require.EqualValues(t, p, d)
	})

	t.Run("test-field", func(t *testing.T) {
		p := ProjectDocBuilder.Field("x", bson.D{
			bson.E{
				Key:   "$meta",
				Value: "textScore",
			},
		}).Doc().ToD()
		d := bson.D{bson.E{
			Key: "x",
			Value: bson.D{
				bson.E{
					Key:   "$meta",
					Value: "textScore",
				},
			},
		}}
		require.EqualValues(t, p, d)
	})
}

func TestProjectMarshal(t *testing.T) {
	t.Run("test-marshal", func(t *testing.T) {
		p := ProjectDocBuilder.Include("x", "y", "x").Doc()
		d := bson.D{{"y", int32(1)}, {"x", int32(1)}}
		pd, err := bson.Marshal(p)
		require.NoError(t, err)
		dd, err := bson.Marshal(d)
		require.NoError(t, err)
		require.EqualValues(t, dd, pd)
	})

	t.Run("test-unmarshal", func(t *testing.T) {
		p := ProjectDocBuilder.Include("x", "y").Doc()
		pd, err := bson.Marshal(p)
		require.NoError(t, err)

		p2 := projectDoc{}
		bson.Unmarshal(pd, &p2)

		require.EqualValues(t, p.ToD(), p2.ToD())
	})
}
