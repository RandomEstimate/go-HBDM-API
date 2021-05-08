package HBDM_API

import (
	goLog "github.com/RandomEstimate/go-log"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type OrderId int
type OrderInfo TradeOrder

type TradeObj struct {
	symbol     string
	accessKey  string
	secretKey  string
	conn       *websocket.Conn
	restClient *http.Client

	pingChan       chan *alivePrivate
	tradeOrderChan chan *TradeOrder
	positionChan   chan *PositionData

	logFile *goLog.FileLog

	orderMap    map[OrderId]OrderInfo
	positionMap struct {
		BuyPrice   float64
		BuyAmount  int
		SellPrice  float64
		SellAmount int
	}

	once       *sync.Once
	m1         *sync.RWMutex
	m2         *sync.RWMutex
	m3         *sync.RWMutex
	certainMap map[OrderId]struct {
		c  chan int
		ts time.Time
	}
}

func NewTradeObj(accessKey, secretKey string, symbol string, logFile *goLog.FileLog) *TradeObj {
	return &TradeObj{
		symbol:         symbol,
		accessKey:      accessKey,
		secretKey:      secretKey,
		conn:           nil,
		logFile:        logFile,
		restClient:     new(http.Client),
		pingChan:       make(chan *alivePrivate, 10),
		tradeOrderChan: make(chan *TradeOrder, 10),
		positionChan:   make(chan *PositionData, 10),
		orderMap:       make(map[OrderId]OrderInfo, 10),
		positionMap: struct {
			BuyPrice   float64
			BuyAmount  int
			SellPrice  float64
			SellAmount int
		}{BuyPrice: 0, BuyAmount: 0, SellPrice: 0, SellAmount: 0},
		once: new(sync.Once),
		m1:   new(sync.RWMutex),
		m2:   new(sync.RWMutex),
		m3:   new(sync.RWMutex),
		certainMap: make(map[OrderId]struct {
			c  chan int
			ts time.Time
		}, 10),
	}

}
