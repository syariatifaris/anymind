package rest

import (
	"errors"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/syariatifaris/anymind/app"
	"github.com/syariatifaris/anymind/internal/common"
	"github.com/syariatifaris/anymind/internal/model"
)

func NewRestApplication(useCase *app.UseCase, port int) app.Application {
	return &restApplication{useCase: useCase, port: port}
}

type restApplication struct {
	useCase *app.UseCase
	port    int
}

func (ra *restApplication) Start() error {
	mx := mux.NewRouter()

	v1 := mx.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/balance", ra.handleUseCase(ra.addBalance)).Methods(http.MethodPost)
	v1.HandleFunc("/balance/hours/{start}/{end}", ra.handleUseCase(ra.getBalanceInHourlyRange)).Methods(http.MethodGet)

	return http.ListenAndServe(fmt.Sprintf(":%d", ra.port), mx)
}

func (ra *restApplication) addBalance(r *http.Request) (*HTTPData, error) {
	var request model.AddBalanceRequest

	err := readPostData(r, &request)
	if err != nil {
		return nil, errors.New(common.BadRequest)
	}

	msg, data, err := ra.useCase.TransactionService.AddBalance(r.Context(), request)
	return &HTTPData{Data: data, Message: msg}, err
}

func (ra *restApplication) getBalanceInHourlyRange(r *http.Request)(*HTTPData, error){
	params := mux.Vars(r)
	start := params["start"]
	end := params["end"]

	if start == "" || end == ""{
		return &HTTPData{Message: "start / end time should not empty"},
			errors.New(common.BadRequest)
	}

	msg, data, err := ra.useCase.TransactionService.GetBalancesInHourlyRange(r.Context(), start, end)
	return &HTTPData{Data: data, Message: msg}, err
}