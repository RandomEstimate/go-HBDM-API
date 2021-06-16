package HBDM_API

import "strings"

//---------------------------------------------Http_Request_Structure------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------
// 下单结构体
type OrderMode struct {
	OrdersData []OrderData `json:"orders_data"`
}

type OrderData struct {
	ContractCode   string  `json:"contract_code"`
	ClientOrderId  int     `json:"client_order_id,omitempty"`
	Price          float64 `json:"price,omitempty"`
	Volume         float64 `json:"volume"`
	Direction      string  `json:"direction"`
	Offset         string  `json:"offset"`
	LeverRate      int     `json:"lever_rate"`
	OrderPriceType string  `json:"order_price_type,omitempty"`
}

//-------------------------------------------------------------------------------------------------------------------
// 批量撤单结构体
type CancelMode struct {
	// 订单ID(多个订单ID中间以","分隔,一次最多允许撤消10个订单)
	OrderId string `json:"order_id,omitempty"`
	// 客户订单ID(多个订单ID中间以","分隔,一次最多允许撤消10个订单)
	ClientOrderId string `json:"client_order_id,omitempty"`
	// 合约代码,支持大小写,"BTC-USD"
	ContractCode string `json:"contract_code"`
}

//-------------------------------------------------------------------------------------------------------------------
// 全部撤单结构体
type CancelAllMode struct {
	ContractCode string `json:"contract_code"`
	// 买卖方向（不填默认全部） "buy":买 "sell":卖
	Direction string `json:"direction,omitempty"`
	// 开平方向（不填默认全部） "open":开 "close":平
	Offset string `json:"offset,omitempty"`
}

//-------------------------------------------------------------------------------------------------------------------
// 获取全部订单结构体
type SearchOrderMode struct {
	ContractCode string `json:"contract_code"`
	// 页码，不填默认第1页
	PageIndex int `json:"page_index,omitempty"`
	// 不填默认20，不得多于50
	PageSize int `json:"page_size,omitempty"`
	// 排序字段，不填默认按创建时间倒序 “created_at”(按照创建时间倒序)，“update_time”(按照更新时间倒序)
	SortBy string `json:"sort_by,omitempty"`
	// 	0:全部,1:买入 开多,2: 卖出开空,3: 买入平空,4: 卖出平多。
	TradeType int `json:"trade_type,omitempty"`
}

//-------------------------------------------------------------------------------------------------------------------
// 获取仓位结构体
type PositionMode struct {
	ContractCode string `json:"contract_code"`
}

//-------------------------------------------------------------------------------------------------------------------
// 账户权益结构体
type MarginMode struct {
	ContractCode string `json:"contract_code"`
}


//-------------------------------------------------------------------------------------------------------------------

//---------------------------------------------Http_Response_Structure-----------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------
// 批量下单返回结构体
type OrderResponse struct {
	Status string `json:"status"`
	Data   struct {
		Errors []struct {
			Index   int    `json:"index"`
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		}
		Success []struct {
			OrderId    int    `json:"order_id"`
			Index      int    `json:"index"`
			OrderIdStr string `json:"order_id_str"`
		}
	} `json:"data"`
	Ts int `json:"ts"`
}

//-------------------------------------------------------------------------------------------------------------------
// 批量撤单返回结构体 (同 撤销所有订单返回结构体)
type CancelResponse struct {
	Status string `json:"status"`
	Data   struct {
		Error []struct {
			OrderId int    `json:"order_id"`
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		}
		Successes string `json:"successes"`
	} `json:"data"`
	Ts int `json:"ts"`
}

