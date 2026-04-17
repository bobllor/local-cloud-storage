package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bobllor/assert"
	"github.com/bobllor/cloud-project/src/session"
)

func TestGetCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", bytes.NewBuffer([]byte{}))
	session := "12345-session-id-here"

	cookie := http.Cookie{
		Name:  CookieSessionKey,
		Value: session,
	}

	r.AddCookie(&cookie)
	v := GetSessionFromCookie(r)

	assert.Equal(t, v, session)
}

func TestSetCookie(t *testing.T) {
	w := httptest.NewRecorder()

	s := session.Session{
		SessionID: "12345-session-id-here",
		AccountID: "12345-account-id",
		CreatedOn: time.Now().UTC(),
		ExpireOn:  time.Now().UTC(),
	}

	SetCookieSession(w, &s)

	cookies := w.Result().Cookies()

	assert.Equal(t, len(cookies), 1)
}
