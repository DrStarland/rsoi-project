package posts

// type Post struct {
// 	ID               string             `json:"id" bson:"id"`                             // ID
// 	Score            int32              `json:"score" bson:"score"`                       // score
// 	Type             string             `json:"type" valid:"in(link|text)"`               // type
// 	Views            uint32             `json:"views" bson:"views"`                       // views
// 	Title            string             `json:"title"`                                    // title
// 	URL              string             `json:"url,omitempty"`                            // url
// 	Author           user.User          `json:"author"`                                   // author
// 	Category         string             `json:"category"`                                 // category
// 	Text             string             `json:"text"`                                     // text
// 	Votes            []Vote             `json:"votes" bson:"votes"`                       // votes
// 	Comments         []comments.Comment `json:"comments" bson:"comments"`                 // comments
// 	Created          string             `json:"created"`                                  // created
// 	UpvotePercentage int32              `json:"upvotePercentage" bson:"upvotePercentage"` // upvotePercentage
// 	mu               *sync.Mutex
// }

// type Vote struct {
// 	UserID string `json:"user"` // user
// 	Vote   int32  `json:"vote"` // vote
// }

// //go:generate mockgen -source=post.go -destination=repo_mock.go -package=posts PostsRepo
// type PostsRepo interface {
// 	GetAll() ([]*Post, error)
// 	GetAllByCategory(category string) ([]*Post, error)
// 	GetAllByUser(username string) ([]*Post, error)
// 	GetByID(id string) (*Post, error)
// 	GetByIDWithUpvote(id, userID string) (*Post, error)
// 	GetByIDWithUndoVote(id, userID string) (*Post, error)
// 	GetByIDWithDownvote(id, userID string) (*Post, error)
// 	Add(post *Post) (string, error)
// 	Update(newPost *Post) (bool, error)
// 	Delete(id string) (bool, error)
// 	AddComment(post *Post, comment *comments.Comment) (string, error)
// 	DeleteComment(commentid string, post *Post) (bool, error)
// }

// type Post struct {
// 	ID               string             `json:"id" bson:"id"`                             // ID
// 	Score            int32              `json:"score" bson:"score"`                       // score
// 	Type             string             `json:"type" valid:"in(link|text)"`               // type
// 	Views            uint32             `json:"views" bson:"views"`                       // views
// 	Title            string             `json:"title"`                                    // title
// 	URL              string             `json:"url,omitempty"`                            // url
// 	Author           user.User          `json:"author"`                                   // author
// 	Category         string             `json:"category"`                                 // category
// 	Text             string             `json:"text"`                                     // text
// 	Votes            []Vote             `json:"votes" bson:"votes"`                       // votes
// 	Comments         []comments.Comment `json:"comments" bson:"comments"`                 // comments
// 	Created          string             `json:"created"`                                  // created
// 	UpvotePercentage int32              `json:"upvotePercentage" bson:"upvotePercentage"` // upvotePercentage
// 	mu               *sync.Mutex
// }

// type Vote struct {
// 	UserID string `json:"user"` // user
// 	Vote   int32  `json:"vote"` // vote
// }

// //go:generate mockgen -source=post.go -destination=repo_mock.go -package=posts PostsRepo
// type PostsRepo interface {
// 	GetAll() ([]*Post, error)
// 	GetAllByCategory(category string) ([]*Post, error)
// 	GetAllByUser(username string) ([]*Post, error)
// 	GetByID(id string) (*Post, error)
// 	GetByIDWithUpvote(id, userID string) (*Post, error)
// 	GetByIDWithUndoVote(id, userID string) (*Post, error)
// 	GetByIDWithDownvote(id, userID string) (*Post, error)
// 	Add(post *Post) (string, error)
// 	Update(newPost *Post) (bool, error)
// 	Delete(id string) (bool, error)
// 	AddComment(post *Post, comment *comments.Comment) (string, error)
// 	DeleteComment(commentid string, post *Post) (bool, error)
// }
