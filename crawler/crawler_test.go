package crawler

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	pid := "RDhsr-Mu4Mdwk"
	vsinger := "Unknown"
	CrawlFromPlayList(pid, vsinger)
	time.Sleep(10 * time.Second)
}
