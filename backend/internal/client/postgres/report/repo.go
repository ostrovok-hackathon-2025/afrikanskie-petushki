package report

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
)

type Repo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Report, bool, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]model.Report, error)
	Get(ctx context.Context, limit, offset int64) ([]model.Report, error)
	Count(ctx context.Context) (int64, error)
	Upsert(ctx context.Context, report model.Report) error
	GetImagesByReportID(ctx context.Context, reportID uuid.UUID) ([]model.Image, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{db: db}
}

const queryGetByID = `
        SELECT 
            r.id,
            a.user_id,
            r.application_id,
            r.expiration_at,
            r.status,
            r.text,
            p.id as "image_id",
            p.s3_link as "image_link",
            a.user_id as "user_id"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
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
		ImageID       *uuid.UUID `db:"image_id"`
		ImageLink     *string    `db:"image_link"`
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
            p.id as "image_id",
            p.s3_link as "image_link"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
        WHERE a.user_id = $1
        ORDER BY r.expiration_at DESC
        LIMIT $2 OFFSET $3
    `

func (r *repo) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]model.Report, error) {
	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		ApplicationID uuid.UUID  `db:"application_id"`
		ExpirationAt  time.Time  `db:"expiration_at"`
		Status        string     `db:"status"`
		Text          string     `db:"text"`
		ImageID       *uuid.UUID `db:"image_id"`
		ImageLink     *string    `db:"image_link"`
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
            p.id as "image_id",
            p.s3_link as "image_link"
        FROM report r
        LEFT JOIN photo p ON r.id = p.report_id
        LEFT JOIN application a ON a.id = r.application_id
        WHERE r.id IN (
            SELECT id FROM report 
            ORDER BY expiration_at DESC 
            LIMIT $1 OFFSET $2
        )
        ORDER BY r.expiration_at DESC, p.id
    `

func (r *repo) Get(ctx context.Context, limit, offset int64) ([]model.Report, error) {
	var rows []struct {
		ID            uuid.UUID  `db:"id"`
		UserID        uuid.UUID  `db:"user_id"`
		ApplicationID uuid.UUID  `db:"application_id"`
		ExpirationAt  time.Time  `db:"expiration_at"`
		Status        string     `db:"status"`
		Text          string     `db:"text"`
		ImageID       *uuid.UUID `db:"image_id"`
		ImageLink     *string    `db:"image_link"`
	}

	err := sqlx.SelectContext(ctx, r.db, &rows, queryGet, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []model.Report{}, nil
	}

	// Группируем строки по отчетам
	reportsMap := make(map[uuid.UUID]*model.Report)
	reportOrder := make([]uuid.UUID, 0, len(rows)) // Для сохранения порядка

	for _, row := range rows {
		if _, exists := reportsMap[row.ID]; !exists {
			reportsMap[row.ID] = &model.Report{
				ID:            row.ID,
				UserID:        row.UserID,
				ApplicationID: row.ApplicationID,
				ExpirationAt:  row.ExpirationAt,
				Status:        row.Status,
				Text:          row.Text,
				Images:        make([]model.Image, 0),
			}
			reportOrder = append(reportOrder, row.ID)
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

	// Восстанавливаем порядок из запроса
	reports := make([]model.Report, 0, len(reportOrder))
	for _, reportID := range reportOrder {
		if report, exists := reportsMap[reportID]; exists {
			reports = append(reports, *report)
		}
	}

	return reports, nil
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

const reportUpsertQuery = `
        INSERT INTO report (id, application_id, expiration_at, status, text)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE SET
            application_id = EXCLUDED.application_id,
            expiration_at = EXCLUDED.expiration_at,
            status = EXCLUDED.status,
            text = EXCLUDED.text
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
		report.ID,
		report.ApplicationID,
		report.ExpirationAt,
		report.Status,
		report.Text,
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
