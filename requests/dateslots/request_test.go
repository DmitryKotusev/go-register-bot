package dateslots

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	test_utils "bot-main/tests/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReservationQueueDateSlots(t *testing.T) {
	globalvars.GetReservationQueueDateSlotsRequestUrl = "https://fake/queue/%s/%s/slots"
	globalvars.HomePageCasesUrl = "https://fake/cases/%s"

	sampleProceeding := &models.DetailedProceedingData{
		ID: "proc123",
	}
	sampleQueue := models.ReservationQueue{
		ID: "queue456",
	}

	testCases := []struct {
		name         string
		client       *http.Client
		sessionToken string
		proceeding   *models.DetailedProceedingData
		queue        models.ReservationQueue
		date         string
		wantList     []models.Slot
		wantErrStr   string
		wantErrType  any
	}{
		{
			name:       "nil client",
			client:     nil,
			proceeding: sampleProceeding,
			queue:      sampleQueue,
			date:       "2025-08-21",
			wantErrStr: "HTTP client is nil",
		},
		{
			name:       "nil proceeding",
			client:     &http.Client{},
			proceeding: nil,
			queue:      sampleQueue,
			date:       "2025-08-21",
			wantErrStr: "proceeding data is nil",
		},
		{
			name: "unauthorized",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "https://fake/queue/queue456/2025-08-21/slots", req.URL.String())
				assert.Equal(t, "https://fake/cases/proc123", req.Header.Get("Referer"))
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
			}),
			sessionToken: "invalid-token",
			proceeding:   sampleProceeding,
			queue:        sampleQueue,
			date:         "2025-08-21",
			wantErrStr:   "GetReservationQueueDateSlots failed because of unauthorized status code: 401 Unauthorized",
			wantErrType:  &modelerrors.UnauthorizedError{},
		},
		{
			name: "server error",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
			}),
			sessionToken: "some-token",
			proceeding:   sampleProceeding,
			queue:        sampleQueue,
			date:         "2025-08-21",
			wantErrStr:   "500 Internal Server Error",
		},
		{
			name: "bad json",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("{invalid json"))),
				}
			}),
			sessionToken: "some-token",
			proceeding:   sampleProceeding,
			queue:        sampleQueue,
			date:         "2025-08-21",
			wantErrStr:   "JSON parcing error",
		},
		{
			name: "success",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				data := []models.Slot{
					{ID: 111, Date: "2025-08-21T08:40:00", Count: 1},
					{ID: 222, Date: "2025-08-21T09:00:00", Count: 2},
				}
				respBytes, _ := json.Marshal(data)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
				}
			}),
			sessionToken: "valid-token",
			proceeding:   sampleProceeding,
			queue:        sampleQueue,
			date:         "2025-08-21",
			wantList: []models.Slot{
				{ID: 111, Date: "2025-08-21T08:40:00", Count: 1},
				{ID: 222, Date: "2025-08-21T09:00:00", Count: 2},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetReservationQueueDateSlots(tc.client, tc.sessionToken, tc.proceeding, tc.queue, tc.date)

			if tc.wantErrStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErrStr)
				if tc.wantErrType != nil {
					assert.Truef(t,
						errors.As(err, tc.wantErrType),
						"expected error of type %T but got %T",
						tc.wantErrType, err,
					)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantList, result)
			}
		})
	}
}
