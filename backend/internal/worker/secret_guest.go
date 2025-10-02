package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/report"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	appModel "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
	reportModel "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
)

type SecretGuestWorker struct {
	offerRepo       offer.Repo
	applicationRepo application.ApplicationRepo
	reportRepo      report.Repo
	userRepo        user.Repo
	scheduler       *gocron.Scheduler
}

func NewSecretGuestWorker(
	offerRepo offer.Repo,
	applicationRepo application.ApplicationRepo,
	reportRepo report.Repo,
	userRepo user.Repo,
) *SecretGuestWorker {
	return &SecretGuestWorker{
		offerRepo:       offerRepo,
		applicationRepo: applicationRepo,
		reportRepo:      reportRepo,
		userRepo:        userRepo,
		scheduler:       gocron.NewScheduler(time.UTC),
	}
}

func (w *SecretGuestWorker) Start() {
	log.Println("Worker started")

	_, err := w.scheduler.Every(1).Minutes().Do(func() {
		w.process()
	})
	if err != nil {
		return
	}

	w.process()

	w.scheduler.StartAsync()

	log.Println("Worker on 5 minutes")
}

func (w *SecretGuestWorker) Stop() {
	log.Println("Worker stopping")
	w.scheduler.Stop()
	log.Println("Worker stopped")
}

func (w *SecretGuestWorker) process() {
	ctx := context.Background()

	log.Printf("\n [%s] Start offer processing\n", time.Now().Format("15:04:05"))

	offers, err := w.offerRepo.GetByExpirationTime(ctx)
	if err != nil {
		log.Printf("❌ Error offers receive %v", err)
		return
	}

	if len(offers) == 0 {
		log.Println("No expired offers")
		return
	}

	for _, i := range offers {
		if err := w.offerRepo.EditStatus(ctx, i.ID, "in_progress"); err != nil {
			log.Printf("❌ Error change status %v", err)
			continue
		}

		applications, err := w.applicationRepo.GetByOfferIDForDraw(ctx, i.ID)
		if err != nil {
			log.Printf("❌Error application receive %v", err)
			continue
		}

		if len(applications) == 0 {
			err := w.offerRepo.EditStatus(ctx, i.ID, "done")
			if err != nil {
				return
			}
			continue
		}

		winner, err := w.selectWinnerByRating(ctx, applications)
		if err != nil {
			log.Printf("❌ Error selecting winner: %v", err)
			continue
		}

		log.Printf("🎉 winner: %s (User %s)", winner.Id, winner.UserId)

		if err := w.offerRepo.EditStatus(ctx, i.ID, "done"); err != nil {
			log.Printf("❌ Error end offer: %v", err)
			continue
		}

		report := reportModel.NewReport(winner.Id, winner.ExpirationAt)

		if err := w.reportRepo.Create(ctx, report); err != nil {
			log.Printf("❌ Error create report: %v", err)
			continue
		}

		winner.Status = appModel.APPLICATION_ACCEPTED

		if err := w.applicationRepo.UpdateApplicationStatus(ctx, winner); err != nil {
			log.Printf("❌ Error update application status: %v", err)
			continue
		}
	}

	log.Println("\n Processing done")
}

func (w *SecretGuestWorker) selectWinnerByRating(_ context.Context, applications []*appModel.ApplicationWithRating) (*appModel.Application, error) {
	if len(applications) == 0 {
		return nil, errors.New("no applications found")
	}

	winnerApp := chooseByRating(applications)

	if winnerApp == nil {
		return nil, errors.New("failed to select winner")
	}

	return &appModel.Application{
		Id:      winnerApp.Id,
		UserId:  winnerApp.UserId,
		OfferId: winnerApp.OfferId,
		Status:  winnerApp.Status,
	}, nil
}
