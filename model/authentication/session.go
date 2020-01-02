package authentication

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/sessions"
)

var SessionsManager *sessions.Sessions

func NewSessionsManager() *sessions.Sessions {
	// attach a session manager
	cookieName := "mycustomsessionid"
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	SessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
	//这里只是给client发送一个session id，然后具体的session映射存储是放在本地内存的
	return SessionsManager
}
