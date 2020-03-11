package model

import "time"

type Transaction struct {
	TrxID       string    `json:"transaction_id" bson:"transaction_id"`
	DatetimeUTC time.Time `json:"datetime_utc" bson:"datetime_utc"`
	Amount      string    `json:"amount" bson:"amount"`
}

type TransactionSnapshot struct {
	HourlyDateTimeUTC time.Time `json:"hourly_datetime_utc" bson:"hourly_datetime_utc"`
	AccumulatedAmount string    `json:"accumulated_amount" bson:"accumulated_amount"`
}

type AddBalanceRequest struct {
	DateTime string  `json:"datetime"`
	Amount   float64 `json:"amount"`
}

type AccumulateBalance struct {
	Key          string `json:"key" bson:"key"`
	TotalBalance string `json:"total_balance" bson:"total_balance"`
}

type BalanceHourlyRange struct {
	Datetime string `json:"datetime"`
	Amount   string `json:"amount"`
}
