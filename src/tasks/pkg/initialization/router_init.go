package initialization

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"tasks/pkg/handlers"
	mid "tasks/pkg/middleware"
	"tasks/pkg/models/task"
	"tasks/pkg/utils"

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
	repoTasks := task.NewMemoryRepository(app.Logger)

	taskHandler := &handlers.TaskMainHandler{
		Logger:  app.Logger,
		Repo:    repoTasks,
		Service: handlers.NewTaskService(repoTasks, *app.Logger),
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
		option.Title("Tasks Service API Doc"),
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
	api.AddTag("Tasks", "")

	api.AddEndpoint(
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
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
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
			endpoint.Response(http.StatusCreated, "was successful", endpoint.SchemaResponseOption(task.Task{})),
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
		// endpoint.New(
		// 	http.MethodPost, "/register",
		// 	endpoint.Handler(authHandler.Register),
		// 	//endpoint.HeaderResponseOption()
		// 	endpoint.Summary("Регистрация пользователя"),
		// 	endpoint.Body(authorization.UserCreateRequest{}, "Структура запроса на создание пользователя", true),

		// 	endpoint.Response(http.StatusOK, "Registration was successful"),
		// 	endpoint.Tags("Authorization"),
		// ),

		// endpoint.New(
		// 	http.MethodPost, "/register",
		// 	endpoint.Handler(authHandler.Register),
		// 	//endpoint.HeaderResponseOption()
		// 	endpoint.Summary("Регистрация пользователя"),
		// 	endpoint.Body(authorization.UserCreateRequest{}, "Структура запроса на создание пользователя", true),

		// 	endpoint.Response(http.StatusOK, "Registration was successful"),
		// 	endpoint.Tags("Authorization"),
		// ),
	// endpoint.New(
	// 	http.MethodPost, "/pet",
	// 	endpoint.Handler(handle),
	// 	endpoint.Summary("Add a new pet to the store"),
	// 	endpoint.Description("Additional information on adding a pet to the store"),
	// 	endpoint.Body(Pet{}, "Pet object that needs to be added to the store", true),
	// 	endpoint.Response(http.StatusOK, "Successfully added pet", endpoint.SchemaResponseOption(Pet{})), //End Schema(P)),
	// 	endpoint.Security("petstore_auth", "read:pets", "write:pets"),
	// 	endpoint.Tags("section"),
	// ),
	// endpoint.New(
	// 	http.MethodGet, "/pet/{petId}",
	// 	endpoint.Handler(handle),
	// 	endpoint.Summary("Find pet by ID"),
	// 	endpoint.Path("petId", "integer", "ID of pet to return", true),
	// 	endpoint.Response(http.StatusOK, "successful operation", endpoint.SchemaResponseOption(Pet{})),
	// 	endpoint.Security("petstore_auth", "read:pets"),
	// ),
	// endpoint.New(
	// 	http.MethodPut, "/pet/{petId}",
	// 	endpoint.Handler(handle),
	// 	endpoint.Path("petId", "integer", "ID of pet to return", true),
	// 	endpoint.Security("petstore_auth", "read:pets"),
	// 	endpoint.ResponseSuccess(endpoint.SchemaResponseOption(struct {
	// 		ID   string `json:"id"`
	// 		Name string `json:"name"`
	// 	}{})),
	// ),
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
