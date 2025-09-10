package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"subscription/internal/config"
	"subscription/internal/model"
	"subscription/internal/pkg/migration"
	"subscription/internal/pkg/servererrors"
	"subscription/internal/repository/dto"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	addServiceQuery    = `INSERT INTO services (name) VALUES ($1) RETURNING *`
	getServiceQuery    = `SELECT service_id,name FROM services WHERE service_id=$1`
	getServicesQuery   = `SELECT service_id,name FROM services`
	updateServiceQuery = `UPDATE services SET name=$2 WHERE service_id=$1`
	removeServiceQuery = `DELETE FROM services WHERE service_id=$1`

	addSubscriptionQuery = `
INSERT INTO subscriptions (service_id,price,user_id,start_date,stop_date) 
VALUES ((SELECT service_id FROM services WHERE "name"=$1),$2,$3,$4,$5) 
RETURNING subscription_id,service_id,price,user_id,start_date,stop_date`
	getSubscriptionQuery = `
SELECT subscription_id,service_id,price,user_id,start_date,stop_date 
FROM subscriptions WHERE subscription_id=$1`
	getSubscriptionsQuery = `
SELECT subscription_id,service_id,price,user_id,start_date,stop_date 
FROM subscriptions OFFSET COALESCE($1,0) LIMIT COALESCE($2,10)`
	getSubscriptionTotalQuery = `
WITH t AS (
	SELECT
		AGE(LEAST(COALESCE(stop_date, $2),$2),GREATEST(start_date,$1)) AS inter,
		price
	FROM subscriptions
	WHERE 
		(start_date<=$2 AND (stop_date IS null OR stop_date>=$1))AND
		($3::uuid IS null or user_id=$3) AND
		($4::character varying IS null or service_id=(SELECT service_id FROM services WHERE "name"=$4)) 
	)
SELECT 	COALESCE (SUM( (EXTRACT(YEAR FROM inter)*12 + EXTRACT (MONTH FROM inter)+1)*price),0) FROM t`
	updateSubscriptionQuery = `
UPDATE subscriptions 
SET service_id=(SELECT service_id FROM services WHERE "name"=$2),price=$3,user_id=$4,start_date=$5,stop_date=$6 
WHERE subscription_id=$1`
	removeSubscriptionQuery = `DELETE FROM subscriptions WHERE subscription_id=$1`
)

type Repository interface {
	AddService(name string) (*model.Service, error)
	GetService(serviceId int) (*model.Service, error)
	GetServices() ([]*model.Service, error)
	UpdateService(dto *dto.UpdateService) error
	RemoveService(serviceId int) error
	AddSubscription(dto *dto.AddSubscription) (*model.Subscription, error)
	GetSubscription(subscriptionId int) (*model.Subscription, error)
	GetSubscriptions(dto *dto.GetSubscriptions) ([]*model.Subscription, error)
	GetSubscriptionTotal(ctx *dto.GetSubscriptionTotal) (int, error)
	UpdateSubscription(dto *dto.UpdateSubscription) error
	RemoveSubscription(subscriptionId int) error
}

type repository struct {
	conn *pgx.Conn
	lg   *zap.SugaredLogger
}

