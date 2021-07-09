package HBDM_API

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

func (a *TradeObj) addCertainMap(id OrderId, c chan int) {
	a.m3.Lock()
	defer a.m3.Unlock()
	a.certainMap[id] = struct {
		c  chan int
		ts time.Time
	}{c: c, ts: time.Now()}
}

func (a *TradeObj) delCertainMap(id OrderId) {
	delete(a.certainMap, id)
}

func (a *TradeObj) GetOrderMap() map[int]OrderInfo {
	a.m1.Lock()
	defer a.m1.Unlock()
	t := make(map[int]OrderInfo, len(a.orderMap))
	for k, v := range a.orderMap {
		t[int(k)] = v
	}
	return t

}

func (a *TradeObj) GetPosition() (BuyPrice float64, BuyAmount int, SellPrice float64, SellAmount int) {
	a.m2.Lock()
	defer a.m2.Unlock()
	return a.positionMap.BuyPrice, a.positionMap.BuyAmount, a.positionMap.SellPrice, a.positionMap.SellAmount
}

func (a *TradeObj) SafeTrade(param *OrderMode) error {
	a.logFile.I("Start SafeTrade time=%v", time.Now())
	id := make(map[OrderId]string, 10)
	idOrder := make([]OrderId, 0, 10)
	c := make(chan int, 10)

	for _, v := range param.OrdersData {
		id[OrderId(v.ClientOrderId)] = ""
		idOrder = append(idOrder, OrderId(v.ClientOrderId))
		a.addCertainMap(OrderId(v.ClientOrderId), c)
	}

	m := sync.RWMutex{}
	c2 := make(chan string, 1)

	errChan := make(chan *OrderResponse, 1)

	go func() {
		resp, err := a.BatchOrder(param)
		if err != nil {
			a.logFile.E("SafeTrade() err=", err)
		}
		m.Lock()
		for _, v := range resp.Data.Errors {
			delete(id, idOrder[v.Index-1])
		}
		m.Unlock()

		// 返回错误信息
		errChan <- resp

	}()

	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case d := <-c:
				m.Lock()
				delete(id, OrderId(d))
				m.Unlock()
				if len(id) == 0 {
					c2 <- "success"
					goto Exit
				}
			case <-t.C:
				c2 <- "time out"
				goto Exit
			}
		}
	Exit:
	}()

	select {
	case d := <-c2:
		switch d {
		case "success":
			a.logFile.I("Trade success time=%v", time.Now())
			return nil
		case "time out":

			// time out 记录订单信息 以及request返回信息
			var errInfo *OrderResponse
			select {
			case errInfo = <-errChan:
			default:
			}
			// 序列化传入log参数
			orderInfo2string, _ := json.Marshal(param)
			errInfo2string, _ := json.Marshal(errInfo)
			a.logFile.E("Trade time out time=%v "+
				"orderInfo=%v errInfo=%v", time.Now(), string(orderInfo2string), errInfo2string)

			return fmt.Errorf("time out")
		}
	}
	return fmt.Errorf("miss err")

}

func (a *TradeObj) SafeCancel(param *CancelMode) error {
	a.logFile.I("Start SafeTrade time=%v", time.Now())
	id := make(map[OrderId]string, 10)
	c := make(chan int, 10)

	for _, v := range strings.Split(param.OrderId, ",") {
		OrderId_, _ := strconv.ParseInt(v, 10, 64)
		id[OrderId(OrderId_)] = ""
		a.addCertainMap(OrderId(OrderId_), c)
	}

	m := sync.RWMutex{}
	c2 := make(chan string, 1)

	go func() {
		resp, err := a.Cancel(param)
		if err != nil {
			a.logFile.E("SafeCancel() err=", err)
		}
		m.Lock()
		for _, v := range resp.Data.Error {
			delete(id, OrderId(v.OrderId))
		}
		m.Unlock()

	}()

	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case d := <-c:
				m.Lock()
				delete(id, OrderId(d))
				m.Unlock()
				if len(id) == 0 {
					c2 <- "success"
					goto Exit
				}
			case <-t.C:
				c2 <- "time out"
				goto Exit
			}
		}
	Exit:
	}()
	select {
	case d := <-c2:
		switch d {
		case "success":
			a.logFile.I("Cancel success time=%v", time.Now())
			return nil
		case "time out":
			a.logFile.I("Cancel time out time=%v", time.Now())
			return fmt.Errorf("time out")
		}
	}
	return fmt.Errorf("miss err")

}

func (a *TradeObj) RepairOrderMap() error {
	o := &SearchOrderMode{
		ContractCode: a.symbol,
	}
	order, err := a.SearchOrder(o)
	if err != nil {
		a.logFile.E("RepairOrderMap() err:", err)
		return err
	}

	a.m1.Lock()
	defer a.m1.Unlock()
	for k, _ := range a.orderMap {
		delete(a.orderMap, k)
	}
	for _, v := range order.Data.Orders {
		o_ := OrderInfo(TradeOrder{
			Symbol:      v.Symbol,
			Direction:   v.Direction,
			Offset:      v.Offset,
			Status:      int(v.Status),
			OrderId:     v.OrderId,
			OrderIdStr:  fmt.Sprint(v.OrderId),
			OrderType:   v.OrderType,
			TradeVolume: v.TradeVolume,
		})
		a.orderMap[OrderId(o_.OrderId)] = o_
	}
	return nil
}
