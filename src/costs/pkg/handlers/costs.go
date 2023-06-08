package handlers

import (
	"context"
	"costs/pkg/models/cost"
	"costs/pkg/myjson"
	"costs/pkg/services"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CostsHandler interface {
	List(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type CostMainHandler struct {
	Logger  *zap.SugaredLogger
	Repo    cost.Repository
	Service services.CostService
}

func (h CostMainHandler) List(w http.ResponseWriter, r *http.Request) {
	// lol := ps.ByName("id")
	elems, err := h.Service.Query(r.Context(), 0, 64)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusOK, elems)
}

func (h CostMainHandler) Show(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	cost, err := h.Service.Get(r.Context(), id)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		return
	}

	myjson.JSONResponce(w, http.StatusOK, cost)
}

func (h *CostMainHandler) Add(w http.ResponseWriter, r *http.Request) {
	cost := &cost.CostCreationRequest{}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		myjson.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = myjson.From(body, cost); err != nil {
		myjson.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	test := context.WithValue(
		context.WithValue(
			r.Context(), "X-UID", r.Header.Get("X-UID"),
		), "X-User-Name", r.Header.Get("X-User-Name"))
	res, err := h.Service.Create(test, cost)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h CostMainHandler) Update(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
	}

	_, err = h.Service.Get(r.Context(), id)
	if err != nil {
		myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		return
	}

	costRequest := &cost.CostCreationRequest{}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		myjson.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = myjson.From(body, costRequest); err != nil {
		myjson.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	test := context.WithValue(
		context.WithValue(r.Context(),
			"X-UID", r.Header.Get("X-UID")),
		"X-User-Name", r.Header.Get("X-User-Name"))
	res, err := h.Service.Update(test, id, costRequest)
	if err != nil {
		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	myjson.JSONResponce(w, http.StatusCreated, res)
}

func (h CostMainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		myjson.JSONResponce(w, http.StatusBadRequest, errors.Wrap(err, "bad ID in URL"))
		return
	}

	cost, err := h.Service.Delete(r.Context(), id)

	if err != nil {
		switch err.Error() {
		case "not exist":
			myjson.JSONResponce(w, http.StatusNoContent, errors.Wrap(err, ""))
			return
		default:
			myjson.JSONResponce(w, http.StatusInternalServerError, errors.Wrap(err, ""))
		}
		return
	}

	myjson.JSONResponce(w, http.StatusNoContent, cost)
}

// func (h CostMainHandler) Create(w http.ResponseWriter, r *http.Request) error {
// 	var input CreateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		r.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}
// 	album, err := r.service.Create(c.Request.Context(), input)
// 	if err != nil {
// 		return err
// 	}

// 	myjson.JSONResponce(cost, http.StatusCreated)
// 	return c.WriteWithStatus(album, http.StatusCreated)
// }

// func (r resource) update(c *routing.Context) error {
// 	var input UpdateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		r.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}

// 	album, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }

// func (r resource) delete(c *routing.Context) error {
// 	album, err := r.service.Delete(c.Request.Context(), c.Param("id"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }

// func (h CostMainHandler) create(c *routing.Context) error {
// 	var input CreateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		r.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}
// 	album, err := h.Repo.Create(c.Request.Context(), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.WriteWithStatus(album, http.StatusCreated)
// }

// func (h CostMainHandler) update(c *routing.Context) error {
// 	var input UpdateAlbumRequest
// 	if err := c.Read(&input); err != nil {
// 		r.logger.With(c.Request.Context()).Info(err)
// 		return errors.BadRequest("")
// 	}

// 	album, err := h.Repo.Update(c.Request.Context(), c.Param("id"), input)
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }

// func (h CostMainHandler) delete(c *routing.Context) error {
// 	album, err := h.Repo.Delete(c.Request.Context(), c.Param("id"))
// 	if err != nil {
// 		return err
// 	}

// 	return c.Write(album)
// }

// type PostsHandler struct {
// 	PostsRepo   posts.PostsRepo
// 	CommentRepo comments.CommentsRepo
// 	Logger      *zap.SugaredLogger
// }

// func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context")
// 	}

// 	post, err := h.PostsRepo.GetByID(id)

// 	switch {
// 	case err != nil:
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	case post == nil:
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	case post.Author.ID != currentSess.User.ID:
// 		myjson.JSONError(w, http.StatusBadRequest, "FORBIDDEN")
// 	}

// 	_, err = h.PostsRepo.Delete(id)
// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, err.Error())
// 	}

// 	myjson.JSONError(w, http.StatusOK, "success")
// }

// func (h *PostsHandler) GetOne(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	item, err := h.PostsRepo.GetByID(id)

// 	switch {
// 	case err != nil:
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error: "+err.Error())
// 		return
// 	case item == nil:
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, item)
// }

