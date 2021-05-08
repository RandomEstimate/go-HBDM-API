package HBDM_API

import (
	"fmt"
	goLog "github.com/RandomEstimate/go-log"
	"strings"
	"testing"
	"time"
)

const (
	accessKey = "44917a14-0684b378-rbr45t6yr4-ae79f"
	secretKey = "09ac5548-edc081d4-49257a73-7cf8a"
)

func TestTradeObj_SafeTrade(t *testing.T) {
	logFile := goLog.NewFileLog("", "log.log", goLog.TRACE)
	logFile.Start()
	defer logFile.Close()
	logFile.I("Start logfile")
	obj := NewTradeObj(accessKey, secretKey, "SUSHI-USD", logFile)
	go obj.StartConn()
	time.Sleep(time.Second)
	o := &OrderMode{
		OrdersData: []OrderData{
			{
				ContractCode:   "SUSHI-USD",
				ClientOrderId:  int(time.Now().UnixNano()) + 1,
				Price:          14,
				Volume:         1,
				Direction:      "buy",
				Offset:         "open",
				LeverRate:      5,
				OrderPriceType: "post_only",
			},
			{
				ContractCode:   "SUSHI-USD",
				ClientOrderId:  int(time.Now().UnixNano()) + 2,
				Price:          14.5,
				Volume:         1,
				Direction:      "buy",
				Offset:         "open",
				LeverRate:      5,
				OrderPriceType: "post_only",
			},
		},
	}
	err := obj.SafeTrade(o)
	fmt.Println(err)
	time.Sleep(time.Second * 5)

}

func TestTradeObj_SafeCancel(t *testing.T) {
	logFile := goLog.NewFileLog("", "log.log", goLog.TRACE)
	logFile.Start()
	defer logFile.Close()
	logFile.I("Start logfile")
	obj := NewTradeObj(accessKey, secretKey, "SUSHI-USD", logFile)
	go obj.StartConn()
	time.Sleep(time.Second)

	o := &OrderMode{
		OrdersData: []OrderData{
			{
				ContractCode:   "SUSHI-USD",
				ClientOrderId:  int(time.Now().UnixNano()) + 1,
				Price:          14,
				Volume:         1,
				Direction:      "buy",
				Offset:         "open",
				LeverRate:      5,
				OrderPriceType: "post_only",
			},
			{
				ContractCode:   "SUSHI-USD",
				ClientOrderId:  int(time.Now().UnixNano()) + 2,
				Price:          14.5,
				Volume:         1,
				Direction:      "buy",
				Offset:         "open",
				LeverRate:      5,
				OrderPriceType: "post_only",
			},
		},
	}
	err := obj.SafeTrade(o)
	fmt.Println(err)
	time.Sleep(time.Second * 5)

	map_ := obj.GetOrderMap()

	idList := make([]string, 0, 10)
	for _, v := range map_ {
		idList = append(idList, fmt.Sprint(v.OrderId))
	}

	o1 := &CancelMode{
		OrderId:      strings.Join(idList, ","),
		ContractCode: "SUSHI-USD",
	}

	err = obj.SafeCancel(o1)
	fmt.Println(err)
	time.Sleep(time.Second * 5)
}

func TestTradeObj_GetOrderMap(t *testing.T) {
	logFile := goLog.NewFileLog("", "log.log", goLog.TRACE)
	logFile.Start()
	defer logFile.Close()
	logFile.I("Start logfile")
	obj := NewTradeObj(accessKey, secretKey, "SUSHI-USD", logFile)
	go obj.StartConn()

	for i := 0; i <= 30; i++ {
		time.Sleep(time.Second * 3)
		fmt.Println(obj.GetOrderMap())
	}

}
