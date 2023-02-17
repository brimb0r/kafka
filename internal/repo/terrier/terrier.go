package terrier

import "go.mongodb.org/mongo-driver/bson/primitive"

type Terrier struct {
	ID             string `bson:"_id"`
	DogName        string
	Activity       string
	Publish_status bool
	LastUpdated    primitive.DateTime `bson:"last_updated,omitempty"`
}
