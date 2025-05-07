package services

type Notification struct {
	userIds []int
	message []byte
}

type ServiceResponse struct {
	Message      []byte
	Notification *Notification
}
