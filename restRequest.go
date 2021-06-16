package HBDM_API

//-------------------------------------------------------------------------------------------------------------------
// 批量下单接口
func (a *TradeObj) BatchOrder(param *OrderMode) (*OrderResponse, error) {
	strUrl := "/swap-api/v1/swap_batchorder"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(OrderResponse))
	return d.(*OrderResponse), err
	//resp, err := a.restClient.Do(req)
	//if err != nil {
	//	return new(OrderResponse), fmt.Errorf("BatchOrder() Request error %v", err)
	//}
	//
	//defer resp.Body.Close()
	//
	//buf, _ := ioutil.ReadAll(resp.Body)
	//d := new(OrderResponse)
	//err = json.Unmarshal(buf, d)
	//if err != nil {
	//	return d, fmt.Errorf("BatchOrder() Unmarshal error %v", err)
	//}
	//return d, nil

}

//-------------------------------------------------------------------------------------------------------------------
// 批量撤单接口
func (a *TradeObj) Cancel(param *CancelMode) (*CancelResponse, error) {
	strUrl := "/swap-api/v1/swap_cancel"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(CancelResponse))
	return d.(*CancelResponse), err

}

//-------------------------------------------------------------------------------------------------------------------
// 全部撤单接口
func (a *TradeObj) CancelAll(param *CancelAllMode) (*CancelResponse, error) {
	strUrl := "/swap-api/v1/swap_cancelall"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(CancelResponse))
	return d.(*CancelResponse), err

}

//-------------------------------------------------------------------------------------------------------------------
// 获取全部订单接口
func (a *TradeObj) SearchOrder(param *SearchOrderMode) (*SearchOrderResponse, error) {

	strUrl := "/swap-api/v1/swap_openorders"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(SearchOrderResponse))
	return d.(*SearchOrderResponse), err
}

//-------------------------------------------------------------------------------------------------------------------
// 获取仓位接口
func (a *TradeObj) Position(param *PositionMode) (*PositionResponse, error) {
	strUrl := "/swap-api/v1/swap_position_info"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(PositionResponse))
	return d.(*PositionResponse), err
}

//-------------------------------------------------------------------------------------------------------------------
// 获取账户权益接口
func (a *TradeObj) MarginPosition(param *MarginMode) (*MarginResponse, error) {

	strUrl := "/swap-api/v1/swap_account_position_info"
	method := "POST"
	strUrl = ApiKeyReady(method, strUrl, a.accessKey, a.secretKey)

	req := ParamReady(strUrl, param, method)
	d, err := RequestDo(a.restClient, req, new(MarginResponse))
	return d.(*MarginResponse), err
}
