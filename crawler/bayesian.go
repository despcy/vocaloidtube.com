package crawler

import (
	"log"

	"github.com/jbrukh/bayesian"
)

var VsingerList []bayesian.Class
var Classifier *bayesian.Classifier

type BayesTagItems struct {
	VsingerName bayesian.Class
	VsingerTags []string
}

func init() {
	classFromFile, err := bayesian.NewClassifierFromFile("vsingerClassifier.dat")
	if err != nil {
		log.Println(err)
	}
	Classifier = classFromFile
}

func ClassifyVsinger(tags []string) string {
	_, likely, _ := Classifier.LogScores(tags)
	return string(Classifier.Classes[likely])
}

// func SaveVocalTags() {
// 	tags := dbops.GetAllVideoTags()
// 	log.Println(len(tags))
// 	file, _ := json.Marshal(tags)
// 	_ = ioutil.WriteFile("tags.json", file, 0644)

// }

// func GetVocalTagsFromFile() []string {
// 	jbyte, _ := ioutil.ReadFile("tags.json")
// 	var tags []string
// 	json.Unmarshal(jbyte, &tags)
// 	log.Println(len(tags))
// 	return tags
// }
