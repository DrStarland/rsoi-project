package posts

// import (
// 	"errors"
// 	"log"
// 	"redditclone/pkg/model/comments"
// 	"sort"

// 	"context"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
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

// func (repo *MongoRepository) countPercentage(post *Post) {
// 	count := 0
// 	for _, val := range post.Votes {
// 		if val.Vote == 1 {
// 			count++
// 		}
// 	}
// 	if n := len(post.Votes); n != 0 {
// 		post.UpvotePercentage = int32(100 * count / n)
// 	} else {
// 		post.UpvotePercentage = 0
// 	}
// }

// func (repo *MongoRepository) UpViewsByID(id string) (*Post, error) {
// 	post, err := repo.hiddenGetByID(id)
// 	if err != nil {
// 		return post, err
// 	}

// 	post.Views++

// 	_, err = repo.Update(post)
// 	if err != nil {
// 		return nil, errors.New("Error with getting post by ID: " + err.Error())
// 	}

// 	return post, nil
// }

// func (repo *MongoRepository) GetAll() ([]*Post, error) {
// 	posts := make([]*Post, 0)

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	cursor, err := repo.source.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = cursor.All(ctx, &posts); err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}

// 	f := func(i, j int) bool { return posts[i].Score > posts[j].Score }
// 	sort.Slice(posts, f)

// 	return posts, nil
// }

// func (repo *MongoRepository) GetByID(id string) (*Post, error) {
// 	return repo.UpViewsByID(id)
// }

// func (repo *MongoRepository) hiddenGetByID(id string) (*Post, error) {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, errors.New("bad ID")
// 	}

// 	post := &Post{}

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	result := repo.source.FindOne(ctx, bson.M{"_id": objectID})
// 	err = result.Decode(post)
// 	if err != nil {
// 		return nil, errors.New("post not found")
// 	}

// 	return post, nil
// }

// func (repo *MongoRepository) GetAllByCategory(category string) ([]*Post, error) {
// 	x := make([]*Post, 0)

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	cursor, err := repo.source.Find(ctx, bson.M{"category": category})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = cursor.All(ctx, &x); err != nil {
// 		return nil, err
// 	}

// 	return x, nil
// }

// func (repo *MongoRepository) GetAllByUser(username string) ([]*Post, error) {
// 	x := make([]*Post, 0)
// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	cursor, err := repo.source.Find(ctx, bson.M{"author.username": username})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = cursor.All(ctx, &x); err != nil {
// 		return nil, err
// 	}
// 	return x, nil
// }

// func (repo *MongoRepository) GetByIDWithUpvote(id, userID string) (*Post, error) {
// 	var index int
// 	alreadyVoted := false
// 	item := Vote{UserID: userID, Vote: 1}

// 	post, err := repo.hiddenGetByID(id)
// 	if post == nil || err != nil {
// 		return nil, err
// 	}

// 	for i, vote := range post.Votes {
// 		if vote.UserID == userID {
// 			alreadyVoted = true
// 			index = i
// 			break
// 		}
// 	}

// 	if alreadyVoted {
// 		post.Votes[index].Vote = 1
// 		post.Score += 2
// 	} else {
// 		post.Votes = append(post.Votes, item)
// 		post.Score += item.Vote
// 	}

// 	repo.countPercentage(post)
// 	_, err = repo.Update(post)
// 	return post, err
// }

// func (repo *MongoRepository) GetByIDWithUndoVote(id, userID string) (*Post, error) {
// 	var index int

// 	post, err := repo.hiddenGetByID(id)
// 	if post == nil {
// 		return nil, err
// 	}

// 	for i, vote := range post.Votes {
// 		if vote.UserID == userID {
// 			index = i
// 			post.Score -= vote.Vote
// 			post.Votes = append((post.Votes)[:index], (post.Votes)[index+1:]...)
// 		}
// 	}
// 	repo.countPercentage(post)
// 	_, err = repo.Update(post)
// 	return post, err
// }

// func (repo *MongoRepository) GetByIDWithDownvote(id, userID string) (*Post, error) {
// 	var index int
// 	alreadyVoted := false
// 	item := Vote{UserID: userID, Vote: -1}

// 	post, err := repo.hiddenGetByID(id)
// 	if post == nil {
// 		return nil, err
// 	} else if err != nil {
// 		log.Println("&&& " + err.Error())
// 		return nil, err
// 	}

// 	for i, vote := range post.Votes {
// 		if vote.UserID == userID {
// 			alreadyVoted = true
// 			index = i
// 			break
// 		}
// 	}

// 	if alreadyVoted {
// 		post.Votes[index].Vote = -1
// 		post.Score -= 2
// 	} else {
// 		post.Votes = append(post.Votes, item)
// 		post.Score += item.Vote
// 	}

// 	repo.countPercentage(post)
// 	_, err = repo.Update(post)
// 	return post, err
// }

// func (repo *MongoRepository) Add(post *Post) (string, error) {
// 	post.MongoID = primitive.NewObjectID()
// 	post.ID = post.MongoID.Hex()

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	insertResult, err := repo.source.InsertOne(ctx, post)
// 	if err != nil {
// 		return "", err
// 	}

// 	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
// }

// func (repo *MongoRepository) Update(newPost *Post) (bool, error) {
// 	id, err := primitive.ObjectIDFromHex(newPost.ID)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return false, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	_, err = repo.source.UpdateOne(ctx,
// 		bson.M{"_id": id},
// 		bson.D{
// 			{Key: "$set", Value: bson.M{
// 				"score":            newPost.Score,
// 				"views":            newPost.Views,
// 				"upvotePercentage": newPost.UpvotePercentage,
// 				"votes":            newPost.Votes,
// 				"comments":         newPost.Comments,
// 			}},
// 		})
// 	if err != nil {
// 		log.Println("!!!! " + err.Error())
// 		return false, err
// 	}

// 	return true, nil
// }

// func (repo *MongoRepository) Delete(id string) (bool, error) {
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return false, errors.New("bad ID")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
// 	defer cancel()
// 	_, err = repo.source.DeleteOne(ctx, bson.M{"_id": objectID})
// 	if err != nil {
// 		log.Println(err.Error())
// 		return true, err
// 	}

// 	return true, nil
// }

// func (repo *MongoRepository) AddComment(post *Post, comment *comments.Comment) (string, error) {
// 	err := post.Attach(comment)
// 	if err != nil {
// 		return "", err
// 	}
// 	_, err = repo.Update(post)
// 	if err != nil {
// 		return "", err
// 	}
// 	return comment.ID, err
// }

// func (repo *MongoRepository) DeleteComment(commentID string, post *Post) (bool, error) {
// 	err := post.Unpin(commentID)
// 	if err != nil {
// 		return false, nil
// 	}
// 	flag, err := repo.Update(post)
// 	return flag, err
// }
