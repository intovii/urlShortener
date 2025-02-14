package postgres

import (
	"URLShortener/config"
	"URLShortener/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	ctx context.Context
	cfg *config.ConfigModel
	log *zap.Logger
	DB  *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, cfg *config.ConfigModel, log *zap.Logger) (*PostgresRepository, error) {
	return &PostgresRepository{
		ctx: ctx,
		cfg: cfg,
		log: log,
		}, nil
}

func (r *PostgresRepository) OnStart(_ context.Context) error {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		r.cfg.Postgres.Host,
		r.cfg.Postgres.Port,
		r.cfg.Postgres.User,
		r.cfg.Postgres.Password,
		r.cfg.Postgres.DBName,
		r.cfg.Postgres.SSLMode)
	pool, err := pgxpool.Connect(r.ctx, connectionUrl)
	if err != nil {
		return err
	}
	r.DB = pool
	r.log.Info("start postgres")
	return nil
}

func (r *PostgresRepository) OnStop(_ context.Context) error {
	r.DB.Close()
	r.log.Info("stop postgres")
	return nil
}

const queryGetOriginalUrl = `
SELECT 
	url
FROM 
	urls
WHERE
	short_url = $1
`

func (r *PostgresRepository) Get(ctx context.Context, shortUrl string) (string, error) {
	var originalUrl string
	err := r.DB.QueryRow(ctx, queryGetOriginalUrl, shortUrl).Scan(&originalUrl)
	if err != nil {
		r.log.Error("can not get original url")
		return "", domain.ErrNotFound
	}
	return originalUrl, nil
}


const queryCreateRow = `
INSERT INTO urls 
	(short_url, url) 
VALUES
	($1, $2)
`

func (r *PostgresRepository) Create(ctx context.Context, shortUrl, url string) error {
	_, err := r.DB.Exec(ctx, queryCreateRow, shortUrl, url)
	if err != nil {
		r.log.Error("can not insert new row")
		return err
	}
	return nil
}

const queryGetCountElems = `
SELECT COUNT(*) FROM urls
`

func (r *PostgresRepository) LenRows(ctx context.Context) (int, error) {
	var count int
	err := r.DB.QueryRow(ctx, queryGetCountElems).Scan(&count)
	if err != nil {
		r.log.Error("can not get row count")
		return 0, err
	}
	return count, nil
}