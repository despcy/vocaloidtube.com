package crawler

import (
	"log"
	"strconv"
	"strings"

	"github.com/yangchenxi/VOCALOIDTube/model/youtubeData"
)

//producer/comsumer model for web crawling and processing
//感觉不是正经的producer/consumer，但是就这么玩吧
// 一个runner分配任务给多个worker，然后workerfinish以后继续给分任务
type Runner struct {
	TaskQueue       []string
	NumberOfWorkers int
}

func NewRunner(workerNum int, tasks []string) *Runner {
	return &Runner{
		TaskQueue:       tasks,
		NumberOfWorkers: workerNum,
	}
}

//producer
func (r *Runner) startRunner() {
	workers := make([]dataChan, r.NumberOfWorkers)
	ctrlChan := make(controlChan, r.NumberOfWorkers) //防止阻塞
	for i := 0; i < r.NumberOfWorkers; i++ {
		workers[i] = make(dataChan, 1)
		workers[i] <- r.TaskQueue[0]
		r.TaskQueue = r.TaskQueue[1:]
		go process(i, workers[i], ctrlChan)
	}

	for {
		select {
		case c := <-ctrlChan:

			if c.errorData != "" {
				log.Println("ctrl channel Error:" + c.errorData)
				if strings.Contains(c.errorData, "Daily Limit Exceeded") {
					//TODO:wait another day

					return
				}
			} else {
				//handle new data to queue
				r.TaskQueue = append(r.TaskQueue, c.newData...)
			}
			//assign new work
			workers[c.workerNum] <- r.TaskQueue[0]
			r.TaskQueue = r.TaskQueue[1:]
			//log.Println(r.TaskQueue)
		}
	}

}

//consumer
func process(id int, dchan dataChan, ctrlChan controlChan) {
	//select 等候分配data
	for {
		select {
		case d := <-dchan:
			log.Println("worker " + strconv.Itoa(id) + "start")
			data, err := processVideoID(d)
			if err != nil {
				ctrlChan <- controlData{
					workerNum: id,
					errorData: err.Error(),
				}
			} else {
				ctrlChan <- controlData{
					workerNum: id,
					newData:   data,
				}
			}
			log.Println("worker " + strconv.Itoa(id) + "end")

		}

	}
}

func processVideoID(vid string) ([]string, error) {
	//time.Sleep(1 * time.Second) //TODO:这里回头调一下，反正每个线程单独一个channel，不怕阻塞
	resp, err := youtubeData.GetSuggestedVideosIDFromVideoID(vid)
	if err != nil {
		return nil, err
	}
	//TODO:add filter and Bayes

	//TODO: store the data in resp to db
	//if in db, remove from resp, remember to add parent 字段
	return resp, nil

}
