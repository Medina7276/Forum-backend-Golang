package middleware

import (
	"net/http"

	"git.01.alem.school/qjawko/forum/jwt"
)

func ContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func Unauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err == nil {
			http.Error(w, "You are already authorized!", http.StatusConflict)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token") //TODO: check user in db
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = jwt.IsValid(cookie.Value, "supersecret"); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func CheckForAdmin(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 		cookie, err := r.Cookie("token") //check if token in request
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		var user model.User
// 		if err := jwt.Unmarshal(cookie.Value, "supersecret", user); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 		}

// 	})
// }

// func CheckForAdmin(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("token") //check if user in ROLE
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		var user *model.User
// 		if err := jwt.Unmarshal(cookie.Value, "supersecret", user); err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		if user.Role < model.ROLE_ADMIN {
// 			http.Error(w, "NOT ADMIN", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func CheckForModerator(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("token")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		var user model.User

// 		if err := jwt.Unmarshal(cookie.Value, "supersecret", &user); err != nil {
// 			fmt.Println(err)
// 			return
// 		}

// 		if user.Role < model.ROLE_MODER {
// 			http.Error(w, "NOT MODER", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
