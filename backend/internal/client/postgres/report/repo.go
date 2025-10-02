package report

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
)

type Repo interface {
	Create(ctx context.Context, createData model.Report) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Report, bool, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]model.Report, error)
	Get(ctx context.Context, limit, offset int64) ([]model.Report, error)
	Count(ctx context.Context) (int64, error)
	CountByUserId(ctx context.Context, userId uuid.UUID) (int64, error)
	Upsert(ctx context.Context, report model.Report) error
	GetImagesByReportID(ctx context.Context, reportID uuid.UUID) ([]model.Image, error)
	UpdateStatus(ctx context.Context, report model.Report) error
	UpdatePromocode(ctx context.Context, report model.Report) error
	GetByApplicationId(ctx context.Context, applicationId uuid.UUID) (uuid.UUID, uuid.UUID, error)
	GetByFilter(ctx context.Context, filter model.Filter) ([]model.Report, error)
	GetCountByFilter(ctx context.Context, filter model.Filter) (int, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

const queryCreate = `
	INSERT INTO report (id, application_id, expiration_at, status, text)
	VALUES ($1, $2, $3, $4, $5)
`

type getRow struct {
	ID            uuid.UUID  `db:"id"`
	UserID        uuid.UUID  `db:"user_id"`
	ApplicationID uuid.UUID  `db:"application_id"`
	ExpirationAt  time.Time  `db:"expiration_at"`
	Status        string     `db:"status"`
	Text          string     `db:"text"`
	Promocode     string     `db:"promocode"`
	ImageID       *uuid.UUID `db:"image_id"`
	ImageLink     *string    `db:"image_link"`
	HotelName     string     `db:"hotel_name"`
	LocationName  string     `db:"location_name"`
	RoomName      string     `db:"room_name"`
	Task          string     `db:"task"`
	CheckInAt     time.Time  `db:"check_in_at"`
	CheckOutAt    time.Time  `db:"check_out_at"`
}

func (r *repo) Create(ctx context.Context, createData model.Report) error {
	_, err := r.db.ExecContext(
		ctx,
		queryCreate,
		createData.ID,
		createData.ApplicationID,
		createData.ExpirationAt,
		createData.Status,
		createData.Text,
	)

	if err != nil {
		return fmt.Errorf("failed to create report: %w", err)
	}

	return nil
}

const queryGetByID = `
        SELECT 
            r.id,
            a.user_id,
            r.application_id,
            r.expiration_at,
            r.status,
            r.text,
			r.promocode,
            p.id as "image_id",
            p.s3_link as "image_link",
            a.user_id as "user_id",
			h.name as "hotel_name",
			l.name as "location_name",
			o.task as "task",
			o.check_in_at,
			o.check_out_at,
			m.name as "room_name"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
		INNER JOIN offer o ON o.id = a.offer_id
		INNER JOIN hotel h ON h.id = o.hotel_id
		INNER JOIN location l ON l.id = h.location_id   
		INNER JOIN room m ON m.id = o.room_id
        WHERE r.id = $1
    `

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (model.Report, bool, error) {
	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		ApplicationID uuid.UUID  `db:"application_id"`
		UserID        uuid.UUID  `db:"user_id"`
		ExpirationAt  time.Time  `db:"expiration_at"`
		Status        string     `db:"status"`
		Text          string     `db:"text"`
		Promocode     string     `db:"promocode"`
		ImageID       *uuid.UUID `db:"image_id"`
		ImageLink     *string    `db:"image_link"`
		HotelName     string     `db:"hotel_name"`
		LocationName  string     `db:"location_name"`
		RoomName      string     `db:"room_name"`
		Task          string     `db:"task"`
		CheckInAt     time.Time  `db:"check_in_at"`
		CheckOutAt    time.Time  `db:"check_out_at"`
	}

	err := sqlx.SelectContext(ctx, r.db, &rows, queryGetByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Report{}, false, nil
		}
		return model.Report{}, false, err
	}

	if len(rows) == 0 {
		return model.Report{}, false, nil
	}

	// Собираем отчет из первой строки (основные данные)
	report := model.Report{
		ID:            rows[0].ID,
		UserID:        rows[0].UserID,
		ApplicationID: rows[0].ApplicationID,
		ExpirationAt:  rows[0].ExpirationAt,
		Status:        rows[0].Status,
		Text:          rows[0].Text,
		Promocode:     rows[0].Promocode,
		LocationName:  rows[0].LocationName,
		HotelName:     rows[0].HotelName,
		RoomName:      rows[0].RoomName,
		Task:          rows[0].Task,
		Images:        make([]model.Image, 0),
	}

	// Собираем изображения из всех строк
	for _, row := range rows {
		if row.ImageID != nil && row.ImageLink != nil {
			report.Images = append(report.Images, model.Image{
				ID:   *row.ImageID,
				Link: *row.ImageLink,
			})
		}
	}

	return report, true, nil
}

