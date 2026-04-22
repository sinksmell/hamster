package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUpdateDocBuilder(t *testing.T) {
	doc := UpdateDocBuilder.
		Set("name", "hamster").
		Inc("count", 1).
		Unset("legacy").
		Rename("old", "new").
		CurrentDate("updated_at").
		Doc()

	require.Len(t, doc.ToD(), 5)
	require.Contains(t, doc.ToD(), bson.E{Key: "$set", Value: bson.D{{Key: "name", Value: "hamster"}}})
	require.Contains(t, doc.ToD(), bson.E{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}})
}

func TestUpdateDocMarshal(t *testing.T) {
	doc := UpdateDocBuilder.Set("x", 1).Max("y", 2).Doc()
	data, err := bson.Marshal(doc)
	require.NoError(t, err)
	std, err := bson.Marshal(bson.D{{Key: "$set", Value: bson.D{{Key: "x", Value: 1}}}, {Key: "$max", Value: bson.D{{Key: "y", Value: 2}}}})
	require.NoError(t, err)
	require.EqualValues(t, std, data)
}
