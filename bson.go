package hamster

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bsoner interface {
	ToM() primitive.M
	ToD() primitive.D
}