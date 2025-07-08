package website_handler

import "github.com/gin-gonic/gin"

func DashboardPage(c *gin.Context) {
	// cookie, err := r.Cookie("token")
	// if err != nil {
	// 	// Cookie not found or other error
	// 	log.Println("Token cookie not found:", err)
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	// // Maybe I can use redis to save the token so that I can retrieve the user information without the need to validate the cookie
	// token, err := auth.ValidateJWT(cookie.Value)
	// if err != nil {
	// 	log.Println("Invalid token:", err)
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	// claims := token.Claims.(jwt.MapClaims)
	// username, ok := claims["username"].(string)
	// if !ok {
	// 	log.Println("Username claim not found or not a string")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	// password, ok := claims["password"].(string)
	// if !ok {
	// 	log.Println("Password claim not found or not a string")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	// log.Println("Current username:", username)
	// log.Println("Current password:", password)
	// templates.ExecuteTemplate(w, "dashboard.html", nil)
}
