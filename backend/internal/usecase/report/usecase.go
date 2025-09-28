package report

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/report"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/s3/image"
	report2 "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
)

var (
	errorNoAccess = errors.New("user has no access")
)

type Usecase interface {
	GetByID(ctx context.Context, id uuid.UUID) (report2.Report, bool, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]report2.Report, error)
	Get(ctx context.Context, limit, offset int64) ([]report2.Report, error)
	GetByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (report2.Report, bool, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, report report2.Report, images []*multipart.FileHeader) error
}

type usecase struct {
	db report.Repo
	s3 image.Repo
}

func New(db report.Repo) Usecase {
	return &usecase{db: db}
}

func (u *usecase) Get(ctx context.Context, limit, offset int64) ([]report2.Report, error) {
	return u.db.Get(ctx, limit, offset)
}

func (u *usecase) Count(ctx context.Context) (int64, error) {
	return u.db.Count(ctx)
}

func (u *usecase) GetByID(ctx context.Context, id uuid.UUID) (report2.Report, bool, error) {
	return u.db.GetByID(ctx, id)
}

func (u *usecase) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int64) ([]report2.Report, error) {
	return u.db.GetByUserID(ctx, userID, limit, offset)
}

func (u *usecase) GetByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (report2.Report, bool, error) {
	rep, ok, err := u.db.GetByID(ctx, id)
	if err != nil {
		return report2.Report{}, false, err
	}

	if !ok {
		return report2.Report{}, false, nil
	}

	if rep.UserID != userID {
		return report2.Report{}, false, errorNoAccess
	}

	return rep, true, nil
}

func (u *usecase) Update(ctx context.Context, report report2.Report, images []*multipart.FileHeader) error {
	if err := u.removeOldImages(ctx, report); err != nil {
		return err
	}

	for _, img := range images {
		url, err := u.saveImage(ctx, img)
		if err != nil {
			return err
		}

		report.Images = append(report.Images, report2.Image{
			ID:   uuid.New(),
			Link: string(url),
		})
	}

	return u.db.Upsert(ctx, report)
}

func (u *usecase) removeOldImages(ctx context.Context, report report2.Report) error {
	oldImages, err := u.db.GetImagesByReportID(ctx, report.ID)
	if err != nil {
		return err
	}

	for _, img := range oldImages {
		if err := u.s3.Delete(ctx, report2.ImageURL(img.Link)); err != nil {
			return err
		}
	}
	return nil
}

func (u *usecase) saveImage(ctx context.Context, image *multipart.FileHeader) (report2.ImageURL, error) {
	ext := filepath.Ext(image.Filename)
	contentType := image.Header.Get("Content-Type")
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	url, err := u.s3.Save(ctx, ext, contentType, file)
	if err != nil {
		return "", err
	}
	return url, nil
}
