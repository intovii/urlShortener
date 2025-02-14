package usecase

import (
	"URLShortener/domain"
	"context"
	"regexp"

	"go.uber.org/zap"
)


const (
	base = 63
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
)

func (uc *Usecase) OnStart(_ context.Context) error {
	uc.log.Info("start usecase")
	return nil
}

func (uc *Usecase) OnStop(_ context.Context) error {
	uc.log.Info("stop usecase")
	return nil
}

func (uc *Usecase) Create(ctx context.Context, url string) (string, error) {
	if isURLValid := validateUrl(url); !isURLValid {
		uc.log.Error("entry data is not url")
		return "", domain.ErrInvalidArgument
	}

	count, err := uc.Repo.LenRows(ctx)
	if err != nil {
		uc.log.Error("can not get rows count from db:", zap.Error(err))
		return "", err
	}

	shortUrl := hash(count)

	if err := uc.Repo.Create(ctx, shortUrl, url); err != nil {
		uc.log.Error("failed to create pair : shoerUrl - url")
		return "", err
	}

	return shortUrl, nil
}

func (uc *Usecase) Get(ctx context.Context, shortUrl string) (string, error) {
	originalUrl, err := uc.Repo.Get(ctx, shortUrl)
	if err != nil {
		uc.log.Error("failed to get url")
		return "", err
	}

	return originalUrl, nil
}

func validateUrl(url string) bool{
	re := regexp.MustCompile(".+\\..+")
	return re.MatchString(url)
}


func hash(num int) string {
	result := make([]byte, 10)

	for i := 9; i >= 0; i-- {
		result[i] = alphabet[num%base]
		num /= base
	}

	return string(result)
}
