package models

type TradeData struct {
	InstId  string `json:"instId"`
	Side    string `json:"side"`
	Sz      string `json:"sz"`
	Px      string `json:"px"`
	TradeId string `json:"tradeId"`
	Ts      string `json:"ts"`
}