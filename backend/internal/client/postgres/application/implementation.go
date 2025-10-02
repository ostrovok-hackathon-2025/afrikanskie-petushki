package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
)

const pgForeginKeyErr = "foreign_key_violation"

type applicationRepo struct {
	db *sqlx.DB
}

type Getter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func NewApplicationRepo(db *sqlx.DB) ApplicationRepo {
	return &applicationRepo{
		db: db,
	}
}

func (r *applicationRepo) CreateApplication(
	ctx context.Context,
	application *application.Application,
) error {
	tx, err := r.db.Beginx()

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT 
		o.participants_limit,
		(SELECT COUNT(*) FROM application a WHERE a.offer_id = o.id) as participants_count
	FROM offer as o WHERE o.id = $1`

	var Limits LimitsDTO

	err = tx.GetContext(ctx, &Limits, query, application.OfferId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOfferNotExist
		}

		return fmt.Errorf("failed to get limits from db: %w", err)
	}

	if Limits.ParticipantsCount >= Limits.ParticipantsLimit {
		err = ErrParticipantsLimit
		return err
	}

	userAppLimitInfo, err := getUserAppLimitInfo(tx, ctx, application.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user app limit info: %w", err)
	}
	if userAppLimitInfo == nil || userAppLimitInfo.Limit-userAppLimitInfo.ActiveAppCount <= 0 {
		err = ErrAppLimit
		return ErrAppLimit
	}

	query = `
	INSERT INTO application (id, user_id, offer_id, status) VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, query, application.Id, application.UserId, application.OfferId, application.Status)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == pgForeginKeyErr {
			return ErrOfferNotExist
		}

		return fmt.Errorf("failed to insert application into db: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit create application")
	}

	return nil
}

func (r *applicationRepo) GetApplications(
	ctx context.Context,
	userId uuid.UUID,
	pageNum, pageSize int,
) ([]*application.Application, int, error) {
	offset := pageNum * pageSize

	query := `
	SELECT a.id, a.user_id, a.offer_id, a.status, o.expiration_at, h.name FROM application as a
	INNER JOIN offer as o ON a.offer_id = o.id
	INNER JOIN hotel as h ON o.hotel_id = h.id
	WHERE a.user_id = $1
	OFFSET $2
	LIMIT $3
	`

	var apps []ApplicationDTO

	err := r.db.SelectContext(ctx, &apps, query, userId, offset, pageSize)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if pageNum == 0 {
				return []*application.Application{}, 0, nil
			}

			return nil, 0, ErrPageNotFound
		}

		return nil, 0, fmt.Errorf("failed to get all applications from db: %w", err)
	}

	var count int
	query = "SELECT COUNT(*) FROM application WHERE user_id = $1"

	err = r.db.Get(&count, query, userId)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count of applications: %w", err)
	}

	pagesCount := count / pageSize
	if count%pageSize != 0 {
		pagesCount++
	}

	res := make([]*application.Application, 0, len(apps))

	for _, app := range apps {
		res = append(res, app.ToApplicationModel())
	}

	return res, pagesCount, nil
}

func (r *applicationRepo) GetApplicationById(
	ctx context.Context,
	applicationId uuid.UUID,
) (*application.Application, error) {
	query := `
	SELECT a.id, a.user_id, a.offer_id, a.status, o.expiration_at, h.name FROM application as a
	INNER JOIN offer as o ON a.offer_id = o.id
	INNER JOIN hotel as h ON o.hotel_id = h.id
	WHERE a.id = $1
	`

	var app ApplicationDTO

	err := r.db.GetContext(ctx, &app, query, applicationId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApplicationNotFound
		}

		return nil, fmt.Errorf("failed to get application from repo: %w", err)
	}

	res := app.ToApplicationModel()

	return res, nil
}

