package posts

// import (
// 	"redditclone/pkg/model/comments"
// 	"sync"
// )

// func (post *Post) InitMutex() {
// 	if post.mu == nil {
// 		post.mu = &sync.Mutex{}
// 	}
// }

// func (post *Post) Attach(comment *comments.Comment) error {
// 	post.Comments = append(post.Comments, *comment)
// 	return nil
// }

// func (post *Post) Unpin(commentID string) error {
// 	i := -1
// 	for idx, comment := range post.Comments {
// 		if comment.ID == commentID {
// 			i = idx
// 			break
// 		}
// 	}
// 	if i < 0 {
// 		return nil
// 	}

// 	if i < len(post.Comments)-1 {
// 		copy(post.Comments[i:], post.Comments[i+1:])
// 	}

// 	post.Comments = post.Comments[:len(post.Comments)-1]
// 	return nil
// }
