package controller

import (
	"github.com/kataras/iris"
	"github.com/yangchenxi/VOCALOIDTube/model"
)

type AuthController struct {
	Service model.AuthService
}

func (c *AuthController) GetGoogle(ctx iris.Context) {
	defer ctx.Next()
	c.Service.BeginAuth(ctx)
}

func (c *AuthController) GetCallback(ctx iris.Context) {
	defer ctx.Next()
	c.Service.CallbackAuth(ctx)
}

func (c *AuthController) GetLogout(ctx iris.Context) {
	defer ctx.Next()
	c.Service.Logout(ctx)

}
