package controller

import (
	"log"
	"strconv"

	"github.com/yangchenxi/VOCALOIDTube/security"

	"github.com/yangchenxi/VOCALOIDTube/model/dbops"
	viewrender "github.com/yangchenxi/VOCALOIDTube/view"

	"github.com/kataras/iris"
	"github.com/yangchenxi/VOCALOIDTube/model"
)

type StatusController struct {
	UserDBService model.UserDBService
}

func (c *StatusController) Get(ctx iris.Context) {
	if c.UserDBService == nil {
		log.Println("null userDBService")

	}

	v := make([]security.BlockedIP, 0, len(security.BlackList))
	for _, value := range security.BlackList {
		v = append(v, value)
	}

	userNum := strconv.Itoa(c.UserDBService.GetUserCount())
	videoNum := strconv.Itoa(int(dbops.GetVideoCount())) //Here is not format and standard
	viewrender.RenderStatusPage(ctx, userNum, videoNum, strconv.FormatUint(model.VisitCount, 10), v)
	defer ctx.Next()

}
