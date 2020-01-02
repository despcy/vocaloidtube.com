package model

import (
	"strconv"

	"github.com/yangchenxi/VOCALOIDTube/model/defs"

	"github.com/kataras/iris"
	"github.com/yangchenxi/VOCALOIDTube/model/dbops"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllVideosInDB(ctx iris.Context) {
	sortBy := ctx.URLParamDefault("sortBy", "click")
	page := ctx.URLParamDefault("page", "1")
	//TODO:Illegal param handling
	pageNum, _ := strconv.Atoi(page)

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
	count, Videos := dbops.GetAllVideos(pageNum, bson.M{sortByBson: -1})

	ctx.JSON(defs.VideosResponse{
		PageNum:   pageNum,
		TotalPage: count,
		Sort:      sortBy,
		Videos:    Videos})

}

func searchVideos(ctx iris.Context) {
	q := ctx.URLParamDefault("q", "")
	searchType := ctx.URLParamDefault("type", "video")
	sortBy := ctx.URLParamDefault("sortBy", "click")
	page := ctx.URLParamDefault("page", "1")
	//TODO:Illegal param handling
	pageNum, _ := strconv.Atoi(page)
	//TODO:sql injection check,去除query中reg字符等
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

	ctx.JSON(defs.VideosResponse{
		PageNum:   pageNum,
		Sort:      sortBy,
		TotalPage: count,
		Videos:    Videos})

}
