package goanda

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"
)

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/pricing-ep/

type Pricings struct {
	Prices []struct {
		Asks []struct {
			Liquidity int    `json:"liquidity"`
			Price     string `json:"price"`
		} `json:"asks"`
		Bids []struct {
			Liquidity int    `json:"liquidity"`
			Price     string `json:"price"`
		} `json:"bids"`
		CloseoutAsk                string `json:"closeoutAsk"`
		CloseoutBid                string `json:"closeoutBid"`
		Instrument                 string `json:"instrument"`
		QuoteHomeConversionFactors struct {
			NegativeUnits string `json:"negativeUnits"`
			PositiveUnits string `json:"positiveUnits"`
		} `json:"quoteHomeConversionFactors"`
		Status         string    `json:"status"`
		Time           time.Time `json:"time"`
		UnitsAvailable struct {
			Default struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"default"`
			OpenOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"openOnly"`
			ReduceFirst struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceFirst"`
			ReduceOnly struct {
				Long  string `json:"long"`
				Short string `json:"short"`
			} `json:"reduceOnly"`
		} `json:"unitsAvailable"`
	} `json:"prices"`
}

func (c *OandaConnection) GetPricingForInstruments(instruments []string) (Pricings, error) {
	instrumentString := strings.Join(instruments, ",")
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + url.QueryEscape(instrumentString)

	data := Pricings{}
	response, err := c.Get(endpoint)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(response, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
