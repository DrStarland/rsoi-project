package initialization

import (
	"fmt"
	"gateway/pkg/handlers"
	mid "gateway/pkg/middleware"
	"gateway/pkg/models/authorization"
	"gateway/pkg/models/cost"
	"gateway/pkg/models/income"
	"gateway/pkg/models/note"
	"gateway/pkg/models/scope"
	"gateway/pkg/models/task"
	"gateway/pkg/myjson"

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
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	authHandler := handlers.NewAuthHandler(app.Logger)
	taskHandler := handlers.NewTaskHandler(app.Logger)
	noteHandler := handlers.NewNoteHandler(app.Logger)
	incomeHandler := handlers.NewIncomeHandler(app.Logger)
	costHandler := handlers.NewCostHandler(app.Logger)
	calcHandler := handlers.NewCalcHandler(app.Logger)

	kafka := InitKafka(app.Logger)

	api := swag.New(
		option.Title("Costs-n-tasks API Doc"),
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
	api.AddTag("Authorization", "")
	api.AddTag("Tasks", "")
	api.AddTag("Notes", "")
	api.AddTag("Costs", "")
	api.AddTag("Incomes", "")
	api.AddTag("Balance", "")

	api.AddEndpoint(
		// СЕРВИС
		endpoint.New(
			http.MethodGet, "/manage/health",
			endpoint.Handler(
				mid.AccessLog(HealthOK, app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary(""),
			endpoint.Response(http.StatusOK, "Server available"),
			endpoint.Tags("Healthcheck and statistics"),
		),
		// АВТОРИЗАЦИЯ
		endpoint.New(
			http.MethodPost, "/authorize",
			endpoint.Handler(
				mid.AccessLog(authHandler.Authorize, app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Авторизация пользователя"),
			endpoint.Body(authorization.AuthRequest{}, "Структура запроса на создание пользователя", true),

			endpoint.Response(http.StatusOK, "Registration was successful"),
			endpoint.Tags("Authorization"),
		),
		endpoint.New(
			http.MethodPost, "/register",
			endpoint.Handler(
				mid.AccessLog(authHandler.Register, app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Регистрация пользователя"),
			endpoint.Body(authorization.UserCreateRequest{}, "Структура запроса на создание пользователя", true),

			endpoint.Response(http.StatusOK, "Registration was successful"),
			endpoint.Tags("Authorization"),
		),
		// ЗАПИСКИ
		endpoint.New(
			http.MethodGet, "/notes",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.List, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Возвращает список заметок"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodGet, "/notes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Show, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(noteHandler.Add, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(noteHandler.Update, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Редактирование существующей заметки"),
			endpoint.Path("id", "integer", "ID of note to edit", true),
			endpoint.Body(note.NoteCreationRequest{},
				"Структура запроса на изменение заметки", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(note.Note{})),
			endpoint.Response(http.StatusBadRequest, ""),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodDelete, "/notes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(noteHandler.Delete, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Удаление заметки"),
			endpoint.Path("id", "integer", "ID of note to delete", true),
			endpoint.Response(http.StatusNoContent, "successful"),
			endpoint.Response(http.StatusNoContent, "Entity is not exist or already deleted", endpoint.SchemaResponseOption("not exist")),
			endpoint.Tags("Notes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		// ЗАДАЧИ
		endpoint.New(
			http.MethodGet, "/tasks",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.List, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Возвращает список заметок"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Tasks"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodGet, "/tasks/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Show, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(taskHandler.Add, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Создание новой задачи"),
			endpoint.Body(task.TaskCreationRequest{}, "Структура запроса на создание задачи", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(task.Task{})),
			endpoint.Tags("Tasks"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/tasks/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(taskHandler.Update, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(taskHandler.Delete, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(taskHandler.AddComment, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
		endpoint.New(
			http.MethodGet, "/costs",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.List, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Возвращает список расходов"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Costs"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodGet, "/gateway/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Show, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(costHandler.Add, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Создание новой записи о расходе"),
			endpoint.Body(cost.CostCreationRequest{}, "Структура запроса на создание записи о расходе", true),
			endpoint.Response(http.StatusOK, "was successful", endpoint.SchemaResponseOption(cost.Cost{})),
			endpoint.Tags("Costs"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodPut, "/gateway/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Update, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
			http.MethodDelete, "/gateway/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(costHandler.Delete, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(incomeHandler.List, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Возвращает список доходов"),
			endpoint.Response(http.StatusOK, ""),
			endpoint.Tags("Incomes"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
		endpoint.New(
			http.MethodGet, "/incomes/{id}",
			endpoint.Handler(
				mid.AccessLog(mid.Auth(incomeHandler.Show, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(incomeHandler.Add, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(incomeHandler.Update, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(incomeHandler.Delete, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
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
				mid.AccessLog(mid.Auth(calcHandler.TotalBalance, app.Logger), app.Logger, kafka.Topic, kafka.Producer),
			),
			endpoint.Summary("Подсчёт баланса. Метод реализован как POST для того, чтобы у запроса могло быть тело."),
			endpoint.Body(scope.Scope{}, "Для какой области видимости вычислить баланс", true),
			endpoint.Response(http.StatusOK, "Balance result", endpoint.SchemaResponseOption(myjson.ResponceForm{})),
			endpoint.Tags("Balance"),
			endpoint.Security("Sophisticated_Service_auth", "read:pets"),
		),
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