// func (h *PostsHandler) GetOneAndUpvote(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context")
// 	}

// 	item, err := h.PostsRepo.GetByIDWithUpvote(id, currentSess.User.ID)
// 	switch {
// 	case err != nil:
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	case item == nil:
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, item)
// }

// func (h *PostsHandler) GetOneAndUndoVote(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context")
// 	}

// 	item, err := h.PostsRepo.GetByIDWithUndoVote(id, currentSess.User.ID)
// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	}
// 	if item == nil {
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, item)
// }

// func (h *PostsHandler) GetOneAndDownvote(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context")
// 	}

// 	item, err := h.PostsRepo.GetByIDWithDownvote(id, currentSess.User.ID)
// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	}
// 	if item == nil {
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, item)
// }

// func (h *PostsHandler) GetAllByCategory(w http.ResponseWriter, r *http.Request) {
// 	category := ps.ByName("category")
// 	elems, err := h.PostsRepo.GetAllByCategory(category)

// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, elems)
// }

// func (h *PostsHandler) GetAllByUser(w http.ResponseWriter, r *http.Request) {
// 	username := ps.ByName("username")
// 	elems, err := h.PostsRepo.GetAllByUser(username)

// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, elems)
// }

// func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
// 	if r.Header.Get("Content-Type") != Applijson {
// 		myjson.JSONError(w, http.StatusBadRequest, "unknown payload")
// 		return
// 	}

// 	id := ps.ByName("id")
// 	post, err := h.PostsRepo.GetByID(id)

// 	if err != nil {
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	}
// 	if post == nil {
// 		myjson.JSONError(w, http.StatusNotFound, "post not found")
// 		return
// 	}

// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context in comment add")
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	r.Body.Close()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	comment := &comments.Comment{}

// 	f := map[string]interface{}{}
// 	err = myjson.From(body, &f)
// 	if err != nil {
// 		myjson.JSONError(w, http.StatusBadRequest, "cant unpack payload")
// 		return
// 	}

// 	if f["comment"] != nil {
// 		comment.Body = f["comment"].(string)
// 		comment.Created = time.Now().Format("2006-01-02T15:04:05.000")
// 		comment.Author = currentSess.User

// 		commentID, err := h.CommentRepo.Add(comment)
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		createdComm, err := h.CommentRepo.GetByID(commentID)
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		_, err = h.PostsRepo.AddComment(post, createdComm)
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 	}

// 	myjson.JSONResponce(w, http.StatusCreated, post)
// }

// func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
// 	id := ps.ByName("id")
// 	currentSess, ok := r.Context().Value(session.SessionKey).(*session.Session)
// 	if !ok {
// 		myjson.JSONError(w, http.StatusInternalServerError, "broken context")
// 	}

// 	post, err := h.PostsRepo.GetByID(id)

// 	switch {
// 	case err != nil:
// 		myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 		return
// 	case post == nil:
// 		myjson.JSONError(w, http.StatusNotFound, "comment not found")
// 		return
// 	}

// 	commentID := ps.ByName("commentid")

// 	for _, comment := range post.Comments {
// 		if comment.ID == commentID {
// 			if comment.Author.ID != currentSess.User.ID {
// 				myjson.JSONError(w, http.StatusBadRequest, "FORBIDDEN")
// 			} else {
// 				_, err = h.PostsRepo.DeleteComment(commentID, post)
// 				if err != nil {
// 					myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 				}
// 				_, err = h.CommentRepo.DeleteFromRepo(commentID)
// 				if err != nil {
// 					myjson.JSONError(w, http.StatusInternalServerError, "DB error")
// 				}
// 				break
// 			}
// 		}
// 	}

// 	myjson.JSONResponce(w, http.StatusOK, post)
// }

// func (h *TicketsHandler) GetTicketsByUsername(w http.ResponseWriter, r *http.Request) {
// 	username := ps.ByName("username")
// 	tickets, err := h.TicketsRepo.GetByUsername(username)
// 	if err != nil {
// 		log.Printf("Failed to get ticket: %s", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Add("Content-Type", "application/json")
// 	myjson.JSONResponce(w, http.StatusOK, tickets)
// }

// func (h *TicketsHandler) BuyTicket(w http.ResponseWriter, r *http.Request) {
// 	if r.Header.Get("Content-Type") != "application/json" {
// 		myjson.JSONError(w, http.StatusBadRequest, "unknown payload")
// 		return
// 	}

