package service_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/service/mocks"
	"github.com/Dyleme/image-coverter/internal/storage"
	"github.com/golang/mock/gomock"
)

var (
	errRepository = errors.New("error in repository")
	errStorage    = errors.New("error in storage")
)

type mockSender struct{}

func (m *mockSender) ProcessImage(data *model.ConversionData) {}

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
			mockRequest := mocks.NewMockRequester(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			srvc := service.NewRequestService(mockRequest, mockStorage, mocks.NewMockImageProcesser(mockCtr))
			ctx := context.Background()

			mockRequest.EXPECT().GetRequests(ctx, tc.userID).Return(tc.repReqs, tc.repErr)

			gotReqs, gotErr := srvc.GetRequests(ctx, tc.userID)
			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("want err: %v, got err: %v", tc.wantErr, gotErr)
			}

			if reflect.DeepEqual(gotReqs, tc.wantErr) {
				t.Errorf("want reqs: %v, got reqs %v", tc.wantReqs, gotReqs)
			}
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
			wantReq:  &model.Request{},
			wantErr:  errRepository,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequester(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			srvc := service.NewRequestService(mockRequest, mockStorage, &mockSender{})
			ctx := context.Background()

			mockRequest.EXPECT().GetRequest(ctx, tc.userID, tc.reqID).Return(tc.repReq, tc.repErr).Times(1)

			gotReqs, gotErr := srvc.GetRequest(ctx, tc.userID, tc.reqID)
			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("want err: %v, got err: %v", tc.wantErr, gotErr)
			}

			if reflect.DeepEqual(gotReqs, tc.wantErr) {
				t.Errorf("want reqs: %v, got reqs %v", tc.wantReq, gotReqs)
			}
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
			wantErr:    service.ErrUnsupportedType,
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
			mockRequest := mocks.NewMockRequester(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)
			mockProcess := mocks.NewMockImageProcesser(mockCtr)

			srvc := service.NewRequestService(mockRequest, mockStorage, mockProcess)
			ctx := context.Background()

			if tc.runUploadFile {
				mockStorage.EXPECT().UploadFile(ctx, tc.userID, tc.fileName, tc.file.Bytes()).Return(tc.imageURL, tc.storageErr)
			}

			if tc.runAddImage {
				mockRequest.EXPECT().AddImage(ctx, tc.userID, gomock.Any()).Return(tc.imageID, tc.imageRepoErr)
			}

			if tc.runAddRequest {
				mockRequest.EXPECT().AddRequest(ctx, gomock.Any(), tc.userID).Return(tc.repoReqID, tc.reqRepoErr)
			}
			if tc.runProcessImage {
				mockProcess.EXPECT().ProcessImage(gomock.Any())
			}

			gotReqID, gotErr := srvc.AddRequest(ctx, tc.userID, tc.file,
				tc.fileName, tc.convInfo)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("want error: %v, got error %v", tc.wantErr, gotErr)
			}

			if gotReqID != tc.wantReqID {
				t.Errorf("want reqID: %v, got reqID %v", tc.wantReqID, gotReqID)
			}
		})
	}
}

