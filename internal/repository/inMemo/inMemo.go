package inMemory

import (
	"URLShortener/domain"
	"context"
	"sync"
	"time"
	"fmt"
	"go.uber.org/zap"
)

type item struct {
	value string
	lastAccess int64
}

type InMemoRepository struct {
	log *zap.Logger
	DB  map[string]*item
	mu  *sync.RWMutex
}

func NewInMemoRepository(ctx context.Context, log *zap.Logger) (*InMemoRepository, error) {
	return &InMemoRepository{
		log: log,
		DB: make(map[string]*item),
		mu: &sync.RWMutex{},
	}, nil
}

func (im *InMemoRepository) OnStart(_ context.Context) error {
	im.log.Info("start in-memory")
	return nil
}

func (im *InMemoRepository) OnStop(_ context.Context) error {
	im.log.Info("stop in-memory")
	return nil
}

func (im *InMemoRepository) Create(ctx context.Context, key, value string) (err error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, ok := im.DB[key]; !ok {
		im.DB[key] = &item{value: value, lastAccess: time.Now().Unix()}
	}
	return nil
}

func (im *InMemoRepository) Get(ctx context.Context, key string) (string, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	it, ok := im.DB[key]
	if !ok {
		im.log.Error(fmt.Sprintf("key does not exist: %s\n", key))
		return "", domain.ErrNotFound
	}
	return it.value, nil
}

func (im *InMemoRepository) Delete(ctx context.Context, key string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	_, ok := im.DB[key]
	if !ok {
		im.log.Error("key does not exist")
		return domain.ErrNotFound
	}
	delete(im.DB, key)
	return nil
}

func (im *InMemoRepository) LenRows(ctx context.Context) (int, error) {
	return len(im.DB), nil
}