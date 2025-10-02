package achievement

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ostrovok-hackathon-2025/afrikanskie-petushki/backend/internal/model/achievement"
)

type Repo interface {
	AddAchievement(ctx context.Context, userId, achievementId uuid.UUID) error
	GetAchievements(ctx context.Context) ([]achievement.Achievement, error)
	GetAchievementsByUserId(ctx context.Context, userId uuid.UUID) ([]achievement.Achievement, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) AddAchievement(ctx context.Context, userId, achievementId uuid.UUID) error {
	query := `
		INSERT INTO user_achievement (id, user_id, achievement_id) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, achievement_id) 
		DO NOTHING
	`

	id := uuid.New()

	_, err := r.db.ExecContext(ctx, query, id, userId, achievementId)

	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetAchievements(ctx context.Context) ([]achievement.Achievement, error) {
	query := `
		SELECT id, name, rating_limit FROM achievement
	`

	var dto []struct {
		Id           uuid.UUID `db:"id"`
		Name         string    `db:"name"`
		RaitingLimit int       `db:"rating_limit"`
	}

	if err := r.db.SelectContext(ctx, &dto, query); err != nil {
		return nil, err
	}

	res := make([]achievement.Achievement, 0, len(dto))

	for _, e := range dto {
		res = append(res, achievement.Achievement{
			Id:           e.Id,
			Name:         e.Name,
			RaitingLimit: e.RaitingLimit,
		})
	}

	return res, nil
}

func (r *repo) GetAchievementsByUserId(ctx context.Context, userId uuid.UUID) ([]achievement.Achievement, error) {
	query := `
		SELECT a.id, a.name, a.rating_limit 
		FROM achievement a
		INNER JOIN user_achievement u ON a.id = u.achievement_id
		WHERE u.user_id = $1
	`

	var dto []struct {
		Id           uuid.UUID `db:"id"`
		Name         string    `db:"name"`
		RaitingLimit int       `db:"rating_limit"`
	}

	if err := r.db.SelectContext(ctx, &dto, query, userId); err != nil {
		return nil, err
	}

	res := make([]achievement.Achievement, 0, len(dto))

	for _, e := range dto {
		res = append(res, achievement.Achievement{
			Id:           e.Id,
			Name:         e.Name,
			RaitingLimit: e.RaitingLimit,
		})
	}

	return res, nil
}
