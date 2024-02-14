package AuthenticationService

import (
	"net/http"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"

	"github.com/gofrs/uuid"
)

func GetUserFromCookie(r *http.Request) *models.User {
	var User *models.User
	if cookie, err := r.Cookie("session"); err == nil {
		User, _ = dbfunctions.GetUserBySessionId(cookie.Value)
	}
	return User
}

func GenerateCookies(user models.User) (*http.Cookie, error) {
	sessionId, _ := uuid.NewV4()
	expiry := time.Now().Add(96 * time.Hour)
	formattedExpiry := expiry.Format("2006-01-02 15:04:05")
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionId.String(),
		Expires:  expiry,
		HttpOnly: true,
		Secure:   false,
	}
	dbfunctions.RemoveCookieData(user.Id)
	var session models.Session
	session.AccountId = user.Id
	session.Expiry = formattedExpiry
	session.Id = sessionId.String()
	if err := dbfunctions.InsertCookieData(session); err != nil {
		return nil, err
	}
	return cookie, nil
}
