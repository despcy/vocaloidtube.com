package crawler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/jbrukh/bayesian"
	"github.com/yangchenxi/VOCALOIDTube/model/dbops"
	"github.com/yangchenxi/VOCALOIDTube/model/youtubeData"
)

// func TestBayes(t *testing.T) {
// 	//	testVideoID := "O0TtDeDiHcE"
// 	// Check the connection
// 	const (
// 		Good bayesian.Class = "Good"
// 		Bad  bayesian.Class = "Badd"
// 	)

// 	//下一步就是从category为10的里边随机抓取几个然后测试一下分类器，顺便找一些反面的tag
// 	classifier := bayesian.NewClassifier(Good, Bad)
// 	goodStuff := GetVocalTagsFromFile()
// 	badStuff := []string{""}
// 	classifier.Learn(goodStuff, Good)
// 	classifier.Learn(badStuff, Bad)

// 	//resp, _ := youtubeData.GetSuggestedVideosIDFromVideoID(testVideoID)

// 	//for _, item := range resp {
// 	runBayes("OmHDuywLloY", classifier)
// 	//}

// }

func runBayes(id string, classifier *bayesian.Classifier) {
	videoInfo, _ := youtubeData.GetVideoInfoFromID(id)
	scores, likely, _ := classifier.LogScores(
		videoInfo.Snippet.Tags,
	)
	//TODO:结果并不是很好，所以贝叶斯就放在vsinger上吧加一层tag的filter，宁可错杀也不放过
	//db.youtubeVideos.find({"snippet.tags":{$in ["VOCALOID","vocaloid"]}})
	log.Println("-------------------")
	log.Println(videoInfo.Id)
	log.Println(videoInfo.Snippet.Title)
	log.Println(videoInfo.Snippet.Tags)
	log.Println(scores)
	log.Println(classifier.Classes[likely])
}

// func TestGetTags(t *testing.T) {
// 	SaveVocalTags()
// }

// func TestReadTags(t *testing.T) {
// 	GetVocalTagsFromFile()
// }

func TestTrainClassfier(t *testing.T) {
	//慎用
	trainClassifier()
}

func TestVsingerBayes(t *testing.T) {
	classifier, _ := bayesian.NewClassifierFromFile("vsingerClassifier.dat")
	runBayes("NMY9PJuoTmk", classifier)
}

func trainClassifier() {
	VsingerList := []bayesian.Class{
		"Unknown",
		"Mo Qingxian",
		"Zhiyu Moke",
		"Haruno Sora",
		"Ken",
		"Kaori",
		"Chris",
		"Amy",
		"Mirai Komachi",
		"Kizuna Akari",
		"LUMi",
		"Kobayashi Matcha",
		"Masaoka Azuki",
		"Yuezheng Longya",
		"Yumemi Nemu",
		"UNI",
		"Macne Petit",
		"CYBER SONGMAN",
		"Otomachi Una",
		"Xingchen",
		"Fukase",
		"Otori Kohaku",
		"DAINA",
		"DEX",
		"RUBY",
		"ARSLOID",
		"Sachiko",
		"Yuezheng Ling",
		"Xin Hua",
		"CYBER DIVA",
		"Chika",
		"Rana",
		"Tohoku Zunko",
		"flower",
		"kanon&anon",
		"kokone",
		"Macne Nana",
		"Merli",
		"MAIKA",
		"YOHIOloid",
		"YANHE",
		"YUU",
		"KYO",
		"WIL",
		"AVANNA",
		"MAYU",
		"galaco",
		"Luo Tianyi",
		"Aoki Lapis",
		"IA",
		"Clara",
		"Yuzuki Yukari",
		"Bruno",
		"CUL",
		"OLIVER",
		"Tone Rion",
		"SeeU",
		"Utatane Piko",
		"Nekomura Iroha",
		"Lily",
		"BIG AL",
		"Hiyama Kiyoteru",
		"Kaai Yuki",
		"SF-A2 miki",
		"GUMI",
		"Megurine Luka",
		"Camui Gackpo",
		"Kagamine Len",
		"Kagamine Rin",
		"Hatsune Miku",
		"KAITO",
		"MEIKO",
	}

	Classifier = bayesian.NewClassifier(VsingerList...)
	log.Println(len(VsingerList))
	//f, _ := os.OpenFile("JData.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	var TrainData []BayesTagItems
	jsonFile, _ := os.Open("data.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &TrainData)
	for _, val := range TrainData {
		//tags := getVsingerTagsFromDB(string(val.Vsing))
		// JsonData := BayesTagItems{
		// 	VsingerName: string(val),
		// 	VsingerTags: tags,
		// }
		//TrainData = append(val.VsingerTags, val.VsingerName)

		Classifier.Learn(val.VsingerTags, val.VsingerName)
		log.Println(val)
	}
	//JByte, _ := json.Marshal(TrainData)
	//f.Write(JByte)
	//f.Close()
	//	Classifier.WriteClassesToFile("vsingerClasses")
	Classifier.WriteToFile("vsingerClassifier.dat")
}

func getVsingerTagsFromDB(vsinger string) []string {
	return dbops.GetVsingerVideoTags(vsinger)
}
