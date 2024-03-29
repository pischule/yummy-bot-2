package ocr

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type taskStruct struct {
	Text                    string `xml:",chardata"`
	ID                      string `xml:"id,attr"`
	RegistrationTime        string `xml:"registrationTime,attr"`
	StatusChangeTime        string `xml:"statusChangeTime,attr"`
	Status                  string `xml:"status,attr"`
	FilesCount              string `xml:"filesCount,attr"`
	Credits                 string `xml:"credits,attr"`
	EstimatedProcessingTime string `xml:"estimatedProcessingTime,attr"`
	ResultUrl               string `xml:"resultUrl,attr"`
}

type taskResponse struct {
	XMLName xml.Name   `xml:"response"`
	Text    string     `xml:",chardata"`
	Task    taskStruct `xml:"task"`
}

type abbyyDocument struct {
	XMLName        xml.Name `xml:"document"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Version        string   `xml:"version,attr"`
	Producer       string   `xml:"producer,attr"`
	Languages      string   `xml:"languages,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	Page           struct {
		Text           string `xml:",chardata"`
		Width          int    `xml:"width,attr"`
		Height         int    `xml:"height,attr"`
		Resolution     string `xml:"resolution,attr"`
		OriginalCoords string `xml:"originalCoords,attr"`
		Block          []struct {
			Chardata  string `xml:",chardata"`
			BlockType string `xml:"blockType,attr"`
			BlockName string `xml:"blockName,attr"`
			L         int    `xml:"l,attr"`
			T         int    `xml:"t,attr"`
			R         int    `xml:"r,attr"`
			B         int    `xml:"b,attr"`
			Region    struct {
				Text string `xml:",chardata"`
				Rect []struct {
					Text string `xml:",chardata"`
					L    int    `xml:"l,attr"`
					T    int    `xml:"t,attr"`
					R    int    `xml:"r,attr"`
					B    int    `xml:"b,attr"`
				} `xml:"rect"`
			} `xml:"region"`
			Text struct {
				Text string `xml:",chardata"`
				Par  []struct {
					Text        string `xml:",chardata"`
					LineSpacing string `xml:"lineSpacing,attr"`
					Align       string `xml:"align,attr"`
					LeftIndent  string `xml:"leftIndent,attr"`
					StartIndent string `xml:"startIndent,attr"`
					Line        []struct {
						Text       string `xml:",chardata"`
						Baseline   string `xml:"baseline,attr"`
						L          int    `xml:"l,attr"`
						T          int    `xml:"t,attr"`
						R          int    `xml:"r,attr"`
						B          int    `xml:"b,attr"`
						Formatting struct {
							Text       string `xml:",chardata"`
							Lang       string `xml:"lang,attr"`
							CharParams []struct {
								Text           string `xml:",chardata"`
								L              int    `xml:"l,attr"`
								T              int    `xml:"t,attr"`
								R              int    `xml:"r,attr"`
								B              int    `xml:"b,attr"`
								Suspicious     string `xml:"suspicious,attr"`
								IsTab          string `xml:"isTab,attr"`
								TabLeaderCount string `xml:"tabLeaderCount,attr"`
							} `xml:"charParams"`
						} `xml:"formatting"`
					} `xml:"line"`
				} `xml:"par"`
			} `xml:"text"`
		} `xml:"block"`
	} `xml:"page"`
}

func (t *taskStruct) isActive() bool {
	return t.Status == "InProgress" || t.Status == "Queued"
}

func processImage(imageReader io.Reader, username string, password string) (taskStruct, error) {
	log.Println("processImage")
	url := "https://cloud-eu.ocrsdk.com/processImage?language=Russian&exportFormat=xml&profile=textExtraction"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, imageReader)

	if err != nil {
		log.Println("create task request creation failed", err)
		return taskStruct{}, err
	}
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "image/jpeg")

	res, err := client.Do(req)
	if err != nil {
		log.Println("task creation request failed", err)
		return taskStruct{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("task create body close failed")
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("body read failed", err)
		return taskStruct{}, err
	}
	var taskResponse taskResponse
	if err := xml.Unmarshal(body, &taskResponse); err != nil {
		log.Println("task creation response unmarshall failed", err)
		return taskStruct{}, err
	}
	return taskResponse.Task, nil
}

func getTaskStatus(task taskStruct, username string, password string) (taskStruct, error) {
	url := "https://cloud-eu.ocrsdk.com/getTaskStatus?taskId=" + task.ID
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("task check request creation failed", err)
		return task, err
	}
	req.SetBasicAuth(username, password)

	res, err := client.Do(req)
	if err != nil {
		log.Println("task check request failed")
		return task, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("task check body close failed")
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("task check request body read failed", err)
		return task, err
	}
	var taskResponse taskResponse
	if err := xml.Unmarshal(body, &taskResponse); err != nil {
		log.Println("task status response unmarshall failed", err)
		return task, err
	}
	return taskResponse.Task, nil
}

func downloadResult(task taskStruct) (abbyyDocument, error) {
	res, err := http.Get(task.ResultUrl)
	if err != nil {
		log.Println("download task result failed", err)
		return abbyyDocument{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("task result body close failed")
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("task result body read failed", err)
		return abbyyDocument{}, err
	}
	var document abbyyDocument
	err = xml.Unmarshal(body, &document)
	if err != nil {
		log.Println("task result unmarshall failed")
		return abbyyDocument{}, err
	}
	return document, nil
}

func recognizeFile(imageReader io.Reader, username string, password string) (abbyyDocument, error) {
	log.Println("uploading...")
	task, err := processImage(imageReader, username, password)
	if err != nil {
		log.Println("processImage failed", err)
		return abbyyDocument{}, err
	}
	if task.Status == "NotEnoughCredits" {
		return abbyyDocument{}, fmt.Errorf("not enough credits to process the document")
	}
	log.Printf("task id = %s", task.ID)
	log.Printf("waiting")
	for i := 0; i < 10; i++ {
		if !task.isActive() {
			break
		}
		time.Sleep(time.Second * 5)
		task, _ = getTaskStatus(task, username, password)
	}

	if task.Status == "Completed" && task.ResultUrl != "" {
		return downloadResult(task)
	} else {
		return abbyyDocument{}, fmt.Errorf("error processing task")
	}
}
