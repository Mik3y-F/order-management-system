package pkg

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusPaid       OrderStatus = "paid"
)