const queryGetByUserID = `
        SELECT 
            r.id,
            r.application_id,
            r.expiration_at,
            r.status,
            r.text,
			r.promocode,
            p.id as "image_id",
            p.s3_link as "image_link",
			h.name as "hotel_name",
			l.name as "location_name",
			o.task as "task",
			o.check_in_at,
			o.check_out_at,
			m.name as "room_name"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
		INNER JOIN offer o ON o.id = a.offer_id
		INNER JOIN hotel h ON h.id = o.hotel_id
		INNER JOIN location l ON l.id = h.location_id   
		INNER JOIN room m ON m.id = o.room_id
        WHERE a.user_id = $1
        ORDER BY r.expiration_at DESC
        LIMIT $2 OFFSET $3
    `

func (r *repo) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]model.Report, error) {
	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		ApplicationID uuid.UUID  `db:"application_id"`
		ExpirationAt  time.Time  `db:"expiration_at"`
		UserID        uuid.UUID  `db:"user_id"`
		Status        string     `db:"status"`
		Text          string     `db:"text"`
		Promocode     string     `db:"promocode"`
		ImageID       *uuid.UUID `db:"image_id"`
		ImageLink     *string    `db:"image_link"`
		HotelName     string     `db:"hotel_name"`
		LocationName  string     `db:"location_name"`
		RoomName      string     `db:"room_name"`
		Task          string     `db:"task"`
		CheckInAt     time.Time  `db:"check_in_at"`
		CheckOutAt    time.Time  `db:"check_out_at"`
	}

	err := sqlx.SelectContext(ctx, r.db, &rows, queryGetByUserID, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []model.Report{}, nil
	}

	// Группируем строки по отчетам
	reportsMap := make(map[uuid.UUID]*model.Report)

	for _, row := range rows {
		if _, exists := reportsMap[row.ID]; !exists {
			reportsMap[row.ID] = &model.Report{
				ID:            row.ID,
				ApplicationID: row.ApplicationID,
				UserID:        userID,
				ExpirationAt:  row.ExpirationAt,
				Status:        row.Status,
				Text:          row.Text,
				Promocode:     row.Promocode,
				LocationName:  row.LocationName,
				HotelName:     row.HotelName,
				RoomName:      row.RoomName,
				Task:          row.Task,
				Images:        make([]model.Image, 0),
			}
		}

		// Добавляем фото, если оно есть
		if row.ImageID != nil && row.ImageLink != nil {
			report := reportsMap[row.ID]
			report.Images = append(report.Images, model.Image{
				ID:   *row.ImageID,
				Link: *row.ImageLink,
			})
		}
	}

	// Конвертируем map в slice
	reports := make([]model.Report, 0, len(reportsMap))
	for _, report := range reportsMap {
		reports = append(reports, *report)
	}

	return reports, nil
}