func (r *applicationRepo) GetByOfferID(
	ctx context.Context,
	offerID uuid.UUID,
) ([]*application.Application, error) {
	query := `
	SELECT id, user_id, offer_id, status
	FROM application
	WHERE offer_id = $1
	`

	var apps []ApplicationDTO

	err := r.db.SelectContext(ctx, &apps, query, offerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*application.Application{}, nil // пустой слайс, если нет заявок. когда отель кал
		}
		return nil, fmt.Errorf("failed to get applications by offer_id: %w", err)
	}

	res := make([]*application.Application, 0, len(apps))
	for _, app := range apps {
		res = append(res, app.ToApplicationModel())
	}

	return res, nil
}

func (r *applicationRepo) GetByOfferIDForDraw(
	ctx context.Context,
	offerID uuid.UUID,
) ([]*application.ApplicationWithRating, error) {
	query := `
	SELECT a.id as id, a.user_id, a.offer_id, a.status, u.rating
	FROM application a JOIN public."user" u on u.id = a.user_id
	WHERE offer_id = $1
	`

	var apps []ApplicationWithRatingDTO

	err := r.db.SelectContext(ctx, &apps, query, offerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*application.ApplicationWithRating{}, nil // пустой слайс, если нет заявок. когда отель кал
		}
		return nil, fmt.Errorf("failed to get applications by offer_id: %w", err)
	}

	res := make([]*application.ApplicationWithRating, 0, len(apps))
	for _, app := range apps {
		res = append(res, app.ToModel())
	}

	return res, nil
}

func (r *applicationRepo) GetUserAppLimitInfo(ctx context.Context, userID uuid.UUID) (*application.UserAppLimitInfo, error) {

	dto, err := getUserAppLimitInfo(r.db, ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToModel(), nil
}

func getUserAppLimitInfo(s Getter, ctx context.Context, userID uuid.UUID) (*UserAppLimitInfoDTO, error) {
	query := `
	SELECT u.app_limit, (
    	SELECT COUNT(*)
    	FROM application a
    	WHERE a.user_id = u.id AND a.status = $2
	) AS active_app_count
	FROM "user" u
	WHERE u.id = $1;
	`
	var userLimitInfo UserAppLimitInfoDTO

	err := s.GetContext(ctx, &userLimitInfo, query, userID, application.APPLICATION_CREATED)
	if err != nil {
		return nil, fmt.Errorf("failed to get count of user applications: %w", err)
	}
	return &userLimitInfo, nil
}

func (r *applicationRepo) UpdateApplicationStatus(ctx context.Context, application *application.Application) error {
	query := `UPDATE application SET status = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, application.Status, application.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *applicationRepo) GetByFilter(
	ctx context.Context,
	filter *application.Filter,
) ([]*application.Application, error) {
	sql := sq.Select(
		"a.id",
		"a.user_id",
		"a.offer_id",
		"a.status",
		"o.expiration_at",
		"h.name",
	).From("application as a").
		Join("offer as o ON a.offer_id = o.id").
		Join("hotel as h ON o.hotel_id = h.id")
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	if roomID, ok := filter.RoomID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.room_id": roomID})
	}
	if hotelID, ok := filter.HotelID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.hotel_id": hotelID})
	}
	if status, ok := filter.Status.Get(); ok {
		sql = sql.Where(sq.Eq{"a.status": status})
	}
	query, args, err := sql.Limit(filter.Limit).Offset(filter.Offset).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	var dtos []*ApplicationDTO
	err = r.db.SelectContext(ctx, &dtos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications by filter: %w", err)
	}

	applications := make([]*application.Application, len(dtos))
	for i, dto := range dtos {
		applications[i] = dto.ToApplicationModel()
	}
	return applications, nil
}

func (r *applicationRepo) GetCountByFilter(
	ctx context.Context,
	filter *application.Filter,
) (int, error) {
	sql := sq.Select("COUNT(*)").From("application as a").
		Join("offer as o ON a.offer_id = o.id").
		Join("hotel as h ON o.hotel_id = h.id")
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	if roomID, ok := filter.RoomID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.room_id": roomID})
	}
	if hotelID, ok := filter.HotelID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.hotel_id": hotelID})
	}
	if status, ok := filter.Status.Get(); ok {
		sql = sql.Where(sq.Eq{"a.status": status})
	}
	query, args, err := sql.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build sql query: %w", err)
	}
	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get count of applications: %w", err)
	}
	return count, nil
}
