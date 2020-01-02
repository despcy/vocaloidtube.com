package youtubeData

import (
	"log"
	"net/http"

	config "github.com/yangchenxi/VOCALOIDTube/config"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	service *youtube.Service
)

func init() {
	developerKey := config.YOUTUBE_API_KEY
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	var conerr error

	service, conerr = youtube.New(client)

	if conerr != nil {
		log.Fatalf("Error creating new YouTube client: %v", conerr)
		panic(conerr.Error())
	}

}

//GetVideoInfoFromID gets the video snippet and statistics info and combine
func GetVideoInfoFromID(id string) (*youtube.Video, error) {

	call := service.Videos.List("snippet,statistic").Id(id)

	response, err := call.Do()
	if err != nil {
		log.Println("Call err", err)
		return nil, err
	}

	if len(response.Items) == 0 {
		log.Println("No such video err", err)
		return nil, err
	}

	return response.Items[0], nil
}

func GetSuggestedVideosIDFromVideoID(id string) ([]string, error) {
	log.Println("get suggestions for " + id)
	pageToken := ""
	call := service.Search.List("id").MaxResults(50).RelatedToVideoId(id).SafeSearch("strict").Type("video").VideoEmbeddable("true")
	var result []string
	for {
		if pageToken != "" {
			call.PageToken(pageToken)
		}
		response, err := call.Do()
		if err != nil {
			log.Println("Call err", err)
			return nil, err
		}

		if len(response.Items) == 0 {
			log.Println("No such video err", err)
			return nil, err
		}

		pageToken = response.NextPageToken

		for _, item := range response.Items {

			result = append(result, item.Id.VideoId)
		}

		if pageToken == "" {
			break
		}
	}

	return result, nil
}

func GetPlayListItemsFromPlayListID(id string, pageToken string) (*youtube.PlaylistItemListResponse, error) {
	call := service.PlaylistItems.List("snippet")
	call.PlaylistId(id)
	call.MaxResults(50)
	if pageToken != "" {
		call.PageToken(pageToken)
	}

	response, err := call.Do()
	if err != nil {
		log.Println("Call err", err)
		return nil, err
	}

	return response, nil
}
