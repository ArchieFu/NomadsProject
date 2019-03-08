package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/thedevsaddam/gojsonq"
	"io/ioutil"
	"net/http"
	"sort"
)

var (
	planeApiUrl = "https://flights.ctrip.com/itinerary/api/12808/products"
)

type DepartureAirportInfo struct {
	CityName    string
	AirportName string
	Terminal    string
}

type ArrivalAirportInfo struct {
	CityName    string
	AirportName string
	Terminal    string
}

type FlightInfo struct {
	FlightNumber  string
	AirlineName   string
	DepartureDate string
	ArrivalDate   string
	DepInfo       DepartureAirportInfo
	ArrInfo       ArrivalAirportInfo
}

type Price struct {
	PrintPrice float64
	Rate       float64
	CabinClass string
}

type PlaneTask struct {
	FliInfo FlightInfo
	Prices  []Price
}

type PlaneTaskResult struct {
	PlaneTasks []PlaneTask
	by         func(p, q *PlaneTask) bool
}

type SortBy func(p, q *PlaneTask) bool

func (r PlaneTaskResult) Len() int {
	return len(r.PlaneTasks)
}

func (r PlaneTaskResult) Less(i, j int) bool {
	return r.by(&r.PlaneTasks[i], &r.PlaneTasks[j])
}

func (r PlaneTaskResult) Swap(i, j int) {
	r.PlaneTasks[i], r.PlaneTasks[j] = r.PlaneTasks[j], r.PlaneTasks[i]
}

func SortPrice(price []PlaneTask, by SortBy) {
	sort.Sort(PlaneTaskResult{price, by})
}

type downloaderService struct {
	Dcity string
	Acity string
	Date  string
}

type AirportParam struct {
	Dcity string `json:"dcity"`
	Acity string `json:"acity"`
	Date  string `json:"date"`
}

type Payload struct {
	FlightWay     string         `json:"flightWay"`
	ClassType     string         `json:"classType"`
	HasChild      bool           `json:"hasChild"`
	HasBaby       bool           `json:"hasBaby"`
	SearchIndex   int            `json:"searchIndex"`
	AirportParams []AirportParam `json:"airportParams"`
}

func NewdownloaderService(dcity, acity, data string) *downloaderService {
	return &downloaderService{Dcity: dcity, Acity: acity, Date: data}
}

func (s *downloaderService) GetTicketResult() (*PlaneTaskResult, error) {
	planeTaskResult := &PlaneTaskResult{}
	payload, err := s.getPayload()
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBufferString(payload)
	req, _ := http.NewRequest("POST", planeApiUrl, buff)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("referer", "https://flights.ctrip.com/itinerary/oneway/bjs-sha,pvg?date=2019-04-04")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("Postman-Token", "29ff4eac-65c9-49bc-b571-34e520f77177")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	planeTaskResult, err = doJsonToData(body)
	if err != nil {
		return nil, err
	}
	return planeTaskResult, nil
}

func (s *downloaderService) getPayload() (string, error) {
	airportParams := AirportParam{}
	airportParams.Acity = s.Acity
	airportParams.Dcity = s.Dcity
	airportParams.Date = s.Date

	payload := Payload{}
	payload.AirportParams = make([]AirportParam, 0) //
	payload.FlightWay = "Oneway"
	payload.ClassType = "ALL"
	payload.HasChild = false
	payload.HasBaby = false
	payload.SearchIndex = 1
	//payload.AirportParams[0] = airportParams
	payload.AirportParams = append(payload.AirportParams, airportParams)

	retStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("json marshal error", err)
		return "", err
	}
	//fmt.Println("request", string(retStr))

	return string(retStr), nil
}

func doJsonToData(body []byte) (planeTaskResult *PlaneTaskResult, err error) {
	data := gojsonq.New().JSONString(string(body)).From("data.routeList").Select("legs")
	infoList, ok := data.Get().([]interface{})
	if !ok {
		fmt.Println("get route list error")
	}

	planeTaskResult = &PlaneTaskResult{}
	planeTaskResult.PlaneTasks = []PlaneTask{}

	//获取到出发地到目的地的所有航班列表
	for _, routes := range infoList {
		infoMap, ok := routes.(map[string]interface{})
		if !ok {
			fmt.Println("get route list error")
		}
		planeTask := PlaneTask{}
		planeTask.Prices = []Price{}

		//获取到legs 单次航班的所有信息
		airInfos := infoMap["legs"].([]interface{})
		if len(airInfos) != 1 {
			// 中转的pass掉 (legs > 2 需要中转)
			continue
		}
		for _, airInfo := range airInfos {
			flight := airInfo.(map[string]interface{})["flight"]
			flightRes := parseflightInfo(flight)
			planeTask.FliInfo = flightRes

			cabins := airInfo.(map[string]interface{})["cabins"]
			cabininfo := parseCabinsInfo(cabins.([]interface{}))
			planeTask.Prices = cabininfo
		}

		planeTaskResult.PlaneTasks = append(planeTaskResult.PlaneTasks, planeTask)
	}

	return planeTaskResult, nil
}

func parseflightInfo(flight interface{}) (resFlight FlightInfo) {
	flightInfo := flight.(map[string]interface{})
	depInfo := flightInfo["departureAirportInfo"].(map[string]interface{})
	terminal := depInfo["terminal"].(map[string]interface{})

	arrInfo := flightInfo["arrivalAirportInfo"].(map[string]interface{})
	arrTerminal := depInfo["terminal"].(map[string]interface{})

	resFlight.FlightNumber = flightInfo["flightNumber"].(string)
	resFlight.AirlineName = flightInfo["airlineName"].(string)
	resFlight.DepartureDate = flightInfo["departureDate"].(string)
	resFlight.ArrivalDate = flightInfo["arrivalDate"].(string)
	resFlight.DepInfo.CityName = depInfo["cityName"].(string)
	resFlight.DepInfo.AirportName = depInfo["airportName"].(string)
	resFlight.DepInfo.Terminal = terminal["name"].(string)

	resFlight.ArrInfo.CityName = arrInfo["cityName"].(string)
	resFlight.ArrInfo.AirportName = arrInfo["airportName"].(string)
	resFlight.ArrInfo.Terminal = arrTerminal["name"].(string)

	return resFlight
}

func parseCabinsInfo(cabins []interface{}) (resPrice []Price) {
	resPrice = []Price{}
	for _, carbin := range cabins {
		info := carbin.(map[string]interface{})
		price := info["price"].(map[string]interface{})

		p := Price{}
		p.PrintPrice = price["printPrice"].(float64)
		p.CabinClass = info["cabinClass"].(string)
		p.Rate = price["rate"].(float64)

		resPrice = append(resPrice, p)
	}

	return resPrice
}

func FloorPrice(r *PlaneTaskResult) {
	SortPrice(r.PlaneTasks,
		func(p, q *PlaneTask) bool {
			return p.Prices[0].PrintPrice < q.Prices[0].PrintPrice	//递增排序
		})
}

// 获取排序后第一个航班列表
func GetOneflight(planeTask *PlaneTaskResult) PlaneTask {
	fmt.Println("len plan task:", len(planeTask.PlaneTasks))
	if len(planeTask.PlaneTasks) == 0 {
		return PlaneTask{}
	}

	return planeTask.PlaneTasks[0]
}