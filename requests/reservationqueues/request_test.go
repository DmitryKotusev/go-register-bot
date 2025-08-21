package reservationqueues

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

func TestGetReservationQueues(t *testing.T) {
	globalvars.GetProceedingReservationQueuesRequestUrl = "https://fake/proceedings/%s/queues"
	globalvars.HomePageCasesUrl = "https://fake/cases/%s"

	sampleProceeding := models.DetailedProceedingData{
		ID: "12345",
	}

	testCases := []struct {
		name         string
		client       *http.Client
		sessionToken string
		proceeding   *models.DetailedProceedingData
		wantList     []models.ReservationQueue
		wantErrStr   string
		wantErrType  any
	}{
		{
			name:       "nil client",
			client:     nil,
			wantErrStr: "GetReservationQueues, HTTP client is nil",
		},
		{
			name:       "nil proceeding",
			client:     &http.Client{},
			proceeding: nil,
			wantErrStr: "GetReservationQueues, proceeding data is nil",
		},
		{
			name: "successful response",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				resp := []models.ReservationQueue{
					{ID: "q1", Prefix: "Queue 1"},
					{ID: "q2", Prefix: "Queue 2"},
				}
				respBytes, _ := json.Marshal(resp)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
					Header:     make(http.Header),
				}
			}),
			sessionToken: "token123",
			proceeding:   &sampleProceeding,
			wantList: []models.ReservationQueue{
				{ID: "q1", Prefix: "Queue 1"},
				{ID: "q2", Prefix: "Queue 2"},
			},
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
			sessionToken: "badtoken",
			proceeding:   &sampleProceeding,
			wantErrStr:   "❌ GetReservationQueues failed because of unauthorized status code",
			wantErrType:  &modelerrors.UnauthorizedError{},
		},
		{
			name: "internal server error",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(bytes.NewReader([]byte(""))),
				}
			}),
			sessionToken: "token123",
			proceeding:   &sampleProceeding,
			wantErrStr:   "GetReservationQueues request failed with status: 500 Internal Server Error",
		},
		{
			name: "bad json",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("{bad json"))),
				}
			}),
			sessionToken: "token123",
			proceeding:   &sampleProceeding,
			wantErrStr:   "GetReservationQueues body JSON parcing error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			queues, err := GetReservationQueues(tc.client, tc.sessionToken, tc.proceeding)

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
				assert.Equal(t, tc.wantList, queues)
			}
		})
	}
}
