package models

type TripSearchParam struct {
	Page     int
	PageSize int
}

func NewTripSearchParam(page int, pageSize int) *TripSearchParam {
	return &TripSearchParam{
		Page:     page,
		PageSize: pageSize,
	}
}
