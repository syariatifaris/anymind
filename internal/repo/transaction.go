package repo

import (
	"context"
	"fmt"
	"github.com/syariatifaris/anymind/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

const (
	transactions = "transactions"
	transactionSnapshots = "transaction_snapshots"
)

func NewTransactionRepository(db *mongo.Database) ITransactionRepository {
	return &transactionRepository{db: db}
}

type ITransactionRepository interface {
	Add(context.Context, model.Transaction) error
	AddUpdateHourlySnapshot(context.Context, time.Time, float64) error
	AccumulateBalance(ctx context.Context, additionAmount float64)(float64, error)
	GetAccumulateBalancesInHour(context.Context, time.Time, time.Time)([]model.TransactionSnapshot, error)
}

type transactionRepository struct {
	db *mongo.Database
}

func (t *transactionRepository) Add(ctx context.Context, topUp model.Transaction) error {
	_, err := t.db.Collection(transactions).InsertOne(ctx, topUp)
	return err
}

func (t *transactionRepository) AccumulateBalance(ctx context.Context, additionAmount float64)(float64, error){
	var a model.AccumulateBalance

	err := t.db.Collection(transactionSnapshots).FindOne(ctx,
		bson.M{
			"key": "accumulate",
		},
	).Decode(&a)
	if err != nil && err != mongo.ErrNoDocuments{
		return 0, err
	}

	if err == mongo.ErrNoDocuments{
		_, err = t.db.Collection(transactionSnapshots).InsertOne(ctx, model.AccumulateBalance{
			Key:          "accumulate",
			TotalBalance: fmt.Sprint(additionAmount),
		})
		if err != nil{
			return 0, err
		}
		return additionAmount, nil
	}

	before, _ := strconv.ParseFloat(a.TotalBalance, 64)
	after := before + additionAmount

	update := bson.D{{"$set", bson.D{{"total_balance", fmt.Sprint(after)}}}}
	_, err = t.db.Collection(transactionSnapshots).UpdateOne(ctx,
		bson.M{"key": "accumulate"}, update)

	return after, nil
}

func (t *transactionRepository) AddUpdateHourlySnapshot(ctx context.Context, hourlyTimeUTC time.Time, latestAmount float64) error {
	var s model.TransactionSnapshot

	update := bson.D{{"$set", bson.D{{"accumulated_amount", fmt.Sprint(latestAmount)}}}}
	err := t.db.Collection(transactionSnapshots).FindOneAndUpdate(ctx,
		bson.M{
			"hourly_datetime_utc": hourlyTimeUTC,
		},update,
	).Decode(&s)
	if err != nil && err != mongo.ErrNoDocuments{
		return err
	}

	if err == mongo.ErrNoDocuments{
		_, err = t.db.Collection(transactionSnapshots).InsertOne(ctx, model.TransactionSnapshot{
			HourlyDateTimeUTC: hourlyTimeUTC,
			AccumulatedAmount: fmt.Sprint(latestAmount),
		})
		if err != nil{
			return err
		}
	}

	return nil
}

func (t *transactionRepository) GetAccumulateBalancesInHour(ctx context.Context, start time.Time, end time.Time)([]model.TransactionSnapshot, error){
	snapshots := make([]model.TransactionSnapshot, 0)
	curr, err := t.db.Collection(transactionSnapshots).Find(ctx, bson.M{
		"hourly_datetime_utc": bson.M{
			"$gte": start,
			"$lte": end,
		},
	})
	if err != nil{
		return nil, err
	}
	if curr != nil{
		for curr.Next(ctx){
			var snapshot model.TransactionSnapshot
			if err := curr.Decode(&snapshot); err != nil{
				return nil, fmt.Errorf("unable to decode result %s", err.Error())
			}
			snapshots = append(snapshots, snapshot)
		}
	}
	return snapshots, nil
}