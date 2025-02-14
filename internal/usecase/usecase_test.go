package usecase

import (
	"URLShortener/domain"
	mock_usecase "URLShortener/internal/usecase/mocks"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGet(t *testing.T) {
	type mockBehavior func(r *mock_usecase.MockRepository, shortUrl string)
	tests := []struct {
		name         string
		in           string
		mockBehavior mockBehavior
		want         string
		wantErr      error
	}{
		{
			name: "ok",
			in:   "aaaaaaaaaa",
			mockBehavior: func(r *mock_usecase.MockRepository, shortUrl string) {
				r.EXPECT().Get(context.Background(), shortUrl).Return("google.com", nil)
			},
			want:    "google.com",
			wantErr: nil,
		},
		{
			name: "Not found",
			in:   "1234567890",
			mockBehavior: func(r *mock_usecase.MockRepository, shortUrl string) {
				r.EXPECT().Get(context.Background(), shortUrl).Return("", domain.ErrNotFound)
			},
			want:    "",
			wantErr: domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_usecase.NewMockRepository(ctrl)
			tt.mockBehavior(mockRepo, tt.in)

			log, _ := zap.NewProduction()
			usecase,_ := NewUsecase(log, mockRepo)

			response, err := usecase.Get(context.Background(), tt.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, response)
		})
	}
}

func TestCreate(t *testing.T) {
	type mockBehavior func(r *mock_usecase.MockRepository, shortUrl, originalUrl string)
	type args struct {
		ShortUrl    string
		OriginalUrl string
	}
	tests := []struct {
		name string
		args
		mockBehavior
		want    string
		wantErr error
	}{
		{
			name: "ok",
			args: args{OriginalUrl: "google.com", ShortUrl: "aaaaaaaaaa"},
			mockBehavior: func(r *mock_usecase.MockRepository, shortUrl, originalUrl string) {
				r.EXPECT().Create(context.Background(), shortUrl, originalUrl).Return(nil)
				r.EXPECT().LenRows(context.Background()).Return(0, nil)
			},
			want:    "aaaaaaaaaa",
			wantErr: nil,
		},
		{
			name:         "invalid format",
			args:         args{OriginalUrl: "exmaplecom", ShortUrl: ""},
			mockBehavior: func(r *mock_usecase.MockRepository, shortUrl, originalUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
		{
			name:         "invalid format",
			args:         args{OriginalUrl: "", ShortUrl: ""},
			mockBehavior: func(r *mock_usecase.MockRepository, shortUrl, originalUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_usecase.NewMockRepository(ctrl)
			tt.mockBehavior(s, tt.args.ShortUrl, tt.args.OriginalUrl)

			log, _ := zap.NewProduction()
			usecase,_ := NewUsecase(log, s)

			response, err := usecase.Create(context.Background(), tt.args.OriginalUrl)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, response)
		})
	}
}