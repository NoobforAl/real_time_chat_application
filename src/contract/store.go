package contract

type Store interface {
	Cache

	StoreUser
	StoreRoom
	StoreMessage
	StoreNotifications

	AuthenticationService

	BrokerRoom
	BrokerMessage
	BrokerNotification
}
