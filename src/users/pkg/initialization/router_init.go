package initialization

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"users/pkg/handlers"
	"users/pkg/models/authorization"
	"users/pkg/utils"

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
		option.Title("User Service API Doc"),
		option.Security("user_auth", "read:pets"),
		option.SecurityScheme("user_auth",
			option.OAuth2Security("accessCode", "http://example.com/oauth/authorize", "http://example.com/oauth/token"),
			option.OAuth2Scope("admin", ""),
			option.OAuth2Scope("user", ""),
		),
		option.BasePath("/api/v1"),
	)

	api.AddTag("Authorization", "")
	api.AddTag("Healthcheck and statistics", "")
	api.AddTag("Authorization", "")

	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(HealthOK),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),

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
			//endpoint.HeaderResponseOption()
			endpoint.Summary("Регистрация пользователя"),
			endpoint.Body(authorization.UserCreateRequest{}, "Структура запроса на создание пользователя", true),

			endpoint.Response(http.StatusOK, "Registration was successful"),
			endpoint.Tags("Authorization"),
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

	app.Router = router

	utils.InitConfig()
	return app
}

func (app App) Stop() {
	// app.DB.Close()
}
