package authentication

import (
	"errors"

	"github.com/yangchenxi/VOCALOIDTube/config"

	"github.com/kataras/iris"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

var providerName = "google"

func init() {

	goth.UseProviders(google.New(config.GoogleOauthClientID, config.GoogleOauthClientKey, config.SERVER_DOMAIN+"/auth/callback"))
	//Session 在这里的作用是判断所有的
}
func BeginAuthHandler(ctx iris.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}

	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}
	//TODO:set state and check state to prevent csrf
	return "woshistate"
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx iris.Context) string {
	return ctx.URLParam("state")
}

func GetAuthURL(ctx iris.Context) (string, error) {

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session := SessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

// Logout invalidates a user session.
func Logout(ctx iris.Context) error {

	session := SessionsManager.Start(ctx)
	session.Clear() //remove all sessions
	return nil
}

var CompleteUserAuth = func(ctx iris.Context) (goth.User, error) {

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := SessionsManager.Start(ctx)
	value := session.GetString(providerName)
	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Set(providerName, sess.Marshal())
	//这里session存的是accesstoken
	return provider.FetchUser(sess)
}