//-------------------------------------------------------------------------------------------------------------------
// 获取全部订单返回结构体
type SearchOrderResponse struct {
	Status string `json:"status"`
	Data   struct {
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		TotalSize   int `json:"total_size"`
		Orders      []struct {
			Symbol          string  `json:"symbol"`
			ContractCode    string  `json:"contract_code"`
			Volume          int     `json:"volume"`
			Price           float64 `json:"price"`
			OrderPriceType  string  `json:"order_price_type"`
			OrderType       int     `json:"order_type"`
			Direction       string  `json:"direction"`
			Offset          string  `json:"offset"`
			LeverRate       int     `json:"lever_rate"`
			OrderId         int     `json:"order_id"`
			ClientOrderId   string  `json:"client_order_id"`
			CreatedAt       int     `json:"created_at"`
			TradeVolume     float64 `json:"trade_volume"`
			TradeTurnover   float64 `json:"trade_turnover"`
			Fee             float64 `json:"fee"`
			TradeAvgPrice   float64 `json:"trade_avg_price"`
			MarginFrozen    float64 `json:"margin_frozen"`
			Profit          float64 `json:"profit"`
			Status          float64 `json:"status"`
			OrderSource     string  `json:"order_source"`
			OrderIdStr      string  `json:"order_id_str"`
			FeeAsset        string  `json:"fee_asset"`
			LiquidationType string  `json:"liquidation_type"`
			CanceledAt      float64 `json:"canceled_at"`
			IsTpsl          int     `json:"is_tpsl"`
			UpdateTime      float64 `json:"update_time"`
			RealProfit      float64 `json:"real_profit"`
		}
	} `json:"data"`
	Ts int `json:"ts"`
}

//-------------------------------------------------------------------------------------------------------------------
// 获取仓位返回结构体
type PositionResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		Volume         float64 `json:"volume"`
		Available      float64 `json:"available"`
		Frozen         float64 `json:"frozen"`
		CostOPen       float64 `json:"cost_o_pen"`
		CostHold       float64 `json:"cost_hold"`
		ProfitUnreal   float64 `json:"profit_unreal"`
		ProfitRate     float64 `json:"profit_rate"`
		LeverRate      int     `json:"lever_rate"`
		PositionMargin float64 `json:"position_margin"`
		Direction      string  `json:"direction"`
		Profit         float64 `json:"profit"`
		LastPrice      float64 `json:"last_price"`
	} `json:"data"`
	Ts int `json:"ts"`
}

//-------------------------------------------------------------------------------------------------------------------
// 账户权益返回结构体
type MarginResponse struct {
	Status string `json:"status"`
	Data   []struct {
		Symbol         string  `json:"symbol"`
		ContractCode   string  `json:"contract_code"`
		MarginBalance  float64 `json:"margin_balance"`
		MarginStatic   float64 `json:"margin_static"`
		MarginPosition float64 `json:"margin_position"`
		// 其他部分目前不需要
	}
}

//---------------------------------------------Conn_Structure-------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------
// ping-pong 结构体
type alivePrivate struct {
	Op string `json:"op"`
	Ts string `json:"ts"`
}

func (a *alivePrivate) Valid() bool {
	if a.Op == "ping" {
		return true
	} else {
		return false
	}
}

//-------------------------------------------------------------------------------------------------------------------
//// OrderInfo
type TradeOrder2 struct {
	TradeVolume float64 `json:"trade_volume"`
	TradePrice  float64 `json:"trade_price"`
	TradeFee    float64 `json:"trade_fee"`
	CreatedAt   float64 `json:"created_at"`
}

type TradeOrder struct {
	Topic         string        `json:"topic"`
	Symbol        string        `json:"symbol"`
	Direction     string        `json:"direction"`
	Offset        string        `json:"offset"`
	Status        int           `json:"status"`
	OrderId       int           `json:"order_id"`
	OrderIdStr    string        `json:"order_id_str"`
	ClientOrderId int           `json:"client_order_id"`
	OrderType     int           `json:"order_type"`
	TradeVolume   float64       `json:"trade_volume"`
	Price         float64       `json:"price"`
	Volume        float64       `json:"volume"`
	Trade         []TradeOrder2 `json:"trade"`
}

func (a *TradeOrder) Valid() bool {
	if strings.Index(a.Topic, "orders") != -1 {
		return true
	} else {
		return false
	}
}

//-------------------------------------------------------------------------------------------------------------------
// positionInfo
type PositionData2 struct {
	Symbol       string  `json:"symbol"`
	Volume       float64 `json:"volume"`
	Available    float64 `json:"available"`
	CostOpen     float64 `json:"cost_open"`
	Direction    string  `json:"direction"`
	ContractCode string  `json:"contract_code"`
}

type PositionData struct {
	Topic string          `json:"topic"`
	Ts    float64         `json:"ts"`
	Event string          `json:"event"`
	Data  []PositionData2 `json:"data"`
}

func (a *PositionData) Valid() bool {
	if strings.Index(a.Topic, "position") != -1 {
		return true
	} else {
		return false
	}
}

//-------------------------------------------------------------------------------------------------------------------
