package report

import (
	"context"
	"errors"
	"fmt"

	"mime/multipart"
	"path/filepath"

	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/handler/rest/validation"

	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/ostrovok"
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
	CountByUserId(ctx context.Context, userId uuid.UUID) (int64, error)
	Update(ctx context.Context, report report2.Report, images []*multipart.FileHeader) error
	UpdateStatus(ctx context.Context, report report2.Report) error
	GetByApplicationId(ctx context.Context, applicationId, userId uuid.UUID) (uuid.UUID, error)
	GetByFilter(ctx context.Context, filter report2.Filter) ([]report2.Report, int, error)
}

type usecase struct {
	db             report.Repo
	s3             image.Repo
	ostrovokClient ostrovok.Client
	userRepo       user.Repo
	appsRepo       application.ApplicationRepo
}

func New(
	db report.Repo,
	s3 image.Repo,
	ostrovokClient ostrovok.Client,
	userRepo user.Repo,
	appsRepo application.ApplicationRepo,
) Usecase {
	return &usecase{
		db:             db,
		s3:             s3,
		ostrovokClient: ostrovokClient,
		userRepo:       userRepo,
		appsRepo:       appsRepo,
	}
}

func (u *usecase) Get(ctx context.Context, limit, offset int64) ([]report2.Report, error) {
	return u.db.Get(ctx, limit, offset)
}

func (u *usecase) Count(ctx context.Context) (int64, error) {
	return u.db.Count(ctx)
}

func (u *usecase) CountByUserId(ctx context.Context, userId uuid.UUID) (int64, error) {
	return u.db.CountByUserId(ctx, userId)
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

func (u *usecase) UpdateStatus(ctx context.Context, report report2.Report) error {
	if report.Status != "accepted" && report.Status != "declined" {
		return errors.New("invalid status")
	}

	if report.Status == "accepted" {
		promocode, err := u.ostrovokClient.GeneratePromocode(ctx)

		if err != nil {
			return err
		}

		report.Promocode = promocode

		err = u.db.UpdatePromocode(ctx, report)

		if err != nil {
			return fmt.Errorf("failed to save promocode: %w", err)
		}
	}

	err := u.db.UpdateStatus(ctx, report)
	if err != nil {
		return err
	}

	user, err := u.userRepo.GetUserByReportId(ctx, report.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	var newRating int
	if report.Status == "accepted" {
		newRating = user.Rating + 20
	} else if report.Status == "declined" {
		newRating = user.Rating - 5
	}

	newRating = validation.ValidateRating(newRating)

	return u.userRepo.UpdateRating(ctx, user.ID, newRating)
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

func (u *usecase) GetByApplicationId(ctx context.Context, applicationId, userId uuid.UUID) (uuid.UUID, error) {
	id, actualUserId, err := u.db.GetByApplicationId(ctx, applicationId)

	if err != nil {
		return uuid.UUID{}, err
	}

	if actualUserId != userId {
		return uuid.UUID{}, errors.New("not owner of application")
	}

	return id, nil
}

func (u *usecase) GetByFilter(ctx context.Context, filter report2.Filter) ([]report2.Report, int, error) {
	reports, err := u.db.GetByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.db.GetCountByFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	// Не считая остаток
	withOutRemainder := count / int(filter.Limit)
	if count%int(filter.Limit) == 0 {
		return reports, withOutRemainder, nil
	}
	return reports, withOutRemainder + 1, nil
}
