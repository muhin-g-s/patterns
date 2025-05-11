package paypal

import (
	"context"
	"errors"
	"factory-method/internal/payment/gateway"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryTransactionStore struct {
	transactions map[string]*gateway.TransactionStatus
	mu           sync.RWMutex
}

func NewTransactionStore() *InMemoryTransactionStore {
	return &InMemoryTransactionStore{
		transactions: make(map[string]*gateway.TransactionStatus),
	}
}

type transactionHandler func(ctx context.Context) (*gateway.TransactionStatus, error)

type Middleware func(transactionHandler) transactionHandler

func (s *InMemoryTransactionStore) Save(ctx context.Context, saveTransaction *SaveTransaction) (*gateway.TransactionStatus, error) {
	handler := s.makeSaveHandler(saveTransaction)
	middlewares := []Middleware{
		withContextCheck,
		withNetworkSimulator,
	}

	wrappedHandler := applyMiddlewares(handler, middlewares...)

	return wrappedHandler(ctx)
}

func (s *InMemoryTransactionStore) Get(ctx context.Context, id string) (*gateway.TransactionStatus, error) {
	handler := s.makeGetHandler(id)
	middlewares := []Middleware{
		withContextCheck,
		withNetworkSimulator,
	}

	wrappedHandler := applyMiddlewares(handler, middlewares...)

	return wrappedHandler(ctx)
}

func (s *InMemoryTransactionStore) Update(ctx context.Context, updateTransaction UpdateTransaction) (*gateway.TransactionStatus, error) {
	handler := s.makeUpdateHandler(updateTransaction)
	middlewares := []Middleware{
		withContextCheck,
		withNetworkSimulator,
	}

	wrappedHandler := applyMiddlewares(handler, middlewares...)

	return wrappedHandler(ctx)
}

func applyMiddlewares(handler transactionHandler, middlewares ...Middleware) transactionHandler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}

func withContextCheck(next transactionHandler) transactionHandler {
	return func(ctx context.Context) (*gateway.TransactionStatus, error) {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return next(ctx)
	}
}

func withNetworkSimulator(next transactionHandler) transactionHandler {
	return func(ctx context.Context) (*gateway.TransactionStatus, error) {
		if shouldFail() {
			return nil, errors.New("paypal: network error")
		}

		networkLatency()

		return next(ctx)
	}
}

func (s *InMemoryTransactionStore) makeSaveHandler(saveTransaction *SaveTransaction) transactionHandler {
	saveHandler := func(ctx context.Context) (*gateway.TransactionStatus, error) {
		now := time.Now().UTC()
		transactionID := generateTransactionID()

		status := &gateway.TransactionStatus{
			TransactionID: transactionID,
			Status:        gateway.StatusPending,
			Amount:        saveTransaction.Amount,
			Currency:      saveTransaction.Currency,
			CreatedAt:     now,
			UpdatedAt:     now,
			ErrorMessage:  "",
		}

		if shouldFail() {
			status.Status = gateway.StatusFailed
			status.ErrorMessage = "paypal: failed save other"
		} else {
			status.Status = gateway.StatusCompleted
		}

		s.saveTransaction(status)

		return status, nil
	}

	return saveHandler
}

func (s *InMemoryTransactionStore) makeGetHandler(id string) transactionHandler {
	getHandler := func(ctx context.Context) (*gateway.TransactionStatus, error) {
		return s.getTransaction(id)
	}

	return getHandler
}

func (s *InMemoryTransactionStore) makeUpdateHandler(updateTransaction UpdateTransaction) transactionHandler {
	updateHandler := func(ctx context.Context) (*gateway.TransactionStatus, error) {
		t, err := s.getTransaction(updateTransaction.TransactionID)
		if err != nil {
			return nil, err
		}

		now := time.Now().UTC()
		t.Status = updateTransaction.Status
		t.UpdatedAt = now

		s.saveTransaction(t)

		return t, nil
	}

	return updateHandler
}

func (s *InMemoryTransactionStore) getTransaction(id string) (*gateway.TransactionStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transaction, ok := s.transactions[id]
	if !ok {
		return nil, errors.New("transaction not found")
	}

	return transaction, nil
}

func (s *InMemoryTransactionStore) saveTransaction(t *gateway.TransactionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transactions[t.TransactionID] = t
}

func networkLatency() {
	delay := time.Duration(500+rand.Intn(1500)) * time.Millisecond
	time.Sleep(delay)
}

func shouldFail() bool {
	return rand.Float64() < 0.5
}

func generateTransactionID() string {
	id, _ := uuid.NewRandom()
	return id.String()
}