const queryGet = `
        SELECT 
            r.id,
            a.user_id,
            r.application_id,
            r.expiration_at,
            r.status,
            r.text,
			r.promocode,
            p.id as "image_id",
            p.s3_link as "image_link",
			h.name as "hotel_name",
			l.name as "location_name",
			o.task as "task",
			o.check_in_at,
			o.check_out_at,
			m.name as "room_name"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
		INNER JOIN offer o ON o.id = a.offer_id
		INNER JOIN hotel h ON h.id = o.hotel_id
		INNER JOIN location l ON l.id = h.location_id   
		INNER JOIN room m ON m.id = o.room_id
        WHERE r.id IN (
            SELECT id FROM report 
            ORDER BY expiration_at DESC 
            LIMIT $1 OFFSET $2
        )
        ORDER BY r.expiration_at DESC, p.id
    `

func (r *repo) Get(ctx context.Context, limit, offset int64) ([]model.Report, error) {
	var rows []getRow

	err := sqlx.SelectContext(ctx, r.db, &rows, queryGet, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []model.Report{}, nil
	}

	return convertGetDtoToModel(rows), nil
}

const queryCount = `SELECT COUNT(*) FROM report`

func (r *repo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count, queryCount)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const queryCountByUserId = `
SELECT COUNT(*) 
FROM report r 
INNER JOIN application a ON a.id = r.application_id
WHERE a.user_id = $1`

func (r *repo) CountByUserId(ctx context.Context, userId uuid.UUID) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, r.db, &count, queryCountByUserId, userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const reportUpsertQuery = `
        UPDATE report SET text = $1, status = $2 WHERE id = $3
    `
const deletePhotosQuery = `DELETE FROM photo WHERE report_id = $1`
const insertPhotoQuery = `
            INSERT INTO photo (id, report_id, s3_link)
            VALUES (:id, :report_id, :s3_link)
        `

