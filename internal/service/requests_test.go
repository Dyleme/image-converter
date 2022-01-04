package service_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/service/mocks"
	"github.com/Dyleme/image-coverter/internal/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errRepository = errors.New("error in repository")
	errStorage    = errors.New("error in storage")
)

func TestGetRequests(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		repReqs  []model.Request
		repErr   error
		wantReqs []model.Request
		wantErr  error
	}{
		{
			testName: "All is good",
			userID:   123,
			repReqs: []model.Request{
				{
					ID:             12,
					OpStatus:       "queued",
					RequestTime:    time.Date(2000, 12, 3, 2, 32, 12, 12, time.Local),
					CompletionTime: time.Time{},
					OriginalID:     5,
					ProcessedID:    0,
					Ratio:          0.23,
					OriginalType:   "jpeg",
					ProcessedType:  "png",
				},
				{
					ID:             23,
					OpStatus:       "done",
					RequestTime:    time.Date(2012, 12, 3, 2, 32, 12, 12, time.Local),
					CompletionTime: time.Date(2012, 12, 3, 2, 35, 12, 12, time.Local),
					OriginalID:     5,
					ProcessedID:    6,
					Ratio:          0.23,
					OriginalType:   "jpeg",
					ProcessedType:  "png",
				},
			},
			repErr: nil,
			wantReqs: []model.Request{
				{
					ID:             12,
					OpStatus:       "queued",
					RequestTime:    time.Date(2000, 12, 3, 2, 32, 12, 12, time.Local),
					CompletionTime: time.Time{},
					OriginalID:     5,
					ProcessedID:    0,
					Ratio:          0.23,
					OriginalType:   "jpeg",
					ProcessedType:  "png",
				},
				{
					ID:             23,
					OpStatus:       "done",
					RequestTime:    time.Date(2012, 12, 3, 2, 32, 12, 12, time.Local),
					CompletionTime: time.Date(2012, 12, 3, 2, 35, 12, 12, time.Local),
					OriginalID:     5,
					ProcessedID:    6,
					Ratio:          0.23,
					OriginalType:   "jpeg",
					ProcessedType:  "png",
				},
			},
			wantErr: nil,
		},
		{
			testName: "Repository error",
			userID:   123,
			repReqs:  nil,
			repErr:   errRepository,
			wantReqs: nil,
			wantErr:  errRepository,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequestRepo(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			srvc := service.NewRequest(mockRequest, mockStorage, mocks.NewMockImageProcesser(mockCtr))
			ctx := context.Background()

			mockRequest.EXPECT().GetRequests(ctx, tc.userID).Return(tc.repReqs, tc.repErr)

			gotReqs, gotErr := srvc.GetRequests(ctx, tc.userID)
			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotReqs, tc.wantReqs)
		})
	}
}

func TestGetRequest(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		reqID    int
		repReq   *model.Request
		repErr   error
		wantReq  *model.Request
		wantErr  error
	}{
		{
			testName: "All is good",
			userID:   123,
			reqID:    4,
			repReq: &model.Request{
				ID:             12,
				OpStatus:       "queued",
				RequestTime:    time.Date(2000, 12, 3, 2, 32, 12, 12, time.Local),
				CompletionTime: time.Time{},
				OriginalID:     5,
				ProcessedID:    0,
				Ratio:          0.23,
				OriginalType:   "jpeg",
				ProcessedType:  "png",
			},
			repErr: nil,
			wantReq: &model.Request{
				ID:             12,
				OpStatus:       "queued",
				RequestTime:    time.Date(2000, 12, 3, 2, 32, 12, 12, time.Local),
				CompletionTime: time.Time{},
				OriginalID:     5,
				ProcessedID:    0,
				Ratio:          0.23,
				OriginalType:   "jpeg",
				ProcessedType:  "png",
			},
			wantErr: nil,
		},
		{
			testName: "Repository error",
			userID:   123,
			reqID:    5,
			repReq:   &model.Request{},
			repErr:   errRepository,
			wantReq:  nil,
			wantErr:  errRepository,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequestRepo(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			srvc := service.NewRequest(mockRequest, mockStorage, &mocks.MockImageProcesser{})
			ctx := context.Background()

			mockRequest.EXPECT().GetRequest(ctx, tc.userID, tc.reqID).Return(tc.repReq, tc.repErr).Times(1)

			gotReq, gotErr := srvc.GetRequest(ctx, tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotReq, tc.wantReq)
		})
	}
}

func loadImage(t *testing.T, path string) []byte {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening an image", err)
	}

	return b
}

