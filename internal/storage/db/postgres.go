package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
	"yandex_gophermart/internal/models"
)

func (m *Manager) GetBalanceInfo(login string) ([]byte, error) {
	userBalance, err := m.getUserBalance(login)
	if err != nil {
		return nil, fmt.Errorf("error while getting current user userBalance: %w", err)
	}

	getUserWithdrawn := `SELECT SUM(amount) AS withdrawn FROM withdraw WHERE login = $1`

	row := m.db.QueryRow(getUserWithdrawn, login)

	var userWithdrawn sql.NullFloat64

	if err = row.Scan(&userWithdrawn); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("error while getting user withdrawn info: %w", err)
		}

		userWithdrawn = sql.NullFloat64{
			Float64: 0,
			Valid:   true,
		}
	}

	info := models.BalanceInfo{
		Withdrawn: userWithdrawn.Float64,
		Current:   userBalance,
	}

	result, err := json.Marshal(info)

	if err != nil {
		return nil, fmt.Errorf("error while marshalling user balance info: %w", err)
	}

	return result, nil
}

func (m *Manager) GetWithdrawals(login string) ([]byte, error) {
	getUserWithdrawals := `SELECT order_id, amount, processed_at FROM withdraw WHERE login = $1 ORDER BY processed_at`

	rows, err := m.db.Query(getUserWithdrawals, login)
	if err != nil {
		return nil, fmt.Errorf("error while searching for userWithdrawals: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			_ = fmt.Errorf("error rows close is: %s", err)
		}

		err = rows.Err()
		if err != nil {
			_ = fmt.Errorf("error rows Err is: %s", err)
		}
	}()

	userWithdrawals := make([]models.WithdrawInfo, 0)
	for rows.Next() {
		var orderID string
		var amount float64
		var processedAt time.Time

		if err = rows.Scan(&orderID, &amount, &processedAt); err != nil {
			return nil, fmt.Errorf("error while scanning rows from userWithdrawals: %w", err)
		}

		userWithdrawals = append(userWithdrawals, models.WithdrawInfo{
			OrderID:     orderID,
			ProcessedAt: &processedAt,
			Amount:      amount,
		})
	}

	if len(userWithdrawals) == 0 {
		return nil, ErrNoData
	}

	result, err := json.Marshal(userWithdrawals)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling user withdrawals info: %w", err)
	}

	return result, nil
}

func (m *Manager) Withdraw(login string, orderID string, sum float64) error {
	userBalance, err := m.getUserBalance(login)
	if err != nil {
		return fmt.Errorf("error while checking user userBalance: %w", err)
	}

	if userBalance < sum {
		return ErrInsufficientBalance
	}

	withdraw := `INSERT INTO withdraw VALUES ($1, $2, now(), $3)`

	if _, err = m.db.Exec(withdraw, login, orderID, sum); err != nil {
		return fmt.Errorf("error while trying to withdraw: %w", err)
	}

	return nil
}

func (m *Manager) GetUserOrders(login string) ([]byte, error) {
	getUserOrders := `SELECT order_id, status, accrual, uploaded_at FROM orders WHERE login = $1`

	rows, err := m.db.Query(getUserOrders, login)

	if err != nil {
		return nil, fmt.Errorf("error while getting orders from db for user %q: %w", login, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			_ = fmt.Errorf("error rows close is: %s", err)
		}

		err = rows.Err()
		if err != nil {
			_ = fmt.Errorf("error rows Err is: %s", err)
		}
	}()

	userOrders := make([]models.OrderInfo, 0)

	for rows.Next() {
		var (
			orderID    string
			status     models.OrderStatus
			accrual    sql.NullFloat64
			uploadedAt time.Time
		)

		if err = rows.Scan(&orderID, &status, &accrual, &uploadedAt); err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		userOrders = append(userOrders, models.OrderInfo{
			OrderID:   orderID,
			Accrual:   accrual.Float64,
			CreatedAt: &uploadedAt,
			Status:    status,
		})
	}

	if len(userOrders) == 0 {
		return nil, ErrNoData
	}

	result, err := json.Marshal(userOrders)
	if err != nil {
		return nil, fmt.Errorf("error while marshalling user orders info: %w", err)
	}

	return result, nil
}

func (m *Manager) GetAllOrders() ([]string, error) {
	getAllOrders := `SELECT order_id FROM orders WHERE status != $1`

	rows, err := m.db.Query(getAllOrders, "PROCESSED")

	if err != nil {
		return nil, fmt.Errorf("error while getting all orders from db: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			_ = fmt.Errorf("error rows close is: %s", err)
		}

		err = rows.Err()
		if err != nil {
			_ = fmt.Errorf("error rows Err is: %s", err)
		}
	}()

	orders := make([]string, 0)

	for rows.Next() {
		var orderID string

		if err = rows.Scan(&orderID); err != nil {
			return nil, fmt.Errorf("error while scanning rows: %w", err)
		}

		orders = append(orders, orderID)
	}

	return orders, nil
}

