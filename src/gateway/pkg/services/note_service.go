package services

import (
	"context"
	"fmt"
	obj "gateway/pkg/models/note"
	"gateway/pkg/myjson"
	"gateway/pkg/utils"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Service encapsulates usecase logic for albums.
type NoteService interface {
	Get(ctx context.Context, id int) (obj.Note, error)
	Query(ctx context.Context, offset, limit int) ([]obj.Note, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input *obj.NoteCreationRequest) (obj.Note, error)
	Update(ctx context.Context, id int, input *obj.NoteCreationRequest) (obj.Note, error)
	Delete(ctx context.Context, id int) (obj.Note, error)

	CreateFromRequest(ctx context.Context, r *http.Request) (obj.Note, error)
	UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Note, error)
	DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Note, error)
}

// NewService creates a new album service.
func NewNoteService(client *http.Client, logger *zap.SugaredLogger) noteService {
	return noteService{client, logger}
}

type noteService struct {
	Client *http.Client
	logger *zap.SugaredLogger
}

func (s noteService) Create(ctx context.Context, req *obj.NoteCreationRequest) (obj.Note, error) {
	return obj.Note{}, errors.New("not implemented")
}

// Update updates the album with the specified ID.
func (s noteService) Update(ctx context.Context, id int, req *obj.NoteCreationRequest) (obj.Note, error) {
	return obj.Note{}, errors.New("not implemented")
}

// Delete deletes the album with the specified ID.
func (s noteService) Delete(ctx context.Context, id int) (obj.Note, error) {
	return obj.Note{}, errors.New("not implemented")
}

// Count returns the number of albums.
func (s noteService) Count(ctx context.Context) (int, error) {
	return 0, errors.New("not implemented")
}

type ByUpdatedAtNote []obj.Note

func (slice ByUpdatedAtNote) Len() int           { return len(slice) }
func (slice ByUpdatedAtNote) Less(i, j int) bool { return slice[i].UpdatedAt < slice[j].UpdatedAt }
func (slice ByUpdatedAtNote) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

// Query returns the albums with the specified offset and limit.
func (s noteService) Query(ctx context.Context, offset, limit int) ([]obj.Note, error) {
	requestURL := fmt.Sprintf("%s/api/v1/notes", utils.Config.NotesEndpoint)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return []obj.Note{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return []obj.Note{}, errors.Wrap(err, "failed request to notes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	items := []obj.Note{}
	if err = myjson.From(body, &items); err != nil {
		return []obj.Note{}, errors.Wrap(err, "failed to decode response")
	}

	sort.Stable(ByUpdatedAtNote(items))
	return items, nil
}

func (s noteService) Get(ctx context.Context, id int) (obj.Note, error) {
	requestURL := fmt.Sprintf("%s/api/v1/notes/%d", utils.Config.NotesEndpoint, id)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		s.logger.Errorln("failed to create an http request")
		return obj.Note{}, errors.Wrap(err, "failed to create an http request")
	}
	req.Header.Set("Authorization", ctx.Value("Authorization").(string))

	res, err := s.Client.Do(req)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "failed request to notes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	note := obj.Note{}
	if err = myjson.From(body, &note); err != nil {
		return obj.Note{}, errors.Wrap(err, "failed to decode response")
	}

	return note, nil
}

func (s noteService) CreateFromRequest(ctx context.Context, r *http.Request) (obj.Note, error) {
	requestURL := fmt.Sprintf("%s/api/v1/notes", utils.Config.NotesEndpoint)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""
	// if err != nil {
	// 	s.logger.Errorln("failed to create an http request")
	// 	return obj.Note{}, errors.Wrap(err, "failed to create an http request")
	// }

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "failed request to notes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	note := obj.Note{}
	if err = myjson.From(body, &note); err != nil {
		return obj.Note{}, errors.Wrap(err, "failed to decode response")
	}

	return note, nil
}

func (s noteService) UpdateFromRequest(ctx context.Context, r *http.Request, id int) (obj.Note, error) {
	requestURL := fmt.Sprintf("%s/api/v1/notes/%d", utils.Config.NotesEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "failed request to notes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	note := obj.Note{}
	if err = myjson.From(body, &note); err != nil {
		return obj.Note{}, errors.Wrap(err, "failed to decode response")
	}

	return note, nil
}

func (s noteService) DeleteFromRequest(ctx context.Context, r *http.Request, id int) (obj.Note, error) {
	requestURL := fmt.Sprintf("%s/api/v1/notes/%d", utils.Config.NotesEndpoint, id)

	URL, _ := url.Parse(requestURL)
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host
	r.URL.Path = URL.Path
	r.RequestURI = ""

	res, err := s.Client.Do(r)
	if err != nil {
		return obj.Note{}, errors.Wrap(err, "failed request to notes service")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Errorln(err.Error())
	}
	res.Body.Close()

	note := obj.Note{}
	if err = myjson.From(body, &note); err != nil {
		if res.StatusCode == 204 {
			return obj.Note{}, errors.New("not exist")
		}
		return obj.Note{}, errors.Wrap(err, "failed to decode response")
	}

	return note, nil
}
