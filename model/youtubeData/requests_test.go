package youtubeData

import (
	"log"
	"testing"
)

// func TestMain(m *testing.M) {

// 	m.Run()

//}

func TestGetVideoInfoFromID(t *testing.T) {
	_, err := GetVideoInfoFromID("vtZQ0DYgjdE")
	if err != nil {
		t.Errorf("Error of : %v", err)
	}
}

func TestGetSuggestedVideosID(t *testing.T) {
	id := "ffCPqEqGgFc"
	resp, _ := GetSuggestedVideosIDFromVideoID(id)

	for _, item := range resp {
		log.Println(item)
	}
}

func TestGetPlayListItemsFromPlayListID(t *testing.T) {

	resp, err := GetPlayListItemsFromPlayListID("PL_vdmUqgn18VVtdSieJfJw1ce9tB4Xppv", "")

	if err != nil {
		t.Errorf("Error of : %v", err)

	}
	for _, item := range resp.Items {
		log.Printf(item.Snippet.Title)
	}
	if resp.NextPageToken != "" {
		log.Printf("Start Next Page")
		resp, err := GetPlayListItemsFromPlayListID("RDARt2fVT33Lw", resp.NextPageToken)

		if err != nil {
			t.Errorf("Error of : %v", err)

		}
		for _, item := range resp.Items {
			log.Printf(item.Snippet.Title)
		}

	}

}
