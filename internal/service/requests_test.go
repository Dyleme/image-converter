package service_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	errRepository = errors.New("error in repository")
	errStorage    = errors.New("error in storage")
	errProc       = errors.New("error while processing")
)

func TestRequest_GetRequests(t *testing.T) {
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

func TestRequest_GetRequest(t *testing.T) {
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

func TestRequest_AddReqeust(t *testing.T) {
	pngTestImage := []byte{1, 2, 3, 4, 5}
	url := "url"
	testCases := []struct {
		testName  string
		userID    int
		fileName  string
		convInfo  model.ConversionInfo
		configure func(*mocks.MockRequestRepo, *mocks.MockStorager, *mocks.MockImageProcesser)
		wantReqID int
		wantErr   error
	}{
		{
			testName: "all is good",
			userID:   123,
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			configure: func(mrr *mocks.MockRequestRepo, ms *mocks.MockStorager, mip *mocks.MockImageProcesser) {
				ctx := context.Background()
				filename := "filename.png"
				imageInfo := &model.ReuquestImageInfo{
					URL:  url,
					Type: "png",
				}
				reqRepoID := 15

				ms.EXPECT().UploadFile(ctx, 123, filename, pngTestImage).Return(url, nil)
				mrr.EXPECT().
					AddImageAndRequest(ctx, 123, imageInfo, gomock.Any()).
					Return(reqRepoID, nil)

				infoToProc := &model.RequestToProcess{
					ReqID:    reqRepoID,
					FileName: filename,
				}
				mip.EXPECT().ProcessImage(ctx, infoToProc).Return(nil)
			},
			wantReqID: 15,
			wantErr:   nil,
		},
		{
			testName: "unknow file type",
			userID:   123,
			fileName: "filename.webm",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			configure: func(mrr *mocks.MockRequestRepo, ms *mocks.MockStorager, mip *mocks.MockImageProcesser) {},
			wantReqID: 0,
			wantErr:   service.UnsupportedTypeError{"webm"},
		},
		{
			testName: "storage error",
			userID:   123,
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			configure: func(mrr *mocks.MockRequestRepo, ms *mocks.MockStorager, mip *mocks.MockImageProcesser) {
				ms.EXPECT().UploadFile(context.Background(), 123, "filename.png", pngTestImage).Return("", errStorage)
			},
			wantReqID: 0,
			wantErr:   errStorage,
		},
		{
			testName: "repo error",
			userID:   123,
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			configure: func(mrr *mocks.MockRequestRepo, ms *mocks.MockStorager, mip *mocks.MockImageProcesser) {
				ctx := context.Background()
				ms.EXPECT().UploadFile(context.Background(), 123, "filename.png", pngTestImage).Return(url, nil)
				imageInfo := &model.ReuquestImageInfo{
					URL:  url,
					Type: "png",
				}
				mrr.EXPECT().
					AddImageAndRequest(ctx, 123, imageInfo, gomock.Any()).
					Return(0, errRepository)
			},
			wantReqID: 0,
			wantErr:   errRepository,
		},
		{
			testName: "process error",
			userID:   123,
			fileName: "filename.png",
			convInfo: model.ConversionInfo{
				Ratio: 0.5,
				Type:  "png",
			},
			configure: func(mrr *mocks.MockRequestRepo, ms *mocks.MockStorager, mip *mocks.MockImageProcesser) {
				ctx := context.Background()
				filename := "filename.png"
				imageInfo := &model.ReuquestImageInfo{
					URL:  url,
					Type: "png",
				}
				reqRepoID := 15

				ms.EXPECT().UploadFile(ctx, 123, filename, pngTestImage).Return(url, nil)
				mrr.EXPECT().
					AddImageAndRequest(ctx, 123, imageInfo, gomock.Any()).
					Return(reqRepoID, nil)

				infoToProc := &model.RequestToProcess{
					ReqID:    reqRepoID,
					FileName: filename,
				}
				mip.EXPECT().ProcessImage(ctx, infoToProc).Return(errProc)
			},
			wantReqID: 0,
			wantErr:   errProc,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequestRepo(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)
			mockProcess := mocks.NewMockImageProcesser(mockCtr)

			tc.configure(mockRequest, mockStorage, mockProcess)

			srvc := service.NewRequest(mockRequest, mockStorage, mockProcess)
			ctx := context.Background()

			gotReqID, gotErr := srvc.AddRequest(ctx, tc.userID, bytes.NewBuffer(pngTestImage),
				tc.fileName, tc.convInfo)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotReqID, tc.wantReqID)
		})
	}
}

func TestRequest_DeleteReqeust(t *testing.T) {
	testCases := []struct {
		testName  string
		userID    int
		reqID     int
		url1      string
		url2      string
		configure func(*mocks.MockRequestRepo, *mocks.MockStorager, int, int, string, string)
		wantErr   error
	}{
		{
			testName: "all is good",
			userID:   1,
			reqID:    2,
			url1:     "first image url",
			url2:     "second image url",
			configure: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
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
			configure: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
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
			configure: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
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
			configure: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
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
			configure: func(mRep *mocks.MockRequestRepo, mStor *mocks.MockStorager, userID, reqID int, url1, url2 string) {
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

			tc.configure(mockRequest, mockStorage, tc.userID, tc.reqID, tc.url1, tc.url2)

			srvc := service.NewRequest(mockRequest, mockStorage, &mocks.MockImageProcesser{})
			ctx := context.Background()

			gotErr := srvc.DeleteRequest(ctx, tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
