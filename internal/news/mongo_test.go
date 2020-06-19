package news

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewMongoQuery(t *testing.T) {
	oID := primitive.NewObjectID()
	cases := []struct {
		name string
		q    Query
		want mongoQuery
	}{
		{
			name: "Given valid hex string within string slice of filter field",
			q: Query{
				Filter: Filter{
					IDs: []string{oID.Hex()},
				},
			},
			want: mongoQuery{
				mongoFilter: mongoFilter{
					IDs: []primitive.ObjectID{oID},
				},
			},
		},
		{
			name: "Given invalid hex string within string slice of filter field",
			q: Query{
				Filter: Filter{
					IDs: []string{"invalidhex"},
				},
			},
			want: mongoQuery{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewMongoQuery(&tc.q)
			if !reflect.DeepEqual(*got, tc.want) {
				t.Errorf("expected mongo query %+v, got %+v", tc.want, *got)
			}
		})
	}
}
