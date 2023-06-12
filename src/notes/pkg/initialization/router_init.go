package initialization

import (
	"fmt"
	"log"
	"net/http"
	"notes/pkg/handlers"
	mid "notes/pkg/middleware"
	"notes/pkg/models/note"
	"notes/pkg/services"
	"notes/pkg/utils"
	"os"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/julienschmidt/httprouter"
	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
	"github.com/zc2638/swag/option"
	"go.uber.org/zap"
)

type App struct {
	Config        utils.Configuration
	Router        *httprouter.Router
	DB            *dbx.DB
	Logger        *zap.SugaredLogger
	ServerAddress string
}

func NewApp(logger *zap.SugaredLogger) App {
	utils.InitConfig()
	return App{Config: utils.Config, Logger: logger}.initDB().initRouter().initAddress()
}

func (app App) GetServerAddress() string {
	ServerAddress := os.Getenv("PORT")
	if ServerAddress == "" || ServerAddress == ":80" {
		ServerAddress = fmt.Sprintf(":%d", app.Config.Port)
	} else if !strings.HasPrefix(ServerAddress, ":") {
		ServerAddress = ":" + ServerAddress
	}
	return ServerAddress
}

func (app App) initAddress() App {
	app.ServerAddress = app.GetServerAddress()
	return app
}

func HealthOK(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app App) initDB() App {
	// db, err := database.CreateConnection()
	// if err != nil {
	// 	app.Logger.DPanicln("Connection was not successfull: "+err.Error(),
	// 		"type", "START",
	// 	)
	// }

	// app.Logger.Infow("Connection to db was successfull",
	// 	"type", "START",
	// )

	// app.DB = db
	return app
}

func (app App) initRouter() App {

	// repoTicket := ticket.NewPostgresRepo(db)
	// tak := dbcontext.New(app.DB)
	repoNotes := note.NewMemoryRepository(app.Logger)

	noteHandler := &handlers.NoteMainHandler{
		Logger:  app.Logger,
		Repo:    repoNotes,
		Service: services.NewNoteService(repoNotes, *app.Logger),
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

	api := swag.New(
		option.Title("Note Service API Doc"),
		option.Security("Sophisticated_Service_auth", "user", "admin"),
		option.SecurityScheme("Sophisticated_Service_auth",
			// option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			option.APIKeySecurity("Authorization", "header"),
			option.OAuth2Scope("admin", ""),
			option.OAuth2Scope("user", ""),
		),
		option.BasePath("/api/v1"),
	)

	api.AddTag("Healthcheck and statistics", "")
	api.AddTag("Notes", "")

	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(mid.AccessLog(HealthOK, app.Logger)),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),

		endpoint.New(
			http.MethodGet, "/notes",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.List, app.Logger), app.Logger),
			),
			endpoint.Summary("Возвращает список заметок"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodGet, "/notes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Show, app.Logger), app.Logger),
			),
			endpoint.Summary("Получение заметки по ID"),
			endpoint.Tags("Notes"),
			endpoint.Path("id", "integer", "ID of note to return", true),
			endpoint.Response(http.StatusOK, "successful operation", endpoint.SchemaResponseOption(note.Note{})),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPost, "/notes",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Add, app.Logger), app.Logger),
			),
			endpoint.Summary("Создание новой заметки"),
			endpoint.Body(note.NoteCreationRequest{}, "Структура запроса на создание заметки", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(note.Note{})),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/notes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Update, app.Logger), app.Logger),
			),
			endpoint.Summary("Редактирование существующей заметки"),
			endpoint.Path("id", "integer", "ID of note to edit", true),
			endpoint.Body(note.NoteCreationRequest{},
				"Структура запроса на изменение заметки", true),
			endpoint.Response(http.StatusCreated, "was successful", endpoint.SchemaResponseOption(note.Note{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodDelete, "/notes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Delete, app.Logger), app.Logger),
			),
			endpoint.Summary("Удаление заметки"),
			endpoint.Path("id", "integer", "ID of note to delete", true),
			endpoint.Response(http.StatusNoContent, "successful"),
			endpoint.Response(http.StatusNoContent, "Entity is not exist or already deleted", endpoint.SchemaResponseOption("not exist")),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
	)

	swag.New()

	api.Walk(func(path string, e *swag.Endpoint) {
		h := e.Handler.(http.Handler)
		path = swag.ColonPath(path)
		router.Handler(e.Method, path, h)
	})

	router.Handler(http.MethodGet, "/swagger/json", api.Handler())
	router.Handler(http.MethodGet, "/swagger/ui/*any", swag.UIHandler("/swagger/ui", "/swagger/json", true))

	// router.GET("/api/v1/tickets/:username", mid.AccessLog(HealthOK, app.Logger))
	// router.DELETE("/api/v1/tickets/:ticketUID", mid.AccessLog(HealthOK, app.Logger))

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	app.Router = router

	return app
}

func (app App) Stop() {
	app.DB.Close()
}
