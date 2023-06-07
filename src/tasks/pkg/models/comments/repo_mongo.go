package comments

// import (
// 	"context"
// 	"errors"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"gopkg.in/mgo.v2/bson"
// )

// type MongoRepository struct {
// 	source *mongo.Collection
// }

// const (
// 	mongoTimeout = 10 * time.Second
// )

// func NewMongoRepo(db *mongo.Collection) *MongoRepository {
// 	return &MongoRepository{
// 		source: db,
// 	}
// }

// func (repo *MongoRepository) GetByID(id string) (*Comment, error) {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, errors.New("bad ID")
// 	}

// 	comment := &Comment{}

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	result := repo.source.FindOne(ctx, bson.M{"_id": objectID})
// 	err = result.Decode(comment)
// 	if err != nil {
// 		return nil, errors.New("post not found")
// 	}

// 	return comment, nil
// }

// func (repo *MongoRepository) AttachTo(commentID string, target Commentable) error {
// 	comment, err := repo.GetByID(commentID)
// 	if err != nil {
// 		return err
// 	}
// 	err = target.Attach(comment)
// 	return err
// }

// func (repo *MongoRepository) RemoveFrom(commentID string, target Commentable) error {
// 	if err := target.Unpin(commentID); err != nil {
// 		return err
// 	}
// 	_, err := repo.DeleteFromRepo(commentID)
// 	return err
// }

// func (repo *MongoRepository) Add(comment *Comment) (string, error) {
// 	comment.MongoID = primitive.NewObjectID()
// 	comment.ID = comment.MongoID.Hex()

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	insertResult, err := repo.source.InsertOne(ctx, comment)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return "", err
// 	}

// 	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
// }

// func (repo *MongoRepository) DeleteFromRepo(commentID string) (bool, error) {
// 	objectID, err := primitive.ObjectIDFromHex(commentID)
// 	if err != nil {
// 		return false, errors.New("bad ID")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	_, err = repo.source.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		log.Println(err.Error())
// 		return false, err
// 	}

// 	return true, nil
// }
