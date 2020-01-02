package dbops

import (
	"context"
	"errors"
	"log"
	"math"

	"github.com/yangchenxi/VOCALOIDTube/crawler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yangchenxi/VOCALOIDTube/model/defs"
	"github.com/yangchenxi/VOCALOIDTube/model/youtubeData"
)

const ITEMS_PER_PAGE = 36 //maximum 50 according to youtube api
//InsertVideo inserts a video to db vid pid vsinger
func InsertVideo(vid string, pid string, vsinger string) error {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	searchVid, _ := collection.Find(context.TODO(), bson.M{"_id": vid})
	if searchVid.Next(context.TODO()) {
		return nil
	}
	videoInfo, err := youtubeData.GetVideoInfoFromID(vid)

	if err != nil {
		log.Println("err geting youtube video data")
		return err

	}
	if videoInfo == nil {
		return errors.New("no such video:" + vid)
	}
	if vsinger == "" {
		vsinger = crawler.ClassifyVsinger(videoInfo.Snippet.Tags)
	}
	v := &defs.Video{
		VideoID:          videoInfo.Id,
		Snippet:          *videoInfo.Snippet,
		Statistic:        *videoInfo.Statistics,
		VideoKind:        "",
		SiteViewCount:    0,
		SitefavCount:     0,
		Sitedislikecount: 0,
		Parent:           pid,
		Vsinger:          vsinger}

	result, err := collection.InsertOne(context.TODO(), v)
	log.Println(result)
	if err != nil {
		log.Println("Video Insert DB ERROR" + err.Error())
	}
	return nil
}

//InsertVideosInPlayList inserts all videos in playlist
func InsertVideosInPlayList(pid string, vsinger string) error {
	//var videos []*youtube.Video
	nextPageToken := ""
	for {
		resp, err := youtubeData.GetPlayListItemsFromPlayListID(pid, nextPageToken)
		if err != nil {
			log.Printf("video playlist get fail")
			return err

		}
		if resp.HTTPStatusCode != 200 {
			log.Println("Error HTTP CODE" + string(resp.HTTPStatusCode))

			return errors.New("ERROR HTTP CODE")
		}

		for _, item := range resp.Items {
			log.Printf(item.Snippet.ResourceId.VideoId)
			go InsertVideo(item.Snippet.ResourceId.VideoId, pid, vsinger)
			// videoInfo, err := youtubeData.GetVideoInfoFromID(item.Snippet.ResourceId.VideoId)

			// if err != nil {
			// 	log.Println("err geting youtube video data")
			// 	return err

			// }
			// videos = append(videos, videoInfo)
		}
		// collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
		// collection.InsertMany(context.TODO(), struct{ videos })
		if resp.NextPageToken == "" {
			break
		}
		//	videos = nil
		nextPageToken = resp.NextPageToken
	}

	return nil
}

//GetAllVideos return totalPageNumber and video list of current page, page start with 1
func GetAllVideos(pageNum int, sortBy bson.M) (int, []*defs.Video) {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	skips := ITEMS_PER_PAGE * (pageNum - 1)
	var results []*defs.Video
	findOptions := options.Find().SetSkip(int64(skips))
	findOptions.SetLimit(ITEMS_PER_PAGE)
	findOptions.SetSort(sortBy)

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Println("collection find error:" + err.Error())
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem defs.Video
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	findOptions.SetLimit(0)

	count, err := collection.CountDocuments(context.TODO(), bson.D{{}}, options.Count())

	if err != nil {
		log.Fatal(err)
	}

	return int(math.Ceil(float64(count) / float64(ITEMS_PER_PAGE))), results

}

func QueryVideos(query bson.M, pageNum int, sortBy bson.M) (int, []*defs.Video) {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	skips := ITEMS_PER_PAGE * (pageNum - 1)
	var results []*defs.Video
	findOptions := options.Find().SetSkip(int64(skips))
	findOptions.SetLimit(ITEMS_PER_PAGE)
	findOptions.SetSort(sortBy)

	cur, err := collection.Find(context.TODO(), query, findOptions)
	if err != nil {
		log.Println("collection find error:" + err.Error())
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem defs.Video
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	findOptions.SetLimit(0)

	count, err := collection.CountDocuments(context.TODO(), query, options.Count())

	if err != nil {
		log.Fatal(err)
	}

	return int(math.Ceil(float64(count) / float64(ITEMS_PER_PAGE))), results
}

func AddUser(userData defs.UserProfile) error {
	collection := dbClient.Database("vocaloidDB").Collection("users")
	result, err := collection.InsertOne(context.TODO(), userData)
	if err != nil {
		log.Println("UserData Insert DB ERROR" + err.Error())
	}
	log.Println(result)
	return nil
}

//Update user if exists, add user if user not exists

func UpdateUser(userData defs.UserProfile) error {
	collection := dbClient.Database("vocaloidDB").Collection("users")
	filter := bson.D{{"_id", userData.UserID}}
	err := collection.FindOne(context.TODO(), filter).Err()
	if err != nil {
		AddUser(userData)
		return nil
	}

	result, err := collection.ReplaceOne(context.TODO(), filter, userData)
	if err != nil {
		log.Println("UserData Update DB ERROR" + err.Error())
	}
	log.Println(result)
	return nil
}

func UpdateVideo(videoData *defs.Video) error {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	filter := bson.D{{"_id", videoData.VideoID}}
	err := collection.FindOne(context.TODO(), filter).Err()
	if err != nil {
		InsertVideo(videoData.VideoID, videoData.Parent, videoData.Vsinger)
		return nil
	}

	result, err := collection.ReplaceOne(context.TODO(), filter, videoData)
	if err != nil {
		log.Println("VideoData Update DB ERROR" + err.Error())
	}
	log.Println(result)
	return nil
}

func GetAllVideoTags() []string {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	cur, err := collection.Find(context.TODO(), bson.D{{}})
	var results []string
	if err != nil {
		log.Println("collection find error:" + err.Error())
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem defs.Video
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem.Snippet.Tags...)
	}

	return results
}

func GetVsingerVideoTags(vsinger string) []string {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	cur, err := collection.Find(context.TODO(), bson.M{"snippet.tags": vsinger})
	var results []string
	if err != nil {
		log.Println("collection find error:" + err.Error())
	}
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem defs.Video
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem.Snippet.Tags...)
	}

	return results
}

func GetVideoCount() int64 {
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")
	result, err := collection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	return result
}

func GetUserDBCount() int64 {
	collection := dbClient.Database("vocaloidDB").Collection("users")
	result, err := collection.EstimatedDocumentCount(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return 0
	}
	return result
}

// func GetVideosByTArtist() {

// }
