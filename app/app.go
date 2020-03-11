package app

import "github.com/syariatifaris/anymind/internal/service"

type Application interface {
	Start() error
}

type UseCase struct {
	TransactionService service.ITransactionService
}