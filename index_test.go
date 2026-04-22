package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestIndexDocBuilder(t *testing.T) {
	idx := IndexDocBuilder.Asc("email").Desc("created_at").Unique().Name("email_created_idx").Doc()
	model := idx.ToModel()
	keys := model.Keys.(bson.D)
	require.EqualValues(t, bson.D{{Key: "email", Value: SortAsc}, {Key: "created_at", Value: SortDesc}}, keys)
	require.NotNil(t, model.Options)
	require.NotNil(t, model.Options.Name)
	require.Equal(t, "email_created_idx", *model.Options.Name)
	require.NotNil(t, model.Options.Unique)
	require.True(t, *model.Options.Unique)
}
