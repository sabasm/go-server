package handlers

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