func (r *repo) Upsert(ctx context.Context, report model.Report) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Upsert для отчета
	_, err = tx.ExecContext(ctx, reportUpsertQuery,
		report.Text,
		report.Status,
		report.ID,
	)
	if err != nil {
		return err
	}

	// Удаляем старые фотографии
	_, err = tx.ExecContext(ctx, deletePhotosQuery, report.ID)
	if err != nil {
		return err
	}

	// Пакетная вставка новых фотографий
	if len(report.Images) > 0 {
		photos := make([]map[string]interface{}, len(report.Images))
		for i, image := range report.Images {
			photos[i] = map[string]interface{}{
				"id":        image.ID,
				"report_id": report.ID,
				"s3_link":   image.Link,
			}
		}

		_, err = tx.NamedExecContext(ctx, insertPhotoQuery, photos)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

const queryGetImagesByReportID = `
        SELECT id, report_id, s3_link 
        FROM photo 
        WHERE report_id = $1
    `

func (r *repo) GetImagesByReportID(ctx context.Context, reportID uuid.UUID) ([]model.Image, error) {
	// Объявляем структуру для работы с БД прямо в методе
	type dbImage struct {
		ID       uuid.UUID `db:"id"`
		ReportID uuid.UUID `db:"report_id"`
		S3Link   string    `db:"s3_link"`
	}

	var dbImages []dbImage
	err := r.db.SelectContext(ctx, &dbImages, queryGetImagesByReportID, reportID)
	if err != nil {
		return nil, err
	}

	// Конвертируем из внутренней структуры БД в доменную модель
	images := make([]model.Image, 0, len(dbImages))
	for _, dbImg := range dbImages {
		images = append(images, model.Image{
			ID:   dbImg.ID,
			Link: dbImg.S3Link,
		})
	}

	return images, nil
}

const reportUpdateStatusQuery = `
        UPDATE report SET status = $1 WHERE id = $2
    `

func (r *repo) UpdateStatus(ctx context.Context, report model.Report) error {
	// TODO: check if status is "created" or "filled"
	_, err := r.db.ExecContext(ctx, reportUpdateStatusQuery,
		report.Status,
		report.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

const reportUpdatePromocodeQuery = `
        UPDATE report SET promocode = $1 WHERE id = $2
    `

func (r *repo) UpdatePromocode(ctx context.Context, report model.Report) error {
	_, err := r.db.ExecContext(ctx, reportUpdatePromocodeQuery,
		report.Promocode,
		report.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

const reportGetByApplicationIdQuery = `
        SELECT r.id, a.user_id 
		FROM report r
		JOIN application a ON a.id = r.application_id 
		WHERE a.id = $1
    `

func (r *repo) GetByApplicationId(ctx context.Context, applicationId uuid.UUID) (uuid.UUID, uuid.UUID, error) {
	var res struct {
		Id     uuid.UUID `db:"id"`
		UserId uuid.UUID `db:"user_id"`
	}

	err := r.db.GetContext(ctx, &res, reportGetByApplicationIdQuery, applicationId)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, err
	}

	return res.Id, res.UserId, nil
}

func (r *repo) GetByFilter(ctx context.Context, filter model.Filter) ([]model.Report, error) {
	sql := sq.Select(
		"r.id",
		"a.user_id",
		"r.application_id",
		"r.expiration_at",
		"r.status",
		"r.text",
		"p.id as image_id",
		"p.s3_link as image_link",
		"h.name as hotel_name",
		"l.name as location_name",
		"o.task as task",
		"o.check_in_at",
		"o.check_out_at",
		"m.name as room_name").
		From("report r").
		LeftJoin("photo p ON r.id = p.report_id").
		LeftJoin("application a ON a.id = r.application_id").
		Join("offer o ON o.id = a.offer_id").
		Join("hotel h ON h.id = o.hotel_id").
		Join("location l ON l.id = h.location_id").
		Join("room m ON m.id = o.room_id").
		OrderBy("r.expiration_at DESC, p.id")
	if status, ok := filter.Status.Get(); ok {
		sql = sql.Where(sq.Eq{"r.status": status})
	}
	if hotelID, ok := filter.HotelID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.hotel_id": hotelID})
	}
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err := sql.Limit(filter.Limit).Offset(filter.Offset).PlaceholderFormat(sq.Dollar).ToSql()
	log.Println(query)
	if err != nil {
		return nil, err
	}
	var rows []getRow
	err = r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []model.Report{}, nil
	}
	return convertGetDtoToModel(rows), nil
}

func (r *repo) GetCountByFilter(ctx context.Context, filter model.Filter) (int, error) {
	sql := sq.Select("count(*)").
		From("report r").
		LeftJoin("photo p ON r.id = p.report_id").
		LeftJoin("application a ON a.id = r.application_id").
		Join("offer o ON o.id = a.offer_id").
		Join("hotel h ON h.id = o.hotel_id").
		Join("location l ON l.id = h.location_id").
		Join("room m ON m.id = o.room_id")
	if status, ok := filter.Status.Get(); ok {
		sql = sql.Where(sq.Eq{"r.status": status})
	}
	if hotelID, ok := filter.HotelID.Get(); ok {
		sql = sql.Where(sq.Eq{"o.hotel_id": hotelID})
	}
	if locationID, ok := filter.LocationID.Get(); ok {
		sql = sql.Where(sq.Eq{"h.location_id": locationID})
	}
	query, args, err := sql.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, err
	}
	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func convertGetDtoToModel(rows []getRow) []model.Report {
	// Группируем строки по отчетам
	reportsMap := make(map[uuid.UUID]*model.Report)

	for _, row := range rows {
		if _, exists := reportsMap[row.ID]; !exists {
			reportsMap[row.ID] = &model.Report{
				ID:            row.ID,
				ApplicationID: row.ApplicationID,
				UserID:        row.UserID,
				ExpirationAt:  row.ExpirationAt,
				Status:        row.Status,
				Text:          row.Text,
				Promocode:     row.Promocode,
				LocationName:  row.LocationName,
				HotelName:     row.HotelName,
				RoomName:      row.RoomName,
				Task:          row.Task,
				Images:        make([]model.Image, 0),
			}
		}

		// Добавляем фото, если оно есть
		if row.ImageID != nil && row.ImageLink != nil {
			report := reportsMap[row.ID]
			report.Images = append(report.Images, model.Image{
				ID:   *row.ImageID,
				Link: *row.ImageLink,
			})
		}
	}

	// Конвертируем map в slice
	reports := make([]model.Report, 0, len(reportsMap))
	for _, report := range reportsMap {
		reports = append(reports, *report)
	}
	return reports
}
