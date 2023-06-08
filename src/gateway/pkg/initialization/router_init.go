package initialization

import (
	"fmt"
	"gateway/pkg/handlers"
	mid "gateway/pkg/middleware"
	"gateway/pkg/models/authorization"

	// "users/pkg/models/authorization"

	// "gateway/pkg/models/authorization"

	"gateway/pkg/utils"
	"log"
	"net/http"
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

func PlugToDo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app App) initRouter() App {

	// repoTicket := ticket.NewPostgresRepo(db)
	// tak := dbcontext.New(app.DB)
	// repoNotes := note.NewRepository(tak, app.Logger)

	// noteHandler := &handlers.NoteMainHandler{
	// 	Logger:  app.Logger,
	// 	Repo:    repoNotes,
	// 	Service: handlers.NewNoteService(repoNotes, *app.Logger),
	// }

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	// client := &http.Client{}

	// auth_var := models.NewAuthM(client)
	// ctrl := &auth.AuthCtrl{Auth: auth_var}

	authHandler := handlers.NewAuthHandler(app.Logger)

	// 	ctrl := &AuthCtrl{auth}
	// 	r.HandleFunc("/register", ctrl.Register).Methods("POST")
	// 	r.HandleFunc("/authorize", ctrl.Authorize).Methods("POST")

	// router.GET("/manage/health", ri\outer\HealthOK)

	// type Category struct {
	// 	ID     int64  `json:"category"`
	// 	Name   string `json:"name" enum:"dog,cat" required:""`
	// 	Exists *bool  `json:"exists" required:""`
	// }

	// // Pet example from the swagger pet store
	// type Pet struct {
	// 	ID        int64     `json:"id"`
	// 	Category  *Category `json:"category" desc:"分类"`
	// 	Name      string    `json:"name" required:"" example:"张三" desc:"名称"`
	// 	PhotoUrls []string  `json:"photoUrls"`
	// 	Tags      []string  `json:"tags" desc:"标签"`
	// }

	// handle := func(w http.ResponseWriter, r *http.Request) {
	// 	_, _ = io.WriteString(w, fmt.Sprintf("[%s] Hello World!", r.Method))
	// }

	api := swag.New(
		option.Title("Costs-n-tasks API Doc"),
		option.Security("user_auth", "read:pets"),
		option.SecurityScheme("user_auth",
			option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			option.OAuth2Scope("admin", ""),
			option.OAuth2Scope("user", ""),
		),
		option.BasePath("/api/v1"),
	)

	api.AddTag("Healthcheck and statistics", "")
	api.AddTag("Authorization", "")
	api.AddTag("Tasks", "")

	api.AddEndpoint(
		// СЕРВИС
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(mid.AccessLog(HealthOK, app.Logger)),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),
		// АВТОРИЗАЦИЯ
		endpoint.New(
			http.MethodPost, "/authorize",
			endpoint.Handler(authHandler.Authorize),
			endpoint.Summary("Авторизация пользователя"),
			endpoint.Body(authorization.AuthRequest{}, "Структура запроса на создание пользователя", true),

			endpoint.Response(http.StatusOK, "Registration was successful"),
			endpoint.Tags("Authorization"),
		),
		endpoint.New(
			http.MethodPost, "/register",
			endpoint.Handler(authHandler.Register),
			endpoint.Summary("Регистрация пользователя"),
			endpoint.Body(authorization.UserCreateRequest{}, "Структура запроса на создание пользователя", true),

			endpoint.Response(http.StatusOK, "Registration was successful"),
			endpoint.Tags("Authorization"),
		),
		// ЗАПИСКИ
		// ЗАДАЧИ
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(mid.AccessLog(HealthOK, app.Logger)),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),

		endpoint.New(
			http.MethodGet, "/tasks",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.List, app.Logger), app.Logger),
			),
			endpoint.Summary("Возвращает список заметок"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Tasks"),
		),
		endpoint.New(
			http.MethodGet, "/tasks/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Show, app.Logger), app.Logger),
			),
			endpoint.Summary("Получение задачи по ID"),
			endpoint.Tags("Tasks"),
			endpoint.Path("id", "integer", "ID of task to return", true),
			endpoint.Response(http.StatusOK, "successful operation", endpoint.SchemaResponseOption(task.Task{})),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPost, "/tasks",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Add, app.Logger), app.Logger),
			),
			endpoint.Summary("Создание новой задачи"),
			endpoint.Body(task.TaskCreationRequest{}, "Структура запроса на создание задачи", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(task.Task{})),
			endpoint.Tags("Tasks"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/tasks/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Update, app.Logger), app.Logger),
			),
			endpoint.Summary("Редактирование существующей задачи"),
			endpoint.Path("id", "integer", "ID of task to edit", true),
			endpoint.Body(task.TaskCreationRequest{},
				"Структура запроса на изменение задачи", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(task.Task{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Tasks"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodDelete, "/tasks/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Delete, app.Logger), app.Logger),
			),
			endpoint.Summary("Удаление задачи"),
			endpoint.Path("id", "integer", "ID of task to delete", true),
			endpoint.Response(http.StatusNoContent, "successful"),
			endpoint.Response(http.StatusNoContent, "Entity is not exist or already deleted", endpoint.SchemaResponseOption("not exist")),
			endpoint.Tags("Tasks"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPost, "/tasks/{id}/comments",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.AddComment, app.Logger), app.Logger),
			),
			endpoint.Summary("Добавление комментария к существующей задаче"),
			endpoint.Path("id", "integer", "ID задачи для добавления комментария", true),
			endpoint.Body(task.CommentCreationRequest{},
				"Структура запроса на создание комментария", true),
			endpoint.Response(http.StatusOK, "Успешное выполнение", endpoint.SchemaResponseOption(task.Task{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Comments"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		// ЗАПИСИ ДОХОДОВ И РАСХОДОВ
		// СТАТИСТИКА
	)

	swag.New()

	api.Walk(func(path string, e *swag.Endpoint) {
		h := e.Handler.(http.Handler)
		path = swag.ColonPath(path)
		router.Handler(e.Method, path, h)
	})

	router.Handler(http.MethodGet, "/swagger/json", api.Handler())
	router.Handler(http.MethodGet, "/swagger/ui/*any", swag.UIHandler("/swagger/ui", "/swagger/json", true))

	app.Router = router
	utils.InitConfig()

	return app
}

func (app App) Stop() {
	// app.DB.Close()
}
