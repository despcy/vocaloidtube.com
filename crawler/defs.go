package crawler

// const (
// 	IDLE = "i"
// 	BUSY = "b"
// )

type controlData struct {
	//state     string
	newData   []string //new video ids
	workerNum int
	errorData string //比如出现token用完之类的
}

type VideoQueue []string

type dataChan chan string

type controlChan chan controlData

type Executor func(data dataChan)
