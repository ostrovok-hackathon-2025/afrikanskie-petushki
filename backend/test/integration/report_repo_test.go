//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	offerRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	reportRepo "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/report"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/offer"
	model "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg/testhelper"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite
	db  *sqlx.DB
	ctx context.Context
}

func TestRepo(t *testing.T) {
	suite.Run(t, &RepoSuite{})
}

func (suite *RepoSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = testhelper.NewPostgreSqlx(suite.T())
}

func (suite *RepoSuite) SetupTest() {
	for _, query := range []string{
		"truncate photo, report;",
	} {
		_, err := suite.db.ExecContext(suite.ctx, query)
		suite.Require().NoError(err, query)
	}
}

func (suite *RepoSuite) TearDownTest() {
	for _, query := range []string{
		"truncate photo, report;",
	} {
		_, err := suite.db.ExecContext(suite.ctx, query)
		suite.Require().NoError(err, query)
	}
}

func (suite *RepoSuite) TestCreate() {
	// Arrange
	ctx, cancel := context.WithTimeout(suite.ctx, time.Minute*3)
	defer cancel()

	reportID := uuid.New()
	applicationID := uuid.MustParse("6fa459ea-ee8a-3ca4-894e-db77e160355e")
	offerID := uuid.New()
	hotelID := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")    // Moscow Grand Hotel
	locationID := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479") // Москва
	roomID := uuid.MustParse("359d2a75-0237-4ce7-9e8f-adbc61357aa2")
	userID := uuid.MustParse("6fa459ea-ee8a-3ca4-894e-db77e160355e")

	report := model.Report{
		ID:            reportID,
		UserID:        userID,
		ApplicationID: applicationID,
		ExpirationAt:  time.Now().Truncate(0),
		Status:        "accepted",
		Text:          "where is some text",
		Images: []model.Image{
			{
				ID:   uuid.New(),
				Link: "vk.com/some_image1",
			},
			{
				ID:   uuid.New(),
				Link: "vk.com/some_image2",
			},
		},
	}

	repo := reportRepo.NewRepo(suite.db)
	offerRepository := offerRepo.New(suite.db, log.Default())

	// Act

	createOfferErr := offerRepository.Create(ctx, offerID, offer.Create{
		Task:              "Сделать классный селфи",
		RoomID:            roomID,
		CheckIn:           time.Now().Add(24 * time.Hour),
		CheckOut:          time.Now().Add(48 * time.Hour),
		ExpirationAT:      time.Now().Add(time.Hour),
		HotelID:           hotelID,
		LocalID:           locationID,
		ParticipantsLimit: 10,
	})

	upsertErr := repo.Upsert(ctx, report)

	getByIDRes, getByIDOk, getByIDErr := repo.GetByID(ctx, report.ID)
	getImagesByReportIDRes, getImagesByReportIDErr := repo.GetImagesByReportID(ctx, report.ID)
	getByUserIDRes, getByUserIDErr := repo.GetByUserID(ctx, report.UserID, 10, 0)
	getRes, getErr := repo.Get(ctx, 10, 0)
	countRes, countErr := repo.Count(ctx)

	// Assert
	suite.Require().NoError(createOfferErr)
	suite.Require().NoError(upsertErr)

	suite.Require().NoError(getByIDErr)
	suite.Require().True(getByIDOk)
	suite.Require().EqualValues(report, getByIDRes)

	suite.Require().NoError(getImagesByReportIDErr)
	suite.Require().EqualValues(report.Images, getImagesByReportIDRes)

	suite.Require().NoError(getByUserIDErr)
	suite.Require().Len(getByUserIDRes, 1)
	suite.Require().EqualValues(report, getByUserIDRes[0])

	suite.Require().NoError(getErr)
	suite.Require().Len(getRes, 1, fmt.Sprintf("К-во записей по Get должно быть 1, получено %d", len(getRes)))
	suite.Require().ElementsMatch(report.Images, getRes[0].Images)
	testReport, testRes := report, getRes[0]
	testReport.Images = nil
	testRes.Images = nil
	suite.Require().EqualValues(testReport, testRes)

	suite.Require().NoError(countErr)
	suite.Require().Equal(int64(1), countRes)
}
