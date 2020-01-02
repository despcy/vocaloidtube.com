package defs

import (
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	VideoID          string                  `bson:"_id" json:"videoID"`
	Snippet          youtube.VideoSnippet    `json:"snippet"`
	Statistic        youtube.VideoStatistics `json:"statistic"`
	VideoKind        string                  `json:"videoKind"` //live or video
	SiteViewCount    int                     `json:"siteViewCount"`
	SitefavCount     int                     `json:"sitefavCount"`
	SitelikeCount    int                     `json:"sitelikeCount"`
	Sitedislikecount int                     `json:"sitedislikeCount"`
	Vsinger          string                  `json:"vsinger"`
	Parent           string                  `json:"parent"`
}
type VideosResponse struct {
	PageNum   int      `json:"pageNum"`
	TotalPage int      `json:"totalPages"`
	Sort      string   `json:"sortBy"`
	Videos    []*Video `json:"videos"`
}

type UserProfile struct {
	UserID    string `bson:"_id" json:"userID"`
	UserEmail string `json:"email"`
	AvatarURL string `json:"avatarURL`
}