func TestAddReqeust(t *testing.T) {
	pngTestImage := loadImage(t, "test_data/x.png")

	testCases := []struct {
		testName        string
		userID          int
		file            *bytes.Buffer
		fileName        string
		imageID         int
		imageRepoErr    error
		imageURL        string
		storageErr      error
		convInfo        model.ConversionInfo
		runUploadFile   bool
		runAddImage     bool
		runAddRequest   bool
		repoReqID       int
		reqRepoErr      error
		runProcessImage bool
		wantReqID       int
		wantErr         error
	}{
		{
			testName: "all is good",
			userID:   123,
			file:     bytes.NewBuffer(pngTestImage),
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			runUploadFile:   true,
			runAddImage:     true,
			runAddRequest:   true,
			reqRepoErr:      nil,
			repoReqID:       15,
			runProcessImage: true,
			wantReqID:       15,
			wantErr:         nil,
		},
		{
			testName: "unknow file type",
			userID:   123,
			file:     bytes.NewBuffer(pngTestImage),
			fileName: "filename.webm",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			reqRepoErr: nil,
			repoReqID:  15,
			wantReqID:  0,
			wantErr:    &service.UnsupportedTypeError{"webm"},
		},
		{
			testName: "storage error",
			userID:   123,
			file:     bytes.NewBuffer(pngTestImage),
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			storageErr:    storage.ErrBucketNotExist,
			runUploadFile: true,
			reqRepoErr:    nil,
			repoReqID:     15,
			wantReqID:     0,
			wantErr:       storage.ErrBucketNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequestRepo(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)
			mockProcess := mocks.NewMockImageProcesser(mockCtr)

			srvc := service.NewRequest(mockRequest, mockStorage, mockProcess)
			ctx := context.Background()

			if tc.runUploadFile {
				mockStorage.EXPECT().UploadFile(ctx, tc.userID, tc.fileName, tc.file.Bytes()).Return(tc.imageURL, tc.storageErr)
			}

			if tc.runAddImage {
				mockRequest.EXPECT().
					AddImageAndRequest(ctx, tc.userID, gomock.Any(), gomock.Any()).
					Return(tc.repoReqID, tc.imageRepoErr)
			}

			if tc.runProcessImage {
				mockProcess.EXPECT().ProcessImage(ctx, gomock.Any())
			}

			gotReqID, gotErr := srvc.AddRequest(ctx, tc.userID, tc.file,
				tc.fileName, tc.convInfo)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotReqID, tc.wantReqID)
		})
	}
}

func TestDeleteReqeust(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		reqID    int
		url1     string
		url2     string
		initMock func(*mocks.MockRequestRepo, *mocks.MockStorager, int, int, string, string)
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   1,
			reqID:    2,
			url1:     "first image url",
			url2:     "second image url",
			initMock: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
				mRep.EXPECT().DeleteRequestAndImage(gomock.Any(), userID, reqID).Return(url1, url2, nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url1).Return(nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url2).Return(nil)
			},
			wantErr: nil,
		},
		{
			testName: "error while deleting request and images",
			userID:   1,
			reqID:    2,
			url1:     "",
			url2:     "",
			initMock: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
				mRep.EXPECT().DeleteRequestAndImage(gomock.Any(), userID, reqID).Return(url1, url2, errRepository)
			},
			wantErr: errRepository,
		},
		{
			testName: "error while deleting first file",
			userID:   1,
			reqID:    2,
			url1:     "first image url",
			url2:     "second image url",
			initMock: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
				mRep.EXPECT().DeleteRequestAndImage(gomock.Any(), userID, reqID).Return(url1, url2, nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url1).Return(errStorage)
			},
			wantErr: errStorage,
		},
		{
			testName: "second image not exists",
			userID:   1,
			reqID:    2,
			url1:     "first image id",
			url2:     "",
			initMock: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
				mRep.EXPECT().DeleteRequestAndImage(gomock.Any(), userID, reqID).Return(url1, url2, nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url1).Return(nil)
			},
			wantErr: nil,
		},
		{
			testName: "error while deleting second file",
			userID:   1,
			reqID:    2,
			url1:     "first image url",
			url2:     "second image url",
			initMock: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
				mRep.EXPECT().DeleteRequestAndImage(gomock.Any(), userID, reqID).Return(url1, url2, nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url1).Return(nil)
				mStor.EXPECT().DeleteFile(gomock.Any(), url2).Return(errStorage)
			},
			wantErr: errStorage,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequestRepo(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			tc.initMock(mockRequest, mockStorage, tc.userID, tc.reqID, tc.url1, tc.url2)

			srvc := service.NewRequest(mockRequest, mockStorage, &mocks.MockImageProcesser{})
			ctx := context.Background()

			gotErr := srvc.DeleteRequest(ctx, tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
