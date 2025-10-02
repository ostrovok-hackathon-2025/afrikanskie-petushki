package worker

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/application"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/offer"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/report"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/client/postgres/user"
	appModel "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/application"
	reportModel "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/report"
	userModel "github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/user"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/pkg"
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

		applications, err := w.applicationRepo.GetByOfferID(ctx, i.ID)
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

func (w *SecretGuestWorker) selectWinnerByRating(ctx context.Context, applications []*appModel.Application) (*appModel.Application, error) {
	if len(applications) == 0 {
		return nil, errors.New("no applications found")
	}

	users := make([]userModel.User, 0, len(applications))
	appByUserID := make(map[uuid.UUID]*appModel.Application)

	for _, app := range applications {
		user, err := w.userRepo.GetUserById(ctx, app.UserId)
		if err != nil {
			log.Printf("⚠️ Warning: failed to get user %s: %v", app.UserId, err)
			continue
		}

		users = append(users, *user)
		appByUserID[user.ID] = app
	}

	if len(users) == 0 {
		return nil, errors.New("no valid users found for applications")
	}

	winnerUserID := pkg.ChooseByRating(users)

	if winnerUserID == uuid.Nil {
		return nil, errors.New("failed to select winner")
	}

	winnerApp, ok := appByUserID[winnerUserID]
	if !ok {
		return nil, errors.New("winner application not found")
	}

	return winnerApp, nil
}
