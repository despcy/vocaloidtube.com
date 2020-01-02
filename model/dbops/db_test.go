package dbops

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/yangchenxi/VOCALOIDTube/crawler"

	"github.com/yangchenxi/VOCALOIDTube/model/defs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMain(m *testing.M) {

	m.Run()

}

func TestInsertVsingerToExistingData(t *testing.T) {

	count, _ := GetAllVideos(1, bson.M{"statistic.likecount": -1})
	log.Println(count)
	for i := 1; i <= count; i++ {
		_, Videos := GetAllVideos(i, bson.M{"statistic.likecount": -1})
		for _, video := range Videos {
			if video.Vsinger == "" {
				video.Vsinger = crawler.ClassifyVsinger(video.Snippet.Tags)
			}
			if video.Parent == "" {
				video.Parent = "Unknown"
			}
			UpdateVideo(video)
			log.Println(video.Snippet.Title + " " + video.Vsinger)
		}
	}
}

func TestInsertVideosInAccount(t *testing.T) {

}

func TestQueryOne(t *testing.T) {
	count, Videos := QueryVideos(bson.M{"_id": "ARt2fVT33Lw"}, 1, bson.M{"statistic.viewcount": -1})
	log.Println(count)
	log.Println(len(Videos))
}

func TestConn(t *testing.T) {
	// Check the connection
	err = dbClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func TestInsert(t *testing.T) {
	InsertVideo("ypYxB3H2UQI", "pid", "Mo Qingxian")

}

func TestInsertUser(t *testing.T) {
	userdata := defs.UserProfile{
		AvatarURL: "testUwRL",
		UserEmail: "test@test.com",
		UserID:    "123",
	}
	AddUser(userdata)
}

func TestUpdateUser(t *testing.T) {
	userdata := defs.UserProfile{
		AvatarURL: "222111dfffftestUwRL",
		UserEmail: "aaa222111test@test.com",
		UserID:    "123ggg",
	}
	//p := bson.M{"$set": bson.M{"avatarurl": userdata.AvatarURL}}
	UpdateUser(userdata)
}

func TestInsertVideosInPlayList(t *testing.T) {
	//PL_vdmUqgn18VVtdSieJfJw1ce9tB4Xppv
	//PLB02wINShjkBKnLfufaEPnCupGO-SK6e4->1.5K
	//PLXHs4FGKMGGzAO0fAk4a3Yb5ktwFQ7pMZ->1k
	//PL7ER-GcyaxAJILUYtt77ijfg8S-WgAG6_
	erro := InsertVideosInPlayList("RDqaQvtckEJB0", "")
	if erro != nil {
		t.Errorf("Error of : %v", err)

	}
	time.Sleep(100 * time.Second)
}

func TestGetAllVideos(t *testing.T) {
	count, Videos := GetAllVideos(2, bson.M{"statistic.likecount": -1})
	log.Println(count)
	for _, video := range Videos {
		log.Println(video.VideoID)
	}

	log.Println(len(Videos))
}

func TestQueryDB(t *testing.T) {
	//video
	//testQueryDB(bson.M{"snippet.title": primitive.Regex{Pattern: ".*miku*.", Options: "i"}}, 1, bson.M{"statistic.viewcount": -1})
	//studio
	//testQueryDB(bson.M{"snippet.channeltitle": primitive.Regex{Pattern: ".*official*.", Options: "i"}}, 1, bson.M{"statistic.viewcount": -1})
	//vsinger(not specified)
	//TODO vsinger query test
	//tag
	testQueryDB(bson.M{"snippet.tags": primitive.Regex{Pattern: "vocaloid", Options: "i"}}, 1, bson.M{"statistic.viewcount": -1})
}

func testQueryDB(query bson.M, pageNum int, sortBy bson.M) {
	count, Videos := QueryVideos(query, pageNum, sortBy)
	log.Println(count)
	for _, video := range Videos {
		log.Println(video.Snippet.Title)
		log.Println(video.Snippet.ChannelTitle)
		log.Println(video.Vsinger)
		log.Println(video.Snippet.Tags)
		log.Println("-------------------")
	}

	log.Println(len(Videos))

}

func TestVsingerTag(t *testing.T) {
	log.Println(GetVsingerVideoTags("Ken"))

}
