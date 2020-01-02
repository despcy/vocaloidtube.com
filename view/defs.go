package viewrender

import (
	"github.com/kataras/iris"
	"github.com/yangchenxi/VOCALOIDTube/model/defs"
)

type HomePageViewModel struct {
	Ctx        iris.Context
	VideoResp  *defs.VideosResponse
	User       UserProfile
	SearchData SearchInfo
}

type SearchInfo struct {
	QueryParam string
	Qtype      string
}

type UserProfile struct {
	AvatarURL string
	Email     string
	UserID    string
}

type ApiBody struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	ReqBody string `json:"req_body"`
}

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

var (
	ErrorRequestNotRecognized   = Err{Error: "api not recognized, bad resquest", ErrorCode: "001"}
	ErrorRequestBodyParseFailed = Err{Error: "request body is not correct", ErrorCode: "002"}
	ErrorInternalFaults         = Err{Error: "internal service error", ErrorCode: "003"}
)

type PageNavItem struct {
	Link   string
	Active bool
	Number string
}

type IPBlockItem struct {
	BlockTime   string
	BlockIPAddr string
	IPLocation  string
	Message     string
}

type StatusPageViewModel struct {
	Ctx         iris.Context
	IPBlackList []IPBlockItem
	User        UserProfile
	VideoCount  string
	UserCount   string
	VisitCount  string
}

type Traffic struct {
	VisitTime     string `json:"time"`
	IpAddr        string `json: "ip"`
	Location      string `json: "location"`
	Path          string `json:"path"`
	SecurityCheck string `json: "security"`
}

type StatusPageWSData struct {
	VisitCount  string  `json: "visit"`
	TrafficData Traffic `json: "traffic"`
}
