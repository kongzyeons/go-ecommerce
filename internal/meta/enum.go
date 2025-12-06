package meta

const (
	ORDER_STATUS_EDITABLE  string = "editable"
	ORDER_STATUS_CONFIRM   string = "confirm"
	ORDER_STATUS_SHIPPING  string = "shipping"
	ORDER_STATUS_COMPLETED string = "completed"
	ORDER_STATUS_CANCEL    string = "cancel"
)

func GetOrderStatus() map[string]bool {
	return map[string]bool{
		ORDER_STATUS_EDITABLE:  true,
		ORDER_STATUS_CONFIRM:   true,
		ORDER_STATUS_SHIPPING:  true,
		ORDER_STATUS_COMPLETED: true,
		ORDER_STATUS_CANCEL:    true,
	}
}
