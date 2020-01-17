package apodRequester

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	apodHost     string = "https://api.nasa.gov/planetary/apod?"
	apiToken     string = "api_key="
	dateToken    string = "date="
	requestLimit int    = 10
)

type apodRequest struct {
	dayOffset int
	request   string
}

var wg = sync.WaitGroup{}
var mtx = sync.Mutex{}
var inProgress = false
var apods = make(map[int]ApodResponse)
var requestChan = make(chan apodRequest, requestLimit)
var doneChan = make(chan struct{})
var pendingRequests int

// GetApodForDateRange returns a range of APOD data
func GetApodForDateRange(numDays int, cb func([]ApodResponse), apiKey string) {
	if inProgress {
		return
	}
	inProgress = true
	go requester()
	wg.Add(numDays)

	today := time.Now()
	for i := 0; i < numDays; i++ {
		// check to see if this data has been requested already
		mtx.Lock()
		_, ok := apods[i]
		mtx.Unlock()
		if ok {
			wg.Done()
			continue
		}
		nextDay := today.AddDate(0, 0, i*-1)
		y, m, d := nextDay.Date()
		dayStr := fmt.Sprintf("%d-%d-%d", y, int(m), d)
		url := apodHost + dateToken + dayStr + "&" + apiToken + apiKey
		for {
			// request limiter
			// prevents too many simulataneous requests to NASA, which could cause the API to temporarily lock
			if pendingRequests >= requestLimit {
				time.Sleep(time.Millisecond * 100)
			} else {
				break
			}
		}
		requestChan <- apodRequest{dayOffset: i, request: url}
	}

	wg.Wait()
	doneChan <- struct{}{}

	rtn := make([]ApodResponse, numDays)
	for k, v := range apods {
		rtn[k] = v
	}
	cb(rtn)
	inProgress = false
}

func requester() {
	for {
		select {
		case req := <-requestChan:
			go requestAndParse(req.dayOffset, req.request)
		case <-doneChan:
			fmt.Println("Done...")
			return
		}
	}
}

func requestAndParse(index int, req string) {
	pendingRequests++
	defer func() {
		wg.Done()
		pendingRequests--
	}()
	d, e := generateGetRequest(req)
	if e != nil {
		fmt.Println("Error generating data: ", index)
		return
	}
	apod, e := UnmarshalApodResponse(d)
	if e != nil {
		return
	}
	mtx.Lock()
	apods[index] = apod
	mtx.Unlock()
}

func generateGetRequest(req string) ([]byte, error) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Get(req)
	if err != nil {
		fmt.Println("Error retrieving data from NASA servers: ", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Could not read response from NASA servers: ", err)
		return nil, err
	}
	return data, nil
}
