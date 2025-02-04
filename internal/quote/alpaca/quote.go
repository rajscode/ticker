package alpaca

import (
    "strings"
    "github.com/go-resty/resty/v2"
    c "github.com/achannarasappa/ticker/v4/internal/common"
)

type Response struct {
    Snapshot struct {
        DailyBar struct {
            Volume int     `json:"v"`
            Low    float64 `json:"l"`
            High   float64 `json:"h"`
            Open   float64 `json:"o"`
            Close  float64 `json:"c"`
        } `json:"dailyBar"`
        LatestQuote struct {
            BidPrice float64 `json:"bp"`
        } `json:"latestQuote"`
        PrevDailyBar struct {
            Close float64 `json:"c"`
        } `json:"prevDailyBar"`
    } `json:"snapshot"`
}

func transformResponse(response Response, symbol string) c.AssetQuote {
    return c.AssetQuote{
        Symbol: symbol,
        QuotePrice: c.QuotePrice{
            Price:          response.Snapshot.LatestQuote.BidPrice,
            PricePrevClose: response.Snapshot.PrevDailyBar.Close,
            PriceOpen:      response.Snapshot.DailyBar.Open,
            PriceDayHigh:   response.Snapshot.DailyBar.High,
            PriceDayLow:    response.Snapshot.DailyBar.Low,
            Change:         response.Snapshot.LatestQuote.BidPrice - response.Snapshot.PrevDailyBar.Close,
        },
        QuoteExtended: c.QuoteExtended{
            Volume: response.Snapshot.DailyBar.Volume,
        },
        QuoteSource: c.QuoteSourceAlpaca,
        Exchange: c.Exchange{
            Name:     "Alpaca",
            IsActive: true,
        },
    }
}

func GetAssetQuotes(client resty.Client, symbols []string) func() []c.AssetQuote {
    return func() []c.AssetQuote {
        quotes := make([]c.AssetQuote, 0)
        for _, symbol := range symbols {
            if strings.HasSuffix(symbol, ".AP") {
                baseSymbol := strings.TrimSuffix(symbol, ".AP")
                res, _ := client.R().
                    SetResult(Response{}).
                    SetHeader("APCA-API-KEY-ID", "keyvalue").
                    SetHeader("APCA-API-SECRET-KEY", "secretvalue").
                    SetQueryParam("symbols", baseSymbol).
                    Get("https://data.alpaca.markets/v1beta1/options/snapshots")
                quotes = append(quotes, transformResponse(*res.Result().(*Response), baseSymbol))
            }
        }
        return quotes
    }
}
