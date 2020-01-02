package controller

import (
	"github.com/yangchenxi/VOCALOIDTube/security"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"

	"github.com/yangchenxi/VOCALOIDTube/model"
	"github.com/yangchenxi/VOCALOIDTube/model/dbops"
	"github.com/yangchenxi/VOCALOIDTube/model/defs"
	viewrender "github.com/yangchenxi/VOCALOIDTube/view"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoController struct {
	VideoDB model.VideoDBService
}

func (c *VideoController) GetSession(ctx iris.Context) {
	defer ctx.Next()
	session := sessions.Get(ctx)
	if session.Len() == 0 {
		ctx.HTML(`no session values stored yet. Navigate to: <a href="/set">set page</a>`)
		return
	}
	ctx.HTML(`<img src="` + session.GetString("UserAvatarURL") + `" ></img>`)

}

func (c *VideoController) Get(ctx iris.Context) {
	c.GetVideos(ctx)

}

func (c *VideoController) GetVsingers(ctx iris.Context) {
	defer ctx.Next()
	viewrender.RenderVsingerPage(ctx)

}

func (c *VideoController) GetVideosBy(ctx iris.Context, vid string) {
	//TODO: filter vid
	defer ctx.Next()

	if len(vid) > 15 {
		ctx.StatusCode(403)
		security.AddIPToBlackList(ctx.RemoteAddr(), "Invalied Video request parameter, seems to be an attack. 打入冷宫！")
		return
	}
	videos, err := c.VideoDB.SearchVideosByID(vid)
	//if vid not in db,return error error这里返回一个彩蛋
	if err != nil {
		ctx.StatusCode(404)
		ctx.View("page-404.html")
		return
	}

	videoData := defs.VideosResponse{
		Videos: videos}
	viewrender.RenderPlayPage(viewrender.HomePageViewModel{
		Ctx:       ctx,
		VideoResp: &videoData})

}

func (c *VideoController) GetVideos(ctx iris.Context) {
	defer ctx.Next()
	sortBy := ctx.URLParamDefault("sortBy", "click")
	//TODO:Illegal param handling
	pageNum, pageerr := ctx.URLParamInt("page")
	if pageerr != nil {
		pageNum = 1
	}
	if !ControllerSecurityCheck(sortBy) {
		ctx.StatusCode(403)
		security.AddIPToBlackList(ctx.RemoteAddr(), "Try to inject Database on home page!")
		return
	}
	var sortByBson string
	if sortBy == "click" || sortBy == "like" || sortBy == "recent" {
		if sortBy == "click" {
			sortByBson = "statistic.viewcount"
		} else if sortBy == "like" {
			//like
			sortByBson = "statistic.likecount"
		} else if sortBy == "recent" {
			sortByBson = "snippet.publishedat"
		}
	} else {
		// ctx.StatusCode(iris.StatusBadRequest)
		// ctx.JSON(defs.ErrorRequestBodyParseFailed)
		return
	}

	count, Videos := c.VideoDB.GetAllVideos(pageNum, bson.M{sortByBson: -1})

	videoData := defs.VideosResponse{
		PageNum:   pageNum,
		Sort:      sortBy,
		TotalPage: count,
		Videos:    Videos}
	viewrender.RenderHomePage(viewrender.HomePageViewModel{
		Ctx:       ctx,
		VideoResp: &videoData})

}

func (c *VideoController) GetSearch(ctx iris.Context) {
	defer ctx.Next()
	q := ctx.URLParamDefault("q", "")
	searchType := ctx.URLParamDefault("type", "video")
	sortBy := ctx.URLParamDefault("sortBy", "click")
	pageNum, pageerr := ctx.URLParamInt("page")
	if pageerr != nil {
		pageNum = 1
	}
	if !ControllerSecurityCheck(q, searchType, sortBy) {
		ctx.StatusCode(403)
		security.AddIPToBlackList(ctx.RemoteAddr(), "Try to inject Database on search page!")
		return
	}

	var sortByBson string
	if sortBy == "click" || sortBy == "like" || sortBy == "recent" {
		if sortBy == "click" {
			sortByBson = "statistic.viewcount"
		} else if sortBy == "like" {
			//like
			sortByBson = "statistic.likecount"
		} else if sortBy == "recent" {
			sortByBson = "snippet.publishedat"
		}
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(defs.ErrorRequestBodyParseFailed)
		return
	}

	var query bson.M
	if searchType == "video" {
		query = bson.M{"snippet.title": primitive.Regex{Pattern: ".*" + q + "*.", Options: "i"}}
	} else if searchType == "studio" {
		query = bson.M{"snippet.channeltitle": primitive.Regex{Pattern: ".*" + q + "*.", Options: "i"}}
	} else if searchType == "vsinger" {
		query = bson.M{"vsinger": primitive.Regex{Pattern: q, Options: "i"}}
	} else if searchType == "tag" {
		query = bson.M{"snippet.tags": primitive.Regex{Pattern: q, Options: "i"}}
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(defs.ErrorRequestBodyParseFailed)
		return
	}
	count, Videos := dbops.QueryVideos(query, pageNum, bson.M{sortByBson: -1})

	videoData := defs.VideosResponse{
		PageNum:   pageNum,
		Sort:      sortBy,
		TotalPage: count,
		Videos:    Videos}
	videoSearchData := viewrender.SearchInfo{
		QueryParam: q,
		Qtype:      searchType,
	}
	viewrender.RenderSearchPage(viewrender.HomePageViewModel{
		Ctx:        ctx,
		VideoResp:  &videoData,
		SearchData: videoSearchData})

}

func ControllerSecurityCheck(params ...string) bool {

	for _, para := range params {
		if !security.DBInjectionKeywordCheck(para) {
			return false
		}
	}

	return true
}

//Honeypot :)
func (c *VideoController) GetAdmin(ctx iris.Context) {
	ctx.StatusCode(403)
	security.AddIPToBlackList(ctx.RemoteAddr(), "Triggered the honeypot, boom!")
	ctx.Next()
}
