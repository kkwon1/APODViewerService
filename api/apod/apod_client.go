package apod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kkwon1/APODViewerService/models"
	log "github.com/sirupsen/logrus"
)

const layoutISO = "2006-01-02"

type ApodClient interface {
	GetBatchImages(w http.ResponseWriter, r *http.Request)
}

type apodClient struct{}

func NewApodClient() ApodClient {
	return &apodClient{}
}

// GetBatchImages will retrieve multiple images from NASA APOD API
func (ac *apodClient) GetBatchImages(w http.ResponseWriter, r *http.Request) {
	count, convErr := strconv.Atoi(r.URL.Query().Get("count"))

	if convErr != nil {
		log.Errorln(convErr)
	}

	page, convErr := strconv.Atoi(r.URL.Query().Get("page"))

	if convErr != nil {
		page = 0
	}

	startDate, endDate := getStartEndDates(count, page)
	log.Printf("End date: %s", endDate)
	log.Printf("Start date: %s", startDate)

	httpResp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s&start_date=%s&end_date=%s", os.Getenv("NASA_API_KEY"), startDate, endDate))
	if err != nil {
		log.Errorln(err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorln(err)
	}

	var responseObjects []models.NASAApodObject
	err = json.Unmarshal(body, &responseObjects)
	if err != nil {
		log.Errorln(err)
	}
	apodObjects := mapFcn(responseObjects, convertToApodObject)

	log.Printf("Successfully retrieved %d number of images", count)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apodObjects)
}

func getStartEndDates(count int, page int) (string, string) {
	// page * count of batch images will be the offset of date
	offset := page * count
	// Get today's date in UTC (That is how APOD operates)
	// convert to specified ISO, and find start date to call API and retrieve
	// the correct number of images
	today := time.Now()
	endDateTime := today.AddDate(0, 0, offset*-1)
	endDate := endDateTime.Format(layoutISO)
	startDate := endDateTime.AddDate(0, 0, (count-1)*-1).Format(layoutISO)

	return startDate, endDate
}

func convertToApodObject(NasaApodObject models.NASAApodObject) models.ApodObject {
	apodObject := models.ApodObject{
		ApodURL:     NasaApodObject.Url,
		ApodHDURL:   NasaApodObject.Hdurl,
		ApodName:    NasaApodObject.Title,
		ApodDate:    NasaApodObject.Date,
		MediaType:   NasaApodObject.Media_type,
		Description: NasaApodObject.Explanation,
		Copyright:   NasaApodObject.Copyright,
	}

	return apodObject
}

func mapFcn(vs []models.NASAApodObject, f func(models.NASAApodObject) models.ApodObject) []models.ApodObject {
	vsm := make([]models.ApodObject, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
