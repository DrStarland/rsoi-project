package posts

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
// )

// //nolint:all
// func TestCountPercentage(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()
// 		startPercent := 0
// 		testedPost := Post{
// 			MongoID:          id,
// 			Votes:            []Vote{{UserID: "He", Vote: -1}, {UserID: "Exactly He", Vote: 1}},
// 			UpvotePercentage: int32(startPercent),
// 		}

// 		repo.countPercentage(&testedPost)

// 		assert.Equal(t, int32(50), testedPost.UpvotePercentage)
// 	})

// 	mt.Run("No votes", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()
// 		startPercent := 0
// 		testedPost := Post{
// 			MongoID:          id,
// 			Votes:            []Vote{},
// 			UpvotePercentage: int32(startPercent),
// 		}

// 		repo.countPercentage(&testedPost)

// 		assert.Equal(t, int32(0), testedPost.UpvotePercentage)
// 	})
// }

// //nolint:all
// func TestGetAll(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id1 := primitive.NewObjectID()
// 		id2 := primitive.NewObjectID()
// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id1}, {Key: "Score", Value: int32(10)},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
// 			{Key: "_id", Value: id2}, {Key: "Score", Value: int32(15)},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		posts, err := repo.GetAll()

// 		assert.Nil(t, err)
// 		// проверяем сортировку (что первый возвращённый пост обладает интересующим нас рейтингом)
// 		assert.Equal(t, posts[0].Score, int32(15))
// 	})

// 	mt.Run("Broken mongo", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id1 := primitive.NewObjectID()
// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id1},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: func() {}},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		_, err := repo.GetAll()

// 		assert.NotNil(t, err)
// 	})

// 	mt.Run("cursor error", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))
// 		_, err := repo.GetAll()
// 		assert.NotNil(t, err)
// 	})
// }

// //nolint:all
// func TestHiddenGetByID(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()
// 		expectedPost := Post{
// 			MongoID: id,
// 		}
// 		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: expectedPost.MongoID},
// 		}))

// 		post, err := repo.hiddenGetByID(id.Hex())

// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedPost.MongoID, post.MongoID)
// 	})

// 	mt.Run("No posts", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, nil))

// 		post, err := repo.hiddenGetByID(primitive.NewObjectID().Hex())

// 		assert.NotNil(t, err)
// 		assert.Nil(t, post)
// 	})

// 	mt.Run("Bad id", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)

// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}})

// 		post, err := repo.hiddenGetByID("1")

// 		assert.NotNil(t, err)
// 		assert.Nil(t, post)
// 	})
// }

// // я не знаю, как покрыть остальное
// //
// //nolint:all
// func TestUpViewsByID(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Not Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()
// 		expectedPost := Post{
// 			MongoID: id, ID: id.Hex(), Views: uint32(2),
// 		}

// 		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: expectedPost.MongoID}, {Key: "id", Value: expectedPost.ID}, {Key: "views", Value: uint32(1)},
// 		},
// 			bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}}))

// 		_, err := repo.UpViewsByID(id.Hex())

// 		assert.NotNil(t, err)
// 	})

// 	mt.Run("Not Good ID", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := ""
// 		_, err := repo.UpViewsByID(id)

// 		assert.Error(t, err)
// 	})
// }

// //nolint:all
// func TestGetAllByCategory(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id1 := primitive.NewObjectID()
// 		id2 := primitive.NewObjectID()

// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id1}, {Key: "category", Value: "music"},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
// 			{Key: "_id", Value: id2}, {Key: "category", Value: "music"},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		_, err := repo.GetAllByCategory("music")

// 		assert.Nil(t, err)
// 	})

// 	mt.Run("Cursed database", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()

// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: func() {}},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		_, err := repo.GetAllByCategory("programming")

// 		assert.NotNil(t, err)
// 	})

// 	mt.Run("Error on get", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

// 		_, err := repo.GetAllByCategory("music")

// 		assert.NotNil(t, err)
// 	})
// }

// //nolint:all
// func TestGetAllByUser(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id1 := primitive.NewObjectID()
// 		id2 := primitive.NewObjectID()

// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id1}, {Key: "author", Value: bson.M{
// 				"username": "rvasily",
// 			}},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
// 			{Key: "_id", Value: id2}, {Key: "author", Value: bson.M{
// 				"username": "rvasily",
// 			}},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		_, err := repo.GetAllByUser("rvasily")

// 		assert.Nil(t, err)
// 	})

// 	mt.Run("Cursed database", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		id := primitive.NewObjectID()

// 		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: id},
// 		})
// 		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
// 			{Key: "_id", Value: func() {}},
// 		})
// 		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
// 		mt.AddMockResponses(first, second, killCursors)

// 		_, err := repo.GetAllByUser("lvasily")

// 		assert.NotNil(t, err)
// 	})

// 	mt.Run("Error on get", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

// 		_, err := repo.GetAllByUser("rvasily")

// 		assert.NotNil(t, err)
// 	})
// }

// //nolint:all
// func TestAddPost(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
// 		post := &Post{}

// 		id, err := repo.Add(post)

// 		assert.Nil(t, err)
// 		assert.NotEqual(t, "", id)
// 	})

// 	mt.Run("Error on add", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))
// 		post := &Post{}

// 		_, err := repo.Add(post)

// 		assert.NotNil(t, err)
// 	})
// }

// //nolint:all
// func TestUpdate(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})
// 		post := &Post{ID: primitive.NewObjectID().Hex()}

// 		flag, err := repo.Update(post)

// 		assert.Nil(t, err)
// 		assert.True(t, flag)
// 	})

// 	mt.Run("Error on update", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

// 		post := &Post{}
// 		flag, err := repo.Update(post)
// 		assert.NotNil(t, err)
// 		assert.False(t, flag)
// 	})
// }

// //nolint:all
// func TestDeletePost(t *testing.T) {
// 	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
// 	defer mt.Close()

// 	// нормальный запуск
// 	mt.Run("Good", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

// 		id := primitive.NewObjectID()
// 		flag, err := repo.Delete(id.Hex())
// 		assert.Nil(t, err)
// 		assert.True(t, flag)
// 	})

// 	// некорректный по структуре ID
// 	mt.Run("Bad id", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

// 		id := primitive.NewObjectID()
// 		flag, err := repo.Delete(id.String())
// 		assert.NotNil(t, err)
// 		assert.False(t, flag)
// 	})

// 	mt.Run("Unknown id", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 0}})

// 		flag, err := repo.Delete("N/A")
// 		assert.NotNil(t, err)
// 		assert.False(t, flag)
// 	})

// 	mt.Run("Error on delete", func(mt *mtest.T) {
// 		postsCollection := mt.Coll
// 		repo := NewMongoRepo(postsCollection)
// 		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{}))

// 		id := primitive.NewObjectID().Hex()
// 		flag, err := repo.Delete(id)
// 		assert.NotNil(t, err)
// 		assert.True(t, flag)
// 	})
// }
