package types

import "time"

type GetDateandCost struct {
	Time time.Time
	Cost float64
}

type GetGraphDataResponse struct {
	Time string
	Cost float64
}

type GetServiceCostMinMAxTotalAvgResponse struct {
	TotalCost   string
	AverageCost string
	MinimumCost string
	MaximumCost string
}