func MustNew(lg *zap.SugaredLogger, cfg *config.Db) *repository {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
	var (
		conn  *pgx.Conn
		err   error
		count int
	)
	for {
		count++
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		conn, err = pgx.Connect(ctx, connString)
		if err != nil {
			lg.Errorf("failed to connect repository(%d): %v", count, err)
			if count > 4 {
				lg.Fatalf("failed to connect repository", err)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	lg.Info("repository connect successfully")

	// migration
	config := conn.Config()
	db := stdlib.OpenDB(*config)
	if err := migration.Up(db, cfg.MigrationPath); err != nil {
		lg.Errorf("migration failed: %v", err)
	} else {
		lg.Info("migration completed successfully")
	}

	return &repository{
		conn: conn,
		lg:   lg,
	}
}
func (r *repository) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.conn.Close(ctx); err != nil {
		r.lg.Errorf("failed to disconnect repository: %v", err)
		return
	}
	r.lg.Info("repository disconnect successfully")
}

func (r *repository) AddService(name string) (*model.Service, error) {
	service := new(model.Service)
	err := r.conn.QueryRow(context.Background(), addServiceQuery, name).Scan(&service.ServiceId, &service.Name)
	if err != nil {
		r.lg.Errorf("failed to add service: %v", err)
		return nil, servererrors.ErrorInternal
	}
	return service, nil
}
func (r *repository) GetService(serviceId int) (*model.Service, error) {
	service := new(model.Service)
	err := r.conn.QueryRow(context.Background(), getServiceQuery, serviceId).Scan(&service.ServiceId, &service.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, servererrors.ErrorRecordNotFound
	}
	if err != nil {
		r.lg.Errorf("failed to get service: %v", err)
		return nil, servererrors.ErrorInternal
	}
	return service, err
}
func (r *repository) GetServices() ([]*model.Service, error) {
	rows, err := r.conn.Query(context.Background(), getServicesQuery)
	if err != nil {
		r.lg.Errorf("failed to get services: %v", err)
		return nil, err
	}
	defer rows.Close()

	services := []*model.Service{}
	for rows.Next() {
		service := new(model.Service)
		err := rows.Scan(&service.ServiceId, &service.Name)
		if err != nil {
			r.lg.Errorf("failed to get services: %v", err)
			return nil, servererrors.ErrorInternal
		}
		services = append(services, service)
	}
	return services, nil
}
func (r *repository) UpdateService(dto *dto.UpdateService) error {
	result, err := r.conn.Exec(context.Background(), updateServiceQuery, dto.ServiceId, dto.Name)
	if err != nil {
		r.lg.Errorf("failed to update service: %v", err)
		return servererrors.ErrorInternal
	}
	if result.RowsAffected() == 0 {
		return servererrors.ErrorRecordNotFound
	}
	return nil
}
func (r *repository) RemoveService(serviceId int) error {
	result, err := r.conn.Exec(context.Background(), removeServiceQuery, serviceId)
	if err != nil {
		r.lg.Errorf("failed to remove service: %v", err)
		return servererrors.ErrorInternal
	}
	if result.RowsAffected() == 0 {
		return servererrors.ErrorRecordNotFound
	}
	return nil
}

func (r *repository) AddSubscription(dto *dto.AddSubscription) (*model.Subscription, error) {

	subscription := new(model.Subscription)
	err := r.conn.QueryRow(
		context.Background(),
		addSubscriptionQuery,
		dto.ServiceName,
		dto.Price,
		dto.UserId,
		dto.StartDate,
		dto.StopDate,
	).Scan(
		&subscription.SubscriptionId,
		&subscription.ServiceId,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.StopDate,
	)
	if err != nil {
		r.lg.Errorf("failed to add subscription: %v", err)
		return nil, servererrors.ErrorInternal
	}

	return subscription, nil
}
func (r *repository) GetSubscription(subscriptionId int) (*model.Subscription, error) {

	subscription := new(model.Subscription)
	err := r.conn.QueryRow(
		context.Background(),
		getSubscriptionQuery,
		subscriptionId,
	).Scan(
		&subscription.SubscriptionId,
		&subscription.ServiceId,
		&subscription.Price,
		&subscription.UserId,
		&subscription.StartDate,
		&subscription.StopDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, servererrors.ErrorRecordNotFound
	}
	if err != nil {
		r.lg.Errorf("failed to get subscription: %v", err)
		return nil, servererrors.ErrorInternal
	}
	return subscription, nil
}

func (r *repository) GetSubscriptions(dto *dto.GetSubscriptions) ([]*model.Subscription, error) {
	rows, err := r.conn.Query(context.Background(), getSubscriptionsQuery, dto.Offset, dto.Limit)
	if err != nil {
		r.lg.Errorf("failed to get subscriptions: %v", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := []*model.Subscription{}
	for rows.Next() {
		subscription := new(model.Subscription)
		err := rows.Scan(
			&subscription.SubscriptionId,
			&subscription.ServiceId,
			&subscription.Price,
			&subscription.UserId,
			&subscription.StartDate,
			&subscription.StopDate,
		)
		if err != nil {
			r.lg.Errorf("failed to get subscriptions: %v", err)
			return nil, servererrors.ErrorInternal
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
func (r *repository) GetSubscriptionTotal(dto *dto.GetSubscriptionTotal) (int, error) {
	var total int
	err := r.conn.QueryRow(
		context.Background(),
		getSubscriptionTotalQuery,
		dto.StartDate,
		dto.StopDate,
		dto.UserId,
		dto.ServiceName,
	).Scan(
		&total,
	)
	if err != nil {
		r.lg.Errorf("failed to get subscription total: %v", err)
		return 0, servererrors.ErrorInternal
	}
	return total, nil
}
func (r *repository) UpdateSubscription(dto *dto.UpdateSubscription) error {
	result, err := r.conn.Exec(
		context.Background(),
		updateSubscriptionQuery,
		dto.SubscriptionId,
		dto.ServiceName,
		dto.Price,
		dto.UserId,
		dto.StartDate,
		dto.StopDate,
	)
	if err != nil {
		r.lg.Errorf("failed to update subscription: %v", err)
		return servererrors.ErrorInternal
	}
	if result.RowsAffected() == 0 {
		return servererrors.ErrorRecordNotFound
	}
	return nil
}
func (r *repository) RemoveSubscription(subscriptionId int) error {
	result, err := r.conn.Exec(context.Background(), removeSubscriptionQuery, subscriptionId)
	if err != nil {
		r.lg.Errorf("failed to remove subscription: %v", err)
		return servererrors.ErrorInternal
	}
	if result.RowsAffected() == 0 {
		return servererrors.ErrorRecordNotFound
	}
	return nil
}
