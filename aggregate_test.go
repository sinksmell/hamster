package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAggregateDocBuilder(t *testing.T) {
	doc := AggregateDocBuilder.
		Match(bson.D{{Key: "year", Value: bson.D{{Key: "$gte", Value: 2000}}}}).
		Sort(bson.D{{Key: "rating", Value: -1}}).
		Limit(10).
		Doc()

	require.Len(t, doc.ToA(), 3)
	require.EqualValues(t, bson.D{{Key: "$limit", Value: int64(10)}}, doc.ToA()[2])
}
