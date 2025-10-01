package analytics

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/analytics"
)

type Repo interface {
	GetAnalytics(ctx context.Context) (analytics.Analytics, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

const queryGetAnalytics = `
SELECT 
	COUNT(DISTINCT CASE WHEN o.status = 'done' 
    AND o.check_in_at >= DATE_TRUNC('month', CURRENT_DATE)
    AND o.check_in_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
    THEN o.id END) as completed_offers,
    
    COUNT(DISTINCT CASE WHEN a.created_at >= DATE_TRUNC('month', CURRENT_DATE)
	AND a.created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
	THEN a.id END) as applications_received,
    
    COUNT(DISTINCT CASE WHEN r.status = 'accepted'
	AND r.created_at >= DATE_TRUNC('month', CURRENT_DATE)
	AND r.created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
	THEN r.id END) as accepted_reports
	
FROM offer o
LEFT JOIN application a ON o.id = a.offer_id
LEFT JOIN report r ON a.id = r.application_id	
`

func (r *repo) GetAnalytics(ctx context.Context) (analytics.Analytics, error) {
	var dto struct {
		CompletedOffers      uint64 `db:"completed_offers"`
		ApplicationsReceived uint64 `db:"applications_received"`
		AcceptedReports      uint64 `db:"accepted_reports"`
	}

	if err := r.db.GetContext(ctx, &dto, queryGetAnalytics); err != nil {
		return analytics.Analytics{}, err
	}

	return analytics.Analytics{
		CompletedOffers:      dto.CompletedOffers,
		ApplicationsReceived: dto.ApplicationsReceived,
		AcceptedReports:      dto.AcceptedReports,
	}, nil
}