func TestDeleteReqeust(t *testing.T) {
	testCases := []struct {
		testName        string
		userID          int
		reqID           int
		repo1ID         int
		repo2ID         int
		deleteReqErr    error
		runDeleteImage1 bool
		url1            string
		deleteIm1RepErr error
		runDeleteImage2 bool
		url2            string
		deleteIm2RepErr error
		runDeleteFile1  bool
		deleteFile1Err  error
		runDeleteFile2  bool
		deleteFile2Err  error
		wantErr         error
	}{
		{
			testName:        "all is good",
			userID:          1,
			reqID:           2,
			repo1ID:         3,
			repo2ID:         4,
			deleteReqErr:    nil,
			runDeleteImage1: true,
			url1:            "first image url",
			deleteIm1RepErr: nil,
			runDeleteImage2: true,
			url2:            "second image url",
			deleteIm2RepErr: nil,
			runDeleteFile1:  true,
			deleteFile1Err:  nil,
			runDeleteFile2:  true,
			deleteFile2Err:  nil,
			wantErr:         nil,
		},
		{
			testName:        "error while deleting request",
			userID:          1,
			reqID:           2,
			repo1ID:         0,
			repo2ID:         0,
			deleteReqErr:    errRepository,
			runDeleteImage1: false,
			runDeleteImage2: false,
			runDeleteFile1:  false,
			runDeleteFile2:  false,
			wantErr:         errRepository,
		},
		{
			testName:        "error while deleting first image",
			userID:          1,
			reqID:           2,
			repo1ID:         3,
			repo2ID:         4,
			deleteReqErr:    nil,
			runDeleteImage1: true,
			deleteIm1RepErr: errRepository,
			runDeleteImage2: false,
			runDeleteFile1:  false,
			runDeleteFile2:  false,
			wantErr:         errRepository,
		},
		{
			testName:        "error while deleting second image",
			userID:          1,
			reqID:           2,
			repo1ID:         3,
			repo2ID:         4,
			deleteReqErr:    nil,
			runDeleteImage1: true,
			url1:            "first image url",
			deleteIm1RepErr: nil,
			runDeleteImage2: true,
			deleteIm2RepErr: errRepository,
			runDeleteFile1:  false,
			runDeleteFile2:  false,
			wantErr:         errRepository,
		},
		{
			testName:        "error while deleting first file",
			userID:          1,
			reqID:           2,
			repo1ID:         3,
			repo2ID:         4,
			deleteReqErr:    nil,
			runDeleteImage1: true,
			url1:            "first image url",
			deleteIm1RepErr: nil,
			runDeleteImage2: true,
			url2:            "second image url",
			deleteIm2RepErr: nil,
			runDeleteFile1:  true,
			deleteFile1Err:  errStorage,
			runDeleteFile2:  false,
			wantErr:         errStorage,
		},
		{
			testName:        "error while deleting second file",
			userID:          1,
			reqID:           2,
			repo1ID:         3,
			repo2ID:         4,
			deleteReqErr:    nil,
			runDeleteImage1: true,
			url1:            "first image url",
			deleteIm1RepErr: nil,
			runDeleteImage2: true,
			url2:            "second image url",
			deleteIm2RepErr: nil,
			runDeleteFile1:  true,
			deleteFile1Err:  nil,
			runDeleteFile2:  true,
			deleteFile2Err:  errStorage,
			wantErr:         errStorage,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockRequest := mocks.NewMockRequester(mockCtr)
			mockStorage := mocks.NewMockStorager(mockCtr)

			srvc := service.NewRequestService(mockRequest, mockStorage, &mockSender{})
			ctx := context.Background()

			mockRequest.EXPECT().DeleteRequest(ctx, tc.userID, tc.reqID).Return(tc.repo1ID, tc.repo2ID, tc.deleteReqErr)

			if tc.runDeleteImage1 {
				mockRequest.EXPECT().DeleteImage(ctx, tc.userID, tc.repo1ID).Return(tc.url1, tc.deleteIm1RepErr)
			}

			if tc.runDeleteImage2 {
				mockRequest.EXPECT().DeleteImage(ctx, tc.userID, tc.repo2ID).Return(tc.url2, tc.deleteIm2RepErr)
			}

			if tc.runDeleteFile1 {
				mockStorage.EXPECT().DeleteFile(ctx, tc.url1).Return(tc.deleteFile1Err)
			}

			if tc.runDeleteFile2 {
				mockStorage.EXPECT().DeleteFile(ctx, tc.url2).Return(tc.deleteFile2Err)
			}

			gotErr := srvc.DeleteRequest(ctx, tc.userID, tc.reqID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("want error: %v, got error %v", tc.wantErr, gotErr)
			}
		})
	}
}
