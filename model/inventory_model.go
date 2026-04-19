package model

type InventoryLevel string

const (
	High   InventoryLevel = "high"
	Medium InventoryLevel = "medium"
	Low    InventoryLevel = "low"
	Out    InventoryLevel = "out"
)
