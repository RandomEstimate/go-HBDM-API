package HBDM_API

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"strings"
	"time"
)

const wsUrl = "wss://api.btcgateway.pro/swap-notification"

func (a *TradeObj) login() error {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = a.accessKey
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp
	hostName := "api.btcgateway.pro"
	Signature := CreateSign(mapParams2Sign, "GET", hostName, "/swap-notification", a.secretKey)
	type tmp struct {
		Op               string `json:"op"`
		Type             string `json:"type"`
		AccessKeyId      string `json:"AccessKeyId"`
		SignatureMethod  string `json:"SignatureMethod"`
		SignatureVersion string `json:"SignatureVersion"`
		Timestamp        string `json:"Timestamp"`
		Signature        string `json:"Signature"`
	}
	sendMessage := tmp{
		Op:               "auth",
		Type:             "api",
		AccessKeyId:      a.accessKey,
		SignatureMethod:  "HmacSHA256",
		SignatureVersion: "2",
		Timestamp:        timestamp,
		Signature:        Signature,
	}

	sendMessageBuf, _ := json.Marshal(sendMessage)

	err := a.conn.WriteMessage(websocket.TextMessage, sendMessageBuf)
	if err != nil {
		return err
	}

	return nil

}

func (a *TradeObj) sub() error {

	err := a.conn.WriteMessage(websocket.TextMessage, []byte("{\"op\": \"sub\", \"cid\": \" \", \"topic\": \"orders."+strings.ToLower(a.symbol)+"\"}"))
	err_ := a.conn.WriteMessage(websocket.TextMessage, []byte("{\"op\": \"sub\", \"cid\": \" \",\"topic\": \"positions."+a.symbol+"\"}"))
	if err == nil && err_ == nil {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("err= %v and %v", err, err_))
}

func (a *TradeObj) deal() {
	scanCertainMap := time.NewTicker(time.Minute)
	for {
		select {
		case data := <-a.tradeOrderChan:
			a.handleTradeOrder(data)
		case data := <-a.positionChan:
			a.handlePosition(data)
		case data := <-a.pingChan:
			a.conn.WriteMessage(websocket.TextMessage, []byte("{\"op\":\"pong\",\"ts\":\""+data.Ts+"\"}"))
		case <-scanCertainMap.C:
			a.handleScan()
		}
	}
}

func (a *TradeObj) handleTradeOrder(d *TradeOrder) {
	a.m1.Lock()
	a.m3.Lock()
	defer a.m1.Unlock()
	defer a.m3.Unlock()

	if d.Status == 3 || d.Status == 5 || d.Status == 6 || d.Status == 7 {
		if d.Status == 3 {
			a.orderMap[OrderId(d.OrderId)] = OrderInfo(*d)
			if s, ok := a.certainMap[OrderId(d.ClientOrderId)]; ok {
				s.c <- d.ClientOrderId
				a.delCertainMap(OrderId(d.ClientOrderId))
			}
		}

		if d.Status == 5 || d.Status == 6 || d.Status == 7 {
			delete(a.orderMap, OrderId(d.OrderId))
			if s, ok := a.certainMap[OrderId(d.OrderId)]; ok {
				s.c <- d.OrderId
				a.delCertainMap(OrderId(d.OrderId))
			}
		}


	}

}

func (a *TradeObj) handleScan() {
	a.m3.Lock()
	defer a.m3.Unlock()
	for k, v := range a.certainMap {
		if time.Since(v.ts) > time.Minute {
			a.delCertainMap(k)
		}
	}
}

func (a *TradeObj) handlePosition(d *PositionData) {
	a.m2.Lock()
	defer a.m2.Unlock()

	for _, v := range d.Data {
		if v.ContractCode == a.symbol {
			if v.Direction == "buy" {
				a.positionMap.BuyPrice = v.CostOpen
				a.positionMap.BuyAmount = int(v.Volume)
			} else if v.Direction == "sell" {
				a.positionMap.SellPrice = v.CostOpen
				a.positionMap.SellAmount = int(v.Volume)
			}

		}
	}

}

func (a *TradeObj) loadMessage(message string) interface{} {
	tmp := TradeOrder{}
	err := json.Unmarshal([]byte(message), &tmp)
	if err == nil && tmp.Valid() {
		return &tmp
	}

	tmp2 := alivePrivate{}
	err = json.Unmarshal([]byte(message), &tmp2)
	if err == nil && tmp2.Valid() {
		return &tmp2
	}

	tmp3 := PositionData{}
	err = json.Unmarshal([]byte(message), &tmp3)
	if err == nil && tmp3.Valid() {
		return &tmp3
	}

	return nil
}

func (a *TradeObj) StartConn() {

	defer func() {
		if err := recover(); err != nil {
			a.logFile.I("startConn() recover err= %v", err)
			a.conn.Close()
			go a.StartConn()
		}
	}()

	var errTime = 0
	for {
		c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			a.logFile.E(fmt.Sprintf("conn fail err: %v", err))
		} else {
			a.conn = c
			break
		}
		errTime++
		if errTime == 10 {
			a.logFile.E("conn fail 10 times err: %v", err)
			panic(err)
		}
	}

	logErr := a.login()
	if logErr != nil {

		a.logFile.E("login err:%v", logErr)
	}


	subErr := a.sub()
	if subErr != nil {

		a.logFile.E("sub err:%v", subErr)
	}

	a.once.Do(func() {
		go a.deal()
	})

	for {
		_, buf, err := a.conn.ReadMessage()
		if err != nil {
			a.logFile.E("ReadMessage() err=%v", err)
			panic(err)
		}
		buf1 := bytes.NewBuffer(buf)
		GzipReader, err := gzip.NewReader(buf1)
		if err != nil {
			a.logFile.E("GzipReader() err=%v", err)
		}
		messageBuf, _ := ioutil.ReadAll(GzipReader)
		parse := a.loadMessage(string(messageBuf))

		switch parse.(type) {
		case *TradeOrder:
			a.tradeOrderChan <- parse.(*TradeOrder)
		case *PositionData:
			a.positionChan <- parse.(*PositionData)
		case *alivePrivate:
			a.pingChan <- parse.(*alivePrivate)
		}

	}

}
