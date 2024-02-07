package forum

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie := &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, usernameCookie)
	http.Redirect(w, r, "/login", http.StatusFound)
}