func (m *Manager) UpdateOrderInfo(orderInfo *models.OrderInfo) error {
	updateOrderInfo := `UPDATE orders SET status = $1, accrual = $2 WHERE order_id = $3`

	if _, err := m.db.Exec(updateOrderInfo, string(orderInfo.Status), orderInfo.Accrual, orderInfo.Order); err != nil {
		return fmt.Errorf("error while updating order info: %w", err)
	}

	return nil
}

func (m *Manager) LoadOrder(login string, orderID string) error {
	getOrderByID := `SELECT login FROM orders WHERE order_id = $1`

	row := m.db.QueryRow(getOrderByID, orderID)

	var userName string

	err := row.Scan(&userName)

	switch err {
	case sql.ErrNoRows:
		loadOrderQuery := `INSERT INTO orders VALUES ($1, $2, now(), $3, $4)`

		if _, err = m.db.Exec(loadOrderQuery, orderID, login, models.OrderStatus("NEW"), 0); err != nil {
			return fmt.Errorf("error while loading order %s: %w", orderID, err)
		}

		return nil
	case nil:
		if userName == login {
			return ErrCreatedBySameUser
		}

		return ErrCreatedDiffUser
	default:
		return fmt.Errorf("error while scanning rows: %w", err)
	}
}

func (m *Manager) Register(login string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("this password is not allowed: %w", err)
	}

	registerUser := `INSERT INTO registered_users VALUES ($1, $2)`

	if _, err = m.db.Exec(registerUser, login, hash); err != nil {
		dublicateKeyErr := ErrDublicateKey{Key: "registered_users_pkey"}
		if err.Error() == dublicateKeyErr.Error() {
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("error while executing register user query: %w", err)
	}

	return nil
}

func (m *Manager) Login(login string, password string) error {
	getRegisteredUser := `SELECT login, password FROM registered_users WHERE login = $1`

	rows, err := m.db.Query(getRegisteredUser, login)
	if err != nil {
		return fmt.Errorf("error while executing search query: %w", err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			_ = fmt.Errorf("error rows close is: %s", err)
		}

		err = rows.Err()
		if err != nil {
			_ = fmt.Errorf("error rows Err is: %s", err)
		}
	}()

	for rows.Next() {
		var loginFromDB, passwordFromDB string

		if err = rows.Scan(&loginFromDB, &passwordFromDB); err != nil {
			return fmt.Errorf("error while scanning rows: %w", err)
		}

		if loginFromDB == login {
			if err = bcrypt.CompareHashAndPassword([]byte(passwordFromDB), []byte(password)); err != nil {
				return ErrInvalidCredentials
			}

			return nil
		}
	}

	return ErrNoSuchUser
}

func (m *Manager) getUserBalance(login string) (float64, error) {
	getUserBalance := `SELECT COALESCE(SUM(accrual), 0) - COALESCE(SUM(amount), 0) AS balance FROM orders o LEFT JOIN withdraw w ON o.login = w.login WHERE o.login = $1 GROUP BY o.login;`

	row := m.db.QueryRow(getUserBalance, login)

	var balance sql.NullFloat64

	if err := row.Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("error while getting user balance: %w", err)
	}

	return balance.Float64, nil
}

func (m *Manager) init(ctx context.Context) error {
	createRegisteredQuery := `CREATE TABLE IF NOT EXISTS registered_users (login text PRIMARY KEY, password TEXT)`
	if _, err := m.db.ExecContext(ctx, createRegisteredQuery); err != nil {
		return fmt.Errorf("error while trying to create table with registered users: %w", err)
	}

	createOrdersQuery := `CREATE TABLE IF NOT EXISTS orders (order_id TEXT UNIQUE, login TEXT, uploaded_at TIMESTAMP WITH TIME ZONE, status TEXT, accrual DOUBLE PRECISION, PRIMARY KEY(order_id))`
	if _, err := m.db.ExecContext(ctx, createOrdersQuery); err != nil {
		return fmt.Errorf("error while trying to create table with orders: %w", err)
	}

	createWithdrawQuery := `CREATE TABLE IF NOT EXISTS withdraw (login TEXT, order_id TEXT UNIQUE, processed_at TIMESTAMP WITH TIME ZONE, amount DOUBLE PRECISION, PRIMARY KEY(login, order_id))`
	if _, err := m.db.ExecContext(ctx, createWithdrawQuery); err != nil {
		return fmt.Errorf("error while trying to create table with orders: %w", err)
	}

	return nil
}

func NewPostgres(ctx context.Context, db *sql.DB) (*Manager, error) {

	m := Manager{
		db: db,
	}

	if err := m.init(ctx); err != nil {
		return nil, err
	}

	return &m, nil
}

type Manager struct {
	db *sql.DB
}
