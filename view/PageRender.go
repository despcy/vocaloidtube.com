package viewrender

import (
	"strconv"
	"strings"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/yangchenxi/VOCALOIDTube/config"
	"github.com/yangchenxi/VOCALOIDTube/model/defs"
	"github.com/yangchenxi/VOCALOIDTube/security"
)

func RenderHomePage(vm HomePageViewModel) {

	AllVideos(vm.Ctx, vm.VideoResp)
}

func RenderVsingerPage(ctx iris.Context) {
	renderUserProfile(ctx)
	ctx.View("vsinger.html")
}

func RenderStatusPage(ctx iris.Context, userCount string, videoCount string, visits string, blacklist []security.BlockedIP) {
	ctx.ViewData("BlackList", blacklist)
	ctx.ViewData("UserCount", userCount)
	ctx.ViewData("VideoCount", videoCount)
	components := strings.Split(config.SERVER_DOMAIN, "://")
	ctx.ViewData("Host", components[1])
	ctx.ViewData("Visits", visits)
	renderUserProfile(ctx)
	ctx.View("status.html")
}

func AllVideos(ctx iris.Context, VideoResp *defs.VideosResponse) {

	renderVideoList(ctx, VideoResp)
	curPage := VideoResp.PageNum
	maxPage := VideoResp.TotalPage
	nextPage := VideoResp.PageNum + 1
	prevPage := VideoResp.PageNum - 1
	sort := VideoResp.Sort
	var url string
	url = "/videos?sortBy=" + sort + "&page="
	renderPageNav(ctx, curPage, maxPage, nextPage, prevPage, url)
	var filterUrls [3]string
	filterUrls[0] = "/videos?sortBy=click&page=" + strconv.Itoa(curPage)
	filterUrls[1] = "/videos?sortBy=like&page=" + strconv.Itoa(curPage)
	filterUrls[2] = "/videos?sortBy=recent&page=" + strconv.Itoa(curPage)
	renderFilterOptions(ctx, filterUrls)
	renderUserProfile(ctx)
	ctx.ViewData("Title", "All Videos Sorted By Most "+sort)
	ctx.View("videos.html")

}

func RenderPlayPage(vm HomePageViewModel) {
	ctx := vm.Ctx
	renderUserProfile(ctx)
	videoTitle := vm.VideoResp.Videos[0].Snippet.Title
	videoID := vm.VideoResp.Videos[0].VideoID
	ctx.ViewData("Title", videoTitle)
	ctx.ViewData("Vid", videoID)
	ctx.View("play.html")
}

func renderUserProfile(ctx iris.Context) {
	session := sessions.Get(ctx)
	if session.GetString("UserEmail") == "" {
		return
	}
	profile := UserProfile{
		Email:     session.GetString("UserEmail"),
		AvatarURL: session.GetString("UserAvatarURL"),
		UserID:    session.GetString("UserID"),
	}
	ctx.ViewData("UserProfile", profile)
}

func renderFilterOptions(ctx iris.Context, filter [3]string) {
	ctx.ViewData("Sort", filter)
}

func renderPageNav(ctx iris.Context, curPage int, maxPage int, nextPage int, prevPage int, url string) {

	var pageData []PageNavItem
	if curPage == 1 {
		pageData = append(pageData, PageNavItem{Link: "#", Active: false, Number: "-1"})
	} else {
		pageData = append(pageData, PageNavItem{Link: url + strconv.Itoa(prevPage), Active: true, Number: "-1"})
	}

	for i := 1; i <= maxPage; i++ {
		if i > 3 && i < prevPage {
			if i == 4 {
				//print .....
				pageData = append(pageData, PageNavItem{Link: url + "4", Active: false, Number: "..."})
			}
		} else if i > nextPage && i < maxPage-2 {
			if i == maxPage-3 {
				pageData = append(pageData, PageNavItem{Link: url + strconv.Itoa(maxPage-3), Active: false, Number: "..."})
			}
		} else if i == curPage {
			pageData = append(pageData, PageNavItem{Link: url + strconv.Itoa(i), Active: true, Number: strconv.Itoa(i)})
		} else {
			pageData = append(pageData, PageNavItem{Link: url + strconv.Itoa(i), Active: false, Number: strconv.Itoa(i)})

		}
	}

	if curPage == maxPage {
		pageData = append(pageData, PageNavItem{Link: "#", Active: false, Number: "-1"})
	} else {
		pageData = append(pageData, PageNavItem{Link: url + strconv.Itoa(nextPage), Active: true, Number: "-1"})
	}
	ctx.ViewData("PageData", pageData)

}

func renderVideoList(ctx iris.Context, videoResp *defs.VideosResponse) {

	ctx.ViewData("VideoResp", videoResp)

}

func RenderSearchPage(vm HomePageViewModel) {
	renderVideoList(vm.Ctx, vm.VideoResp)
	VideoResp := vm.VideoResp
	ctx := vm.Ctx
	curPage := VideoResp.PageNum
	maxPage := VideoResp.TotalPage
	nextPage := VideoResp.PageNum + 1
	prevPage := VideoResp.PageNum - 1
	sort := VideoResp.Sort
	var url string
	url = "/search?q=" + vm.SearchData.QueryParam + "&type=" + vm.SearchData.Qtype + "&sortBy=" + sort + "&page="
	renderPageNav(ctx, curPage, maxPage, nextPage, prevPage, url)
	var filterUrls [3]string
	filterUrls[0] = "/search?q=" + vm.SearchData.QueryParam + "&type=" + vm.SearchData.Qtype + "&sortBy=click&page=" + strconv.Itoa(curPage)
	filterUrls[1] = "/search?q=" + vm.SearchData.QueryParam + "&type=" + vm.SearchData.Qtype + "&sortBy=like&page=" + strconv.Itoa(curPage)
	filterUrls[2] = "/search?q=" + vm.SearchData.QueryParam + "&type=" + vm.SearchData.Qtype + "&sortBy=recent&page=" + strconv.Itoa(curPage)
	renderFilterOptions(ctx, filterUrls)
	renderUserProfile(ctx)
	ctx.ViewData("Title", "Search Results For: "+vm.SearchData.QueryParam+" ,Sorted By Most "+sort)
	ctx.View("videos.html")
}
