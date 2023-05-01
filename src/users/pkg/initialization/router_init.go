package initialization

import (
	"log"
	"net/http"
	"os"
	"strings"
	"users/pkg/database"
	"users/pkg/dbcontext"
	"users/pkg/handlers"
	mid "users/pkg/middleware"
	"users/pkg/models/note"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type App struct {
	Router        *httprouter.Router
	DB            *dbx.DB
	Logger        *zap.SugaredLogger
	ServerAddress string
}

func NewApp(logger *zap.SugaredLogger) App {
	return App{Logger: logger}.initDB().initRouter().initAddress()
}

func GetServerAddress() string {
	ServerAddress := os.Getenv("PORT")
	if ServerAddress == "" || ServerAddress == ":80" {
		ServerAddress = ":8080"
	} else if !strings.HasPrefix(ServerAddress, ":") {
		ServerAddress = ":" + ServerAddress
	}
	return ServerAddress
}

func (app App) initAddress() App {
	ServerAddress := GetServerAddress()
	app.ServerAddress = ServerAddress
	return app
}

func HealthOK(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func (app App) initDB() App {
	db, err := database.CreateConnection()
	if err != nil {
		app.Logger.DPanicln("Connection was not successfull: "+err.Error(),
			"type", "START",
		)
	}

	app.Logger.Infow("Connection to db was successfull",
		"type", "START",
	)

	app.DB = db
	return app
}

func (app App) initRouter() App {

	// repoTicket := ticket.NewPostgresRepo(db)
	tak := dbcontext.New(app.DB)
	repoNotes := note.NewRepository(tak, app.Logger)

	noteHandler := &handlers.NoteMainHandler{
		Logger:  app.Logger,
		Repo:    repoNotes,
		Service: handlers.NewNoteService(repoNotes, *app.Logger),
	}

	// ticketHandler := &handlers.TicketsHandler{
	// 	Logger:      logger,
	// 	TicketsRepo: repoTicket,
	// }

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	router.GET("/api/v1/notes", mid.AccessLog(noteHandler.List, app.Logger))
	router.GET("/api/v1/users/:id", mid.AccessLog(noteHandler.Show, app.Logger))

	router.GET("/api/v1/tickets/:username", mid.AccessLog(HealthOK, app.Logger))
	router.DELETE("/api/v1/tickets/:ticketUID", mid.AccessLog(HealthOK, app.Logger))

	router.GET("/manage/health", HealthOK)

	app.Router = router

	return app
}

func (app App) Stop() {
	app.DB.Close()
}
