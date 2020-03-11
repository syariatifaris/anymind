package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"github.com/syariatifaris/anymind/internal/common"
	"github.com/syariatifaris/anymind/internal/model"
	"github.com/syariatifaris/anymind/internal/repo"
)

func NewTransactionService(dal repo.ITransactionRepository, defLocation *time.Location) ITransactionService {
	return &transactionService{
		trxRepo:     dal,
		defLocation: defLocation,
	}
}

type ITransactionService interface {
	AddBalance(ctx context.Context, request model.AddBalanceRequest) (string, interface{}, error)
	GetBalancesInHourlyRange(ctx context.Context, start string, end string) (string, interface{}, error)
}

type transactionService struct {
	trxRepo     repo.ITransactionRepository
	defLocation *time.Location
}

func (t *transactionService) AddBalance(ctx context.Context, request model.AddBalanceRequest) (string, interface{}, error) {
	logrus.Infoln("time by request", request.DateTime)
	requestTime, err := time.Parse(time.RFC3339, request.DateTime)
	if err != nil {
		return "wrong datetime format", nil, errors.New(common.BadRequest)
	}

	diff := time.Now().UTC().Sub(requestTime.UTC())
	if diff > time.Hour{
		return "request time should not older than last hour", nil, errors.New(common.BadRequest)
	}

	logrus.Infoln("time in UTC", requestTime.UTC())
	uid := uuid.NewV4()
	trxID := uid.String()

	err = t.trxRepo.Add(ctx, model.Transaction{
		TrxID:       trxID,
		Amount:      fmt.Sprint(request.Amount),
		DatetimeUTC: requestTime.UTC(),
	})
	if err != nil {
		return "", nil, errors.New(common.InternalServerError)
	}

	after, err := t.trxRepo.AccumulateBalance(ctx, request.Amount)
	if err != nil{
		logrus.Errorln("unable to handle accumulate balance", err.Error())
		return "", nil, errors.New(common.InternalServerError)
	}

	//get time utc hour only
	utcTime := requestTime.UTC()
	utcTimeHourRounded := utcTime.Add(-1 * time.Duration(utcTime.Minute()) * time.Minute)
	hourlyUTC := utcTimeHourRounded.Format(time.RFC3339)
	logrus.Infoln("hourly format in UTC: ", hourlyUTC)

	err = t.trxRepo.AddUpdateHourlySnapshot(ctx, utcTimeHourRounded, after)
	if err != nil{
		logrus.Errorln("unable to handle add update", err.Error())
		return "", nil, errors.New(common.InternalServerError)
	}

	return "balance was added successfully", trxID, nil
}

func(t *transactionService) GetBalancesInHourlyRange(ctx context.Context, start string, end string) (string, interface{}, error){
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return "wrong start datetime request", nil, errors.New(common.BadRequest)
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil{
		return "wrong end datetime request", nil, errors.New(common.BadRequest)
	}

	startUTC := startTime.UTC()
	endUTC := endTime.UTC()

	startUTC = startUTC.Add(-1 * time.Duration(startUTC.Minute()) * time.Minute)
	endUTC = endUTC.Add(-1 * time.Duration(endUTC.Minute()) * time.Minute)

	logrus.Infoln("start time in utc = ", startUTC.String(), " / end time in utc = ", endUTC.String())

	snapshots, err := t.trxRepo.GetAccumulateBalancesInHour(ctx, startUTC, endUTC)
	if err != nil{
		logrus.Errorln("unable to get accumulated balance in hour range", err.Error())
		return "", nil, errors.New(common.InternalServerError)
	}

	var balances []model.BalanceHourlyRange
	for _, s := range snapshots{
		f, _ := strconv.ParseFloat(s.AccumulatedAmount, 1)
		balances = append(balances, model.BalanceHourlyRange{
			Datetime: s.HourlyDateTimeUTC.Format(time.RFC3339),
			Amount:   fmt.Sprintf("%.1f", f),
		})
	}

	return fmt.Sprintf("%d record(s) retrieved", len(balances)), balances, nil
}
