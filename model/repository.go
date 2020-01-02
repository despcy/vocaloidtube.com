package model

import (
	"errors"
	"log"
	"sync"

	"github.com/yangchenxi/VOCALOIDTube/model/authentication"

	"github.com/kataras/iris"

	"github.com/yangchenxi/VOCALOIDTube/model/dbops"
	"github.com/yangchenxi/VOCALOIDTube/model/defs"
	"go.mongodb.org/mongo-driver/bson"
)

var VisitCount uint64 = 0
var VisitCountMutex = &sync.Mutex{}

type UserDBService interface {
	GetUserCount() int
}

func NewUserDBService() UserDBService {
	return &userDBService{}
}

type userDBService struct {
}

func (s *userDBService) GetUserCount() int {
	return int(dbops.GetUserDBCount())

}

type VideoDBService interface {
	GetAllVideos(pageNum int, sortBy bson.M) (int, []*defs.Video)
	SearchVideosByID(vid string) ([]*defs.Video, error)
	GetVideoCount() int
}

func NewVideoDBService() VideoDBService {
	return &videoDBService{}
}

type videoDBService struct {
}

func (s *videoDBService) GetVideoCount() int {
	return int(dbops.GetVideoCount())
}

func (s *videoDBService) GetAllVideos(pageNum int, sortBy bson.M) (int, []*defs.Video) {
	return dbops.GetAllVideos(pageNum, sortBy)
}

func (s *videoDBService) SearchVideosByID(vid string) ([]*defs.Video, error) {
	count, Videos := dbops.QueryVideos(bson.M{"_id": vid}, 1, bson.M{"statistic.viewcount": -1})
	if count == 0 {
		return nil, errors.New("video id" + vid + " not found")
	}
	return Videos, nil
}

type AuthService interface {
	BeginAuth(ctx iris.Context)
	CallbackAuth(ctx iris.Context)
	Logout(ctx iris.Context)
}

func NewAuthService() AuthService {
	return &authService{}
}

type authService struct {
}

func (s *authService) BeginAuth(ctx iris.Context) {
	// try to get the user without re-authenticating
	if gothUser, err := authentication.CompleteUserAuth(ctx); err == nil {
		ctx.ViewData("", gothUser)
		if err := ctx.View("user.html"); err != nil {
			ctx.Writef("%v", err)
		}
	} else {
		authentication.BeginAuthHandler(ctx)
	}
}

func (s *authService) CallbackAuth(ctx iris.Context) {
	user, err := authentication.CompleteUserAuth(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Writef("%v", err)
		return
	}
	session := authentication.SessionsManager.Start(ctx)
	session.Set("UserEmail", user.Email)
	session.Set("UserAvatarURL", user.AvatarURL)
	session.Set("UserID", user.UserID)
	//Set session ID and csrf , store/update user profile data to maongodb ->session model
	// ctx.ViewData("", user)
	// if err := ctx.View("user.html"); err != nil {
	// 	ctx.Writef("%v", err)
	userdata := defs.UserProfile{
		AvatarURL: user.AvatarURL,
		UserEmail: user.Email,
		UserID:    user.UserID,
	}
	go dbops.UpdateUser(userdata)
	ctx.Redirect("/")

}

func (s *authService) Logout(ctx iris.Context) {

	err := authentication.Logout(ctx)
	if err != nil {
		log.Println(err)
	}
	ctx.Redirect("/")
	//invalidate cookie ->session model
}
