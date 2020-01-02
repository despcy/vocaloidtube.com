package security

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/yangchenxi/VOCALOIDTube/config"
)

//放一个hashmap减少request次数

type GeoData struct {
	GeopluginRequest                string `json:"geoplugin_request"`
	GeopluginStatus                 int    `json:"geoplugin_status"`
	GeopluginDelay                  string `json:"geoplugin_delay"`
	GeopluginCity                   string `json:"geoplugin_city"`
	GeopluginRegionName             string `json:"geoplugin_regionName"`
	GeopluginCountryName            string `json:"geoplugin_countryName"`
	GeopluginLatitude               string `json:"geoplugin_latitude"`
	GeopluginLongitude              string `json:"geoplugin_longitude"`
	GeopluginLocationAccuracyRadius string `json:"geoplugin_locationAccuracyRadius"`
}

type BlockedIP struct {
	Ip          string
	TimeBlocked string
	Reason      string
	Geo         GeoData
}

var IpMap map[string]GeoData //Redis is better actually
var BlackList map[string]BlockedIP
var mutex = &sync.Mutex{}
var LRUIPCache *lru.Cache
var DBInjectHashSet map[string]bool

func init() {
	IpMap = make(map[string]GeoData)
	BlackList = make(map[string]BlockedIP)
	DBInjectHashSet = make(map[string]bool)
	set := readFileLineByLine("DBInjectionKeywords.txt")
	for i := 0; i < len(set); i++ {
		val := set[i]
		DBInjectHashSet[val] = true
	}
	var err error
	LRUIPCache, err = lru.New(config.RateLimit)
	if err != nil {
		log.Println(err)
	}
}

func AddIPToBlackList(ip string, reason string) {
	geo, err := requestGeoIP(ip)
	if err != nil {
		geo = GeoData{
			GeopluginRequest:                "unknown",
			GeopluginStatus:                 500,
			GeopluginDelay:                  "unknown",
			GeopluginCity:                   "unknown",
			GeopluginRegionName:             "unknown",
			GeopluginCountryName:            "unknown",
			GeopluginLatitude:               "unknown",
			GeopluginLongitude:              "unknown",
			GeopluginLocationAccuracyRadius: "unknown",
		}
	}
	t := time.Now()
	timeB := t.Format("2006-01-02 15:04:05")
	mutex.Lock()
	BlackList[ip] = BlockedIP{
		Ip:          ip,
		TimeBlocked: timeB,
		Reason:      reason,
		Geo:         geo,
	}
	mutex.Unlock()
}

func GetLocationFromIP(ip string) string {
	data, err := requestGeoIP(ip)
	if err != nil {
		log.Println(err)
		return "Unknown"
	}
	return data.GeopluginCity + "," + data.GeopluginRegionName + "," + data.GeopluginCountryName
}

func GetCoordFromIP(ip string) string {
	data, err := requestGeoIP(ip)
	if err != nil {
		log.Println(err)
		return "33.6404996,-117.8464902"
	}
	return data.GeopluginLatitude + "," + data.GeopluginLongitude
}

func requestGeoIP(ip string) (GeoData, error) {
	mutex.Lock()
	val, ok := IpMap[ip]
	mutex.Unlock()
	if ok {
		return val, nil
	}
	resp, err := http.Get("http://www.geoplugin.net/json.gp?ip=" + ip)

	if err != nil {
		log.Println(err)
		return GeoData{}, errors.New("http Geoip request fail for " + ip)
	}
	defer resp.Body.Close()
	var data GeoData
	if json.NewDecoder(resp.Body).Decode(&data) != nil {
		return GeoData{}, errors.New("http Geoip request decode fail for " + ip)
	}
	mutex.Lock()
	IpMap[ip] = data
	mutex.Unlock()
	return data, nil
}

func FormatDisplayIP(ip string) string {
	data := strings.Split(ip, ".")
	if len(data) != 4 {
		return ip
	}
	return data[0] + "." + data[1] + "." + "*.*"
}

type LRUItem struct {
	Lpath string
	Ltime time.Time
}

//Limit the Access Using LRU Cache, if the access rate of the same ip is too high, block
func RequestIpAntiRobot(ip string, accessPath string) {
	//if same ip access the same path in short time,block
	if LRUIPCache.Contains(ip) {
		item, ok := LRUIPCache.Get(ip)
		if ok {
			interval := time.Now().Sub(item.(LRUItem).Ltime).Seconds()
			if item.(LRUItem).Lpath == accessPath && interval < 0.01 {
				//block the ip
				go AddIPToBlackList(ip, "Request interval is "+strconv.FormatFloat(interval, 'f', -1, 64)+" seconds , this is not human behavior! 瞧瞧这是人干的事么！")
			}
		}
	}
	LRUIPCache.Add(ip, LRUItem{
		Lpath: accessPath,
		Ltime: time.Now(),
	})
}

func DBInjectionKeywordCheck(param string) bool {
	_, ok := DBInjectHashSet[param]

	return !ok

}

func readFileLineByLine(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}
