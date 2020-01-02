package security

import (
	"log"
	"testing"
)

func TestGeoIPRequest(t *testing.T) {
	log.Println(GetLocationFromIP("115.216.215.145"))
}

func TestHandleIP(t *testing.T) {
	log.Println(FormatDisplayIP("192.168.1.1"))
}

func TestIPCoor(t *testing.T) {
	log.Println(GetCoordFromIP("135.210.245.142"))
}

func TestDBInjectKeyword(t *testing.T) {
	log.Println(DBInjectionKeywordCheck("=%20--"))
}
