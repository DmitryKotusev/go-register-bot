package dates

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

func TestGetReservationQueueDates(t *testing.T) {
	globalvars.GetReservationQueueDatesRequestUrl = "https://fake/queue/%s/dates"
	globalvars.HomePageCasesUrl = "https://fake/cases/%s"

	sampleProceeding := models.DetailedProceedingData{
		ID: "abc123",
	}

	sampleQueue := models.ReservationQueue{
		ID:     "queue001",
		Prefix: "Test Queue",
	}

	testCases := []struct {
		name         string
		client       *http.Client
		sessionToken string
		proceeding   *models.DetailedProceedingData
		queue        models.ReservationQueue
		wantDates    []string
		wantErrStr   string
		wantErrType  any
	}{
		{
			name:       "nil client",
			client:     nil,
			wantErrStr: "GetReservationQueueDates, HTTP client is nil",
		},
		{
			name:       "nil proceeding",
			client:     &http.Client{},
			proceeding: nil,
			wantErrStr: "GetReservationQueueDates, proceeding data is nil",
		},
		{
			name: "successful response",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				body, _ := json.Marshal([]string{"2025-08-10T00:00:00", "2025-08-12T00:00:00"})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}
			}),
			sessionToken: "valid-token",
			proceeding:   &sampleProceeding,
			queue:        sampleQueue,
			wantDates:    []string{"2025-08-10", "2025-08-12"},
		},
		{
			name: "unauthorized",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			}),
			sessionToken: "invalid-token",
			proceeding:   &sampleProceeding,
			queue:        sampleQueue,
			wantErrStr:   "❌ GetReservationQueueDates failed because of unauthorized status code",
			wantErrType:  &modelerrors.UnauthorizedError{},
		},
		{
			name: "500 error",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			}),
			sessionToken: "valid-token",
			proceeding:   &sampleProceeding,
			queue:        sampleQueue,
			wantErrStr:   "GetReservationQueueDates request failed with status: 500 Internal Server Error",
		},
		{
			name: "bad json",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("{invalid_json"))),
				}
			}),
			sessionToken: "valid-token",
			proceeding:   &sampleProceeding,
			queue:        sampleQueue,
			wantErrStr:   "GetReservationQueueDates body JSON parcing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dates, err := GetReservationQueueDates(tc.client, tc.sessionToken, tc.proceeding, tc.queue)

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
				assert.Equal(t, tc.wantDates, dates)
			}
		})
	}
}