// 	body, _ := io.ReadAll(r.Body)
// 	r.Body.Close()

// 	ticket := &ticket.Ticket{}
// 	err := myjson.From(body, ticket)

// 	if err != nil {
// 		h.Logger.Errorln("STRANDING ", string(body))
// 		myjson.JSONError(w, http.StatusBadRequest, "cant unpack payload")
// 		return
// 	}

// 	if err := h.TicketsRepo.Add(ticket); err != nil {
// 		log.Printf("Failed to create ticket: %s", err)

// 		myjson.JSONError(w, http.StatusInternalServerError, "Failed to create ticket: "+err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func (h *TicketsHandler) DeleteTicket(w http.ResponseWriter, r *http.Request) {
// 	ticketUID := ps.ByName("ticketUID")

// 	if err := h.TicketsRepo.Delete(ticketUID); err != nil {
// 		h.Logger.Errorln("Failed to create ticket: " + err.Error())

// 		myjson.JSONError(w, http.StatusInternalServerError, "failed to create ticket: "+err.Error())
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func SearchServer(w http.ResponseWriter, r *http.Request) {
// 	if !checkToken(r.Header) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		if _, err := w.Write([]byte("Неправильный токен!")); err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// log.Println(err.Error())
// 		}
// 		return
// 	}

// 	req, err := parseRequest(r.URL.Query())
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		js, nestedErr := ToJSON(SearchErrorResponse{Error: err.Error()})
// 		if nestedErr != nil {
// 			if _, deepErr := w.Write([]byte(textBadJSON)); deepErr != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				// log.Println(deepErr.Error())
// 			}
// 		}
// 		if _, nestedErr := w.Write(js); nestedErr != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// log.Println(err.Error())
// 		}
// 		return
// 	}

// 	UserInfoStorage, err := ParseDataFromFile(PathToDataset)
// 	if err != nil {
// 		// log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		js, nestedErr := ToJSON(SearchErrorResponse{Error: "Ошибка чтения из файла."})
// 		if nestedErr != nil {
// 			if _, deepErr := w.Write([]byte(textBadJSON)); deepErr != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				// log.Println(err.Error())
// 			}
// 		}
// 		if _, nestedErr := w.Write(js); nestedErr != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// log.Println(err.Error())
// 		}
// 		return
// 	}
// 	UserStorage = UserInfoStorage.toUsers()

// 	result := UserStorage.FindByQueryAndGetSlice(req.Query).Sort(req.OrderField, req.OrderBy).DoOffset(req.Offset).CutToLimit(req.Limit)

// 	bdata, err := ToJSON(result)
// 	// log.Println(bdata)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		if _, deepErr := w.Write([]byte(textBadJSON)); deepErr != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// log.Println(deepErr.Error())
// 		}
// 	} else {
// 		if _, err = w.Write(bdata); err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			// log.Println(err.Error())
// 		} else {
// 			w.WriteHeader(http.StatusOK)
// 		}
// 	}
// 	w.WriteHeader(http.StatusInternalServerError)
// }

// func checkToken(head http.Header) bool {
// 	if token := head.Get("AccessToken"); len(token) != 0 {
// 		// log.Printf("Token: %s\n", token)
// 		if token != SecretKey {
// 			return false
// 		}
// 	}
// 	return true
// }

// func parseRequest(src url.Values) (SearchRequest, error) {
// 	var (
// 		order, offset, limit int
// 		err                  error
// 	)

// 	var req SearchRequest

// 	if order, err = strconv.Atoi(src.Get("order_by")); err != nil {
// 		// log.Println(errors.New("Empty order_by"))
// 		return req, errors.New("empty order_by")
// 	}

// 	if offset, err = strconv.Atoi(src.Get("offset")); err != nil {
// 		// log.Println(err.Error())
// 		return req, errors.New("empty offset")
// 	}

// 	if limit, err = strconv.Atoi(src.Get("limit")); err != nil {
// 		// log.Println(err.Error())
// 		return req, errors.New("empty limit")
// 	}
// 	req = SearchRequest{
// 		Query:      src.Get("query"),
// 		OrderField: src.Get("order_field"),
// 		OrderBy:    order,
// 		Offset:     offset,
// 		Limit:      limit,
// 	}

// 	switch req.OrderField {
// 	case caseID:
// 	case caseAge:
// 	case caseName:
// 	case "":
// 		req.OrderField = caseName
// 	default:
// 		return req, errors.New(ErrorBadOrderField)
// 	}

// 	return req, err
// }
