package session

import (
	"github.com/gorilla/sessions"
)


var Store = sessions.NewCookieStore(
	[]byte("top-s3cr3t"))