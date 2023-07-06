package storage

type Storage interface {
	GetBalanceInfo(login string) ([]byte, error)
	GetWithdrawals(login string) ([]byte, error)
	Withdraw(login string, orderID string, sum float64) error
	GetUserOrders(login string) ([]byte, error)
	LoadOrder(login string, orderID string) error
	Register(login string, password string) error
	Login(login string, password string) error
}
