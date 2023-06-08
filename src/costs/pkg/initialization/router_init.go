package initialization

import (
	"costs/pkg/handlers"
	mid "costs/pkg/middleware"
	"costs/pkg/models/cost"
	"costs/pkg/models/income"
	"costs/pkg/models/scope"
	"costs/pkg/myjson"
	"costs/pkg/services"
	"costs/pkg/utils"
	"fmt"
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

func (app App) initRouter() App {

	// repoTicket := ticket.NewPostgresRepo(db)
	// tak := dbcontext.New(app.DB)
	repoCosts := cost.NewMemoryRepository(app.Logger)
	repoIncomes := income.NewMemoryRepository(app.Logger)

	costHandler := &handlers.CostMainHandler{
		Logger:  app.Logger,
		Repo:    repoCosts,
		Service: services.NewCostService(repoCosts, *app.Logger),
	}

	incomeHandler := &handlers.IncomeMainHandler{
		Logger:  app.Logger,
		Repo:    repoIncomes,
		Service: services.NewIncomeService(repoIncomes, *app.Logger),
	}

	calcHandler := &handlers.CalcMainHandler{
		Logger:        app.Logger,
		CostsSource:   costHandler.Service,
		IncomesSource: incomeHandler.Service,
	}

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
		option.Title("Costs Service API Doc"),
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
	api.AddTag("Costs", "")
	api.AddTag("Incomes", "")
	api.AddTag("Balance", "")

	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(mid.AccessLog(HealthOK, app.Logger)),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),

		endpoint.New(
			http.MethodGet, "/costs",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.List, app.Logger), app.Logger),
			),
			endpoint.Summary("Возвращает список расходов"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Costs"),
		),
		endpoint.New(
			http.MethodGet, "/costs/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Show, app.Logger), app.Logger),
			),
			endpoint.Summary("Получение записи о расходе по ID"),
			endpoint.Tags("Costs"),
			endpoint.Path("id", "integer", "ID of cost to return", true),
			endpoint.Response(http.StatusOK, "successful operation", endpoint.SchemaResponseOption(cost.Cost{})),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPost, "/costs",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Add, app.Logger), app.Logger),
			),
			endpoint.Summary("Создание новой записи о расходе"),
			endpoint.Body(cost.CostCreationRequest{}, "Структура запроса на создание записи о расходе", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(cost.Cost{})),
			endpoint.Tags("Costs"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/costs/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Update, app.Logger), app.Logger),
			),
			endpoint.Summary("Редактирование существующей записи о расходе"),
			endpoint.Path("id", "integer", "ID of cost to edit", true),
			endpoint.Body(cost.CostCreationRequest{},
				"Структура запроса на изменение записи о расходе", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(cost.Cost{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Costs"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodDelete, "/costs/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Delete, app.Logger), app.Logger),
			),
			endpoint.Summary("Удаление записи о расходе"),
			endpoint.Path("id", "integer", "ID of cost to delete", true),
			endpoint.Response(http.StatusNoContent, "successful"),
			endpoint.Response(http.StatusNoContent, "Entity is not exist or already deleted", endpoint.SchemaResponseOption("not exist")),
			endpoint.Tags("Costs"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		/////
		endpoint.New(
			http.MethodGet, "/incomes",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.List, app.Logger), app.Logger),
			),
			endpoint.Summary("Возвращает список доходов"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Incomes"),
		),
		endpoint.New(
			http.MethodGet, "/incomes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.Show, app.Logger), app.Logger),
			),
			endpoint.Summary("Получение записи о доходе по ID"),
			endpoint.Tags("Incomes"),
			endpoint.Path("id", "integer", "ID of income to return", true),
			endpoint.Response(http.StatusOK, "successful operation", endpoint.SchemaResponseOption(income.Income{})),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPost, "/incomes",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.Add, app.Logger), app.Logger),
			),
			endpoint.Summary("Создание новой записи о доходе"),
			endpoint.Body(income.IncomeCreationRequest{}, "Структура запроса на создание записи о доходе", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(income.Income{})),
			endpoint.Tags("Incomes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/incomes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.Update, app.Logger), app.Logger),
			),
			endpoint.Summary("Редактирование существующей записи о доходе"),
			endpoint.Path("id", "integer", "ID of income to edit", true),
			endpoint.Body(income.IncomeCreationRequest{},
				"Структура запроса на изменение записи о доходе", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(income.Income{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Incomes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodDelete, "/incomes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.Delete, app.Logger), app.Logger),
			),
			endpoint.Summary("Удаление записи о доходе"),
			endpoint.Path("id", "integer", "ID of income to delete", true),
			endpoint.Response(http.StatusNoContent, "successful"),
			endpoint.Response(http.StatusNoContent, "Entity is not exist or already deleted", endpoint.SchemaResponseOption("not exist")),
			endpoint.Tags("Incomes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),

		endpoint.New(
			http.MethodPost, "/balance",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(calcHandler.TotalBalance, app.Logger), app.Logger),
			),
			endpoint.Summary("Подсчёт баланса. Метод реализован как POST для того, чтобы у запроса могло быть тело."),
			endpoint.Body(scope.Scope{}, "Для какой области видимости вычислить баланс", true),
			endpoint.Response(http.StatusOK, "Balance result", endpoint.SchemaResponseOption(myjson.ResponceForm{})),
			endpoint.Tags("Balance"),
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
