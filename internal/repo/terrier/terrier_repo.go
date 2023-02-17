package terrier

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const TerrierCollection = "terrier"

type ITerrierRepo interface {
	QueryTerriers(ctx context.Context) ([]*Terrier, error)
	UpdateTerrierPublished(terrier *Terrier) error
}

type Repo struct {
	*mongo.Database
}

func (repo *Repo) QueryTerriers(ctx context.Context) ([]*Terrier, error) {
	collection := repo.Collection(TerrierCollection)
	results := make([]*Terrier, 0)

	filter := bson.M{
		"DogName": "June",
	}

	opts := options.Find().SetSort(bson.D{{"_id", 1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = cursor.Close(ctx)
	}()

	for {
		if cursor.Err() != nil {
			return results, err
		}
		if ok := cursor.Next(ctx); !ok {
			break
		}

		t := &Terrier{}
		if err := cursor.Decode(t); err != nil {
			return nil, err
		}
		results = append(results, t)
	}

	return results, err
}

func (repo *Repo) UpdateTerrierPublished(terrier *Terrier) error {
	log.Printf("Updated terrier %v", terrier.ID)
	_, err := repo.Collection(TerrierCollection).UpdateOne(
		context.Background(),
		bson.D{{"_id", terrier.ID}},
		bson.D{{"$set",
			bson.D{{"Publish_status", true}},
		}},
	)
	if err != nil {
		return fmt.Errorf("[%s] not updated in mongo: %w", terrier.ID, err)
	}
	return nil
}
