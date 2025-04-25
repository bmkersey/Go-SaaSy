package auth

import "net/http"

func GetUserID(r *http.Request) (string, bool) {
	val := r.Context().Value(UserIDKey)
	id, ok := val.(string)
	return id, ok
}
