package hamster

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildSortDoc(t *testing.T) {
	sortDoc := SortDocBuilder.OrderBy("a", SortAsc).
		OrderBy("b", SortDesc).
		OrderAscBy("c", "d").
		OrderDescBy("e", "f").
		Doc()

	require.NotEmpty(t, sortDoc.ToD())
	require.Equal(t, len(sortDoc.ToD()), 6)
	require.Equal(t, len(sortDoc.ToD()), len(sortDoc.ToM()))
}

func TestSortDocMarshalBson(t *testing.T) {
	sortDoc := SortDocBuilder.OrderBy("a", SortAsc).
		OrderBy("b", SortDesc).
		OrderAscBy("c", "d").
		OrderDescBy("e", "f").
		Doc()

	data, err := bson.Marshal(sortDoc)
	if err != nil {
		t.Error(err)
	}

	s2 := bson.D{}
	s2 = append(s2, bson.E{Key: "a", Value: SortAsc},
		bson.E{Key: "b", Value: SortDesc},
		bson.E{Key: "c", Value: SortAsc},
		bson.E{Key: "d", Value: SortAsc},
		bson.E{Key: "e", Value: SortDesc},
		bson.E{Key: "f", Value: SortDesc})

	data2, err := bson.Marshal(s2)
	if err != nil {
		t.Error(err)
	}

	require.EqualValues(t, data, data2)
}

func TestSortDocUnmarshalBson(t *testing.T) {
	s := SortDocBuilder.OrderBy("a", SortAsc).
		OrderBy("b", SortDesc).
		OrderAscBy("c", "d").
		OrderDescBy("e", "f").
		Doc()

	data, err := bson.Marshal(s)
	if err != nil {
		t.Error(err)
	}

	s1 := sortDoc{}
	s2 := bson.D{}

	if err = bson.Unmarshal(data, &s1); err != nil {
		t.Error(err)
	}
	if err = bson.Unmarshal(data, &s2); err != nil {
		t.Error(err)
	}

	require.NotEmpty(t, s1.ToD())
	require.ElementsMatch(t, s1.ToD(), s2)
}
