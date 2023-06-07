package comments

// import (
// 	"redditclone/pkg/model/user"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// type Comment struct {
// 	Created string             `json:"created"`
// 	MongoID primitive.ObjectID `json:"-" bson:"_id"`
// 	ID      string             `json:"id"`
// 	Author  user.User          `json:"author"`
// 	Body    string             `json:"body"`
// }

// //go:generate mockgen -source=comment.go -destination=repo_mock.go -package=comments CommentsRepo
// type CommentsRepo interface {
// 	GetByID(commentID string) (*Comment, error)
// 	Add(comment *Comment) (string, error)
// 	DeleteFromRepo(commentID string) (bool, error)
// 	// для прикрепления комментариев к чему-то
// 	AttachTo(commentID string, target Commentable) error
// 	RemoveFrom(commentID string, target Commentable) error
// }

// type Commentable interface {
// 	Attach(comment *Comment) error
// 	Unpin(commentID string) error
// }
