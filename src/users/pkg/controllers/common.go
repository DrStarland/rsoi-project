package controllers

// func initControllers(r *mux.Router, models *models.Models) {
// 	r.Use(utils.LogHandler)
// 	api1_r := r.PathPrefix("/api/v1/").Subrouter()

// 	InitAuth(api1_r, models.Auth)
// }

// func InitRouter() *mux.Router {
// 	router := mux.NewRouter()
// 	models := models.InitModels()

// 	initControllers(router, models)
// 	return router
// }

// func RunRouter(r *mux.Router, port uint16) error {
// 	c := cors.New(cors.Options{})
// 	handler := c.Handler(r)
// 	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
// }
