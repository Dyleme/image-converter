package service_test

import (
	"context"
	"testing"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	originalImageURL  = "./test_data/x.png"
	convertedImageURL = "./test_data/xconv.png"
)

func TestConvertRequest_Convert(t *testing.T) {
	testCases := []struct {
		testName  string
		reqID     int
		filename  string
		configure func(*mocks.MockStorager, *mocks.MockConvertRepo)
		wantErr   error
	}{
		{
			testName: "all is good",
			reqID:    12,
			filename: "file",
			configure: func(ms *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				info := &model.ConvImageInfo{
					OldURL:  originalImageURL,
					OldType: "png",
					NewType: "jpeg",
					OldImID: 32,
					UserID:  1,
					Ratio:   0.5,
				}
				newURL := "newURL"
				originalBytes := loadImage(t, originalImageURL)
				convertedBytes := loadImage(t, convertedImageURL)
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(info, nil)
				ms.EXPECT().GetFile(ctx, originalImageURL).Times(1).
					Return(originalBytes, nil)
				mcr.EXPECT().SetImageResolution(ctx, info.OldImID, 1152, 648).
					Return(nil)
				ms.EXPECT().UploadFile(ctx, info.UserID, "file", convertedBytes).
					Return(newURL, nil)
				newImgInfo := &model.ReuquestImageInfo{
					URL:  newURL,
					Type: info.NewType,
				}
				mcr.EXPECT().AddProcessedImage(ctx, info.UserID, 12, newImgInfo,
					576, 324, repository.StatusDone, gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			testName: "error in get conv info",
			reqID:    12,
			filename: "file",
			configure: func(_ *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(nil, errRepository)
			},
			wantErr: errRepository,
		},
		{
			testName: "error in get file",
			reqID:    12,
			filename: "file",
			configure: func(ms *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				info := &model.ConvImageInfo{
					OldURL:  originalImageURL,
					OldType: "png",
					NewType: "jpeg",
					OldImID: 32,
					UserID:  1,
					Ratio:   0.5,
				}
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(info, nil)
				ms.EXPECT().GetFile(ctx, originalImageURL).Times(1).
					Return(nil, errStorage)
			},
			wantErr: errStorage,
		},
		{
			testName: "error in set image resolution",
			reqID:    12,
			filename: "file",
			configure: func(ms *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				info := &model.ConvImageInfo{
					OldURL:  originalImageURL,
					OldType: "png",
					NewType: "jpeg",
					OldImID: 32,
					UserID:  1,
					Ratio:   0.5,
				}
				bts := loadImage(t, originalImageURL)
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(info, nil)
				ms.EXPECT().GetFile(ctx, originalImageURL).Times(1).
					Return(bts, nil)
				mcr.EXPECT().SetImageResolution(ctx, info.OldImID, 1152, 648).
					Return(errRepository)
			},
			wantErr: errRepository,
		},
		{
			testName: "error in file uploading",
			reqID:    12,
			filename: "file",
			configure: func(ms *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				info := &model.ConvImageInfo{
					OldURL:  originalImageURL,
					OldType: "png",
					NewType: "jpeg",
					OldImID: 32,
					UserID:  1,
					Ratio:   0.5,
				}
				bts := loadImage(t, originalImageURL)
				convbts := loadImage(t, convertedImageURL)
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(info, nil)
				ms.EXPECT().GetFile(ctx, originalImageURL).Times(1).
					Return(bts, nil)
				mcr.EXPECT().SetImageResolution(ctx, info.OldImID, 1152, 648).
					Return(nil)
				ms.EXPECT().UploadFile(ctx, info.UserID, "file", convbts).
					Return("", errStorage)
			},
			wantErr: errStorage,
		},
		{
			testName: "error in add processed image",
			reqID:    12,
			filename: "file",
			configure: func(ms *mocks.MockStorager, mcr *mocks.MockConvertRepo) {
				ctx := context.Background()
				info := &model.ConvImageInfo{
					OldURL:  originalImageURL,
					OldType: "png",
					NewType: "jpeg",
					OldImID: 32,
					UserID:  1,
					Ratio:   0.5,
				}
				bts := loadImage(t, originalImageURL)
				convbts := loadImage(t, convertedImageURL)
				newURL := "newURL"
				mcr.EXPECT().GetConvInfo(ctx, 12).Times(1).
					Return(info, nil)
				ms.EXPECT().GetFile(ctx, originalImageURL).Times(1).
					Return(bts, nil)
				mcr.EXPECT().SetImageResolution(ctx, info.OldImID, 1152, 648).
					Return(nil)
				ms.EXPECT().UploadFile(ctx, info.UserID, "file", convbts).
					Return(newURL, nil)
				newImgInfo := &model.ReuquestImageInfo{
					URL:  newURL,
					Type: info.NewType,
				}
				mcr.EXPECT().AddProcessedImage(ctx, info.UserID, 12, newImgInfo,
					576, 324, repository.StatusDone, gomock.Any()).
					Return(errRepository)
			},
			wantErr: errRepository,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			storMock := mocks.NewMockStorager(mockCtr)
			convRepoMock := mocks.NewMockConvertRepo(mockCtr)

			tc.configure(storMock, convRepoMock)

			conv := service.NewConvertRequest(convRepoMock, storMock)

			gotErr := conv.Convert(context.Background(), tc.reqID, tc.filename)
			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
