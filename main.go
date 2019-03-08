package main

import (
	"./downloader"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

var CityCode = map[string]string{
	//"贵阳": "KWE",
	//"北京": "BJS",
	"重庆":   "CKG",
	"成都":   "CTU",
	"武汉":   "WUH",
	"海口":   "HAK",
	"郑州":   "CGO",
	"上海浦东": "PVG",
	"上海虹桥": "SHA",
	"长沙":   "CSX",
	"杭州":   "HGH",
	"厦门":   "XMN",
	"昆明":   "KMG",
	"桂林":   "KWL",
	"南昌":   "KHN",
	"兰州":   "LHW",
	"济南":   "TNA",
	"青岛":   "TAO",
	"大连":   "DLC",
	"南通":   "NTG",
	"长春":   "CGQ",
	"沈阳":   "SHE",
	"烟台":   "YNT",
	"西安":   "SIA",
	"宁波":   "NGB",
	"天津":   "TSH",
	"常德":   "CGD",
	"襄樊":   "XFN",
	"无锡":   "WUX",
	"义乌":   "YIW",
	"三亚":   "SYX",
	"南京":   "NKG",
	"丽江":   "LJG",
	"湛江":   "ZHA",
	"银川":   "INC",
	"常州":   "CZX",
	"宜昌":   "YIH",
	"南阳":   "NNY",
	"宜宾":   "YBP",
	"南宁":   "NNG",
	"西宁":   "XNN",
	"黄山":   "TXN",
	"太原":   "TYN",
	"温州":   "WUZ",
	"晋江":   "JJN",
	"衢州":   "JUZ",
	"合肥":   "HFE",
	"福州":   "FOC",
	"丹东":   "DDG",
	"洛阳":   "LYA",
	"北海":   "BHY",
	"伊宁":   "YIN",
	"梧州":   "WUZ",
	"延安":   "ENY",
	"西昌":   "XIC",
	"威海":   "WEH",
	"苏州":   "SZV",
	"大理":   "DLU",
	"拉萨":   "LXA",
	"吉林":   "JIL",
	"广州":   "CAN",
	"衡阳":   "HNY",
	"黄岩":   "HYN",
	"九江":   "JIU",
	"柳州":   "LZH",
	"汕头":   "SWA",
	"澳门":   "MFM",
	"珠海":   "ZUH",
	"临沧":   "LNJ",
	"铜仁":   "TEN",
	"延吉":   "YNJ",
	"汉中":   "HZG",
	"乌鲁木齐": "URC",
	"齐齐哈尔": "NDG",
	"呼和浩特": "HET",
	"西双版纳": "JHG",
	"龙岩连城": "LCX",
	"西安咸阳": "XIY",
	"哈尔滨":  "HRB",
	"张家界":  "DYG",
	"景德镇":  "JDZ",
	"武夷山":  "WUS",
	"石家庄":  "SJW",
	"佳木斯":  "JUM",
}

var CityCode2 = map[string]string{
	"重庆":   "CKG",
	"成都":   "CTU",
	"武汉":   "WUH",
	"海口":   "HAK",
	"郑州":   "CGO",
	"上海浦东": "PVG",
	"上海虹桥": "SHA",
	"长沙":   "CSX",
	"杭州":   "HGH",
	"厦门":   "XMN",
	"昆明":   "KMG",
}

var CityCode3 = map[string]string{
	"厦门": "XMN",
	"昆明": "KMG",
}

type TripParam struct {
	DepartureCity string
	ArrivalCity   string
	Date          string
}

func ceshi() {
	//fmt.Println("---------------------------")
	//SerachMinPricePlane("NKG", "2019-04-30")	//南京
	//
	//fmt.Println("---------------------------")
	//SerachMinPricePlane("PVG", "2019-04-30")	//广州

	fmt.Println("---------------------------")
	totalCost := SerachMinPricePlane(CityCode["哈尔滨"], "2019-04-30") //哈尔滨
	fmt.Printf("北京和贵阳到%s 机票和最低价格:%v", "哈尔滨", totalCost)

	//fmt.Println("---------------------------")
	//totalCost := SerachMinPricePlane(CityCode["丹东"], "2019-04-30") //哈尔滨
	//fmt.Printf("北京和贵阳到%s 机票和最低价格:%v", "丹东", totalCost)
}

func main() {
	SerchAllCity("2019-04-30")
}

func SerachAirPlane(tripParam TripParam) *downloader.PlaneTaskResult {
	dls := downloader.NewdownloaderService(tripParam.DepartureCity, tripParam.ArrivalCity, tripParam.Date)
	res, err := dls.GetTicketResult()
	if err != nil {
		fmt.Println("get ticket result error", err)
	}

	// 根据机票价格递增排序航班信息
	downloader.FloorPrice(res)
	return res

	//打印所有航班信息
	//for i, ticket := range res.PlaneTasks {
	//	//fmt.Print(ticket.FliInfo, ticket.Prices[0])
	//	flightInfo := ticket.FliInfo
	//	price := ticket.Prices[0]
	//	fmt.Printf("%02d:航班号:%s,航空公司:%s,日期:%s-%s, 始发地:%s,目的地:%s,价格:%v,折扣:%v\n",
	//		i+1,
	//		flightInfo.FlightNumber,
	//		flightInfo.AirlineName,
	//		flightInfo.DepartureDate,
	//		flightInfo.ArrivalDate,
	//		flightInfo.DepInfo.CityName,
	//		flightInfo.ArrInfo.CityName,
	//		price.PrintPrice,
	//		price.Rate)
	//}
	//
	//fmt.Println("\n\ndo over and success!\n")
}

// 搜索某城市到北京和贵阳当日机票最少的两个航班
func SerachMinPricePlane(city, data string) float64 {
	tripParam1 := TripParam{DepartureCity: "BJS", ArrivalCity: city, Date: data} // 北京 -- city
	trips := SerachAirPlane(tripParam1)
	onePlaneTask := downloader.GetOneflight(trips)

	tripParam2 := TripParam{DepartureCity: "KWE", ArrivalCity: city, Date: data} // 贵阳 -- city
	trips2 := SerachAirPlane(tripParam2)
	onePlaneTask2 := downloader.GetOneflight(trips2)

	if len(onePlaneTask.Prices) == 0 {
		fmt.Println("北京 no flight!!")
	}else{
		PrintAirPlane(onePlaneTask)
	}

	if len(onePlaneTask2.Prices) == 0 {
		fmt.Println("贵阳 no flight!!")
	}else{
		PrintAirPlane(onePlaneTask2)
	}
	if len(onePlaneTask.Prices) == 0 || len(onePlaneTask2.Prices) == 0 {
		return 0
	}

	return onePlaneTask.Prices[0].PrintPrice + onePlaneTask2.Prices[0].PrintPrice
}

func PrintAirPlane(onePlaneTask downloader.PlaneTask) {
	flightInfo := onePlaneTask.FliInfo
	price := onePlaneTask.Prices[0]
	fmt.Printf("航班号:%s,航空公司:%s,日期:%s-%s, 行程:%s--%s,价格:%v,折扣:%v\n",
		flightInfo.FlightNumber,
		flightInfo.AirlineName,
		flightInfo.DepartureDate,
		flightInfo.ArrivalDate,
		flightInfo.DepInfo.CityName,
		flightInfo.ArrInfo.CityName,
		price.PrintPrice,
		price.Rate)
}

func SerchAllCity(date string) {
	fn := "I:\\a.txt"
	i := 1
	name := []string{}
	for k, _ := range CityCode {
		name = append(name, k)
	}
	sort.Strings(name)
	buf := bytes.NewBufferString("")
	for _, v := range name {
		totalCost := SerachMinPricePlane(CityCode[v], date)
		fmt.Printf("%02d.北京和贵阳到--%s-- 机票和最低价格:%v\n\n", i, v, totalCost)
		fmt.Fprintf(buf, "%02d.北京和贵阳到--%s-- 机票和最低价格:%v\n", i, v, totalCost)
		i++
		time.Sleep(100 * time.Millisecond)
		err := ioutil.WriteFile(fn, buf.Bytes(), os.ModeAppend)
		if err != nil {
			fmt.Println("Write File err:", err)
		}
	}
	fmt.Println("Serch All City done!!")
}
