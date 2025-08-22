package reserve

import (
	"bot-main/globalvars"
	"bot-main/models"
	modelerrors "bot-main/models/errors"
	test_utils "bot-main/tests/utils"
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sampleProceeding() *models.DetailedProceedingData {
	return &models.DetailedProceedingData{
		ID: "proc-123",
		Person: models.Person{
			FirstName:   "Bob",
			Surname:     "Smith",
			DateOfBirth: "1990-01-01T00:00:00Z",
		},
	}
}

func sampleSlot() models.Slot {
	return models.Slot{
		ID:   42,
		Date: "2025-10-03T11:30:00",
	}
}

func sampleQueue() models.ReservationQueue {
	return models.ReservationQueue{ID: "queue-1"}
}

func TestReserveDateSlot(t *testing.T) {
	globalvars.ReserveAppointmentRequestUrl = "https://fake/reservation/%s"
	globalvars.HomePageCasesUrl = "https://fake/cases/%s"

	testCases := []struct {
		name        string
		client      *http.Client
		session     string
		proceeding  *models.DetailedProceedingData
		queue       models.ReservationQueue
		slot        models.Slot
		wantErrStr  string
		wantErrType any
	}{
		{
			name:       "nil client",
			client:     nil,
			session:    "tok",
			proceeding: sampleProceeding(),
			queue:      sampleQueue(),
			slot:       sampleSlot(),
			wantErrStr: "HTTP client is nil",
		},
		{
			name:       "nil proceeding",
			client:     &http.Client{},
			session:    "tok",
			proceeding: nil,
			queue:      sampleQueue(),
			slot:       sampleSlot(),
			wantErrStr: "proceeding data is nil",
		},
		{
			name: "successful reservation",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				// Проверим метод и заголовки
				if req.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", req.Method)
				}
				if req.Header.Get("Authorization") != "Bearer tok" {
					t.Errorf("expected Authorization header")
				}
				resp := `{"ok":true}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(resp))),
					Header:     make(http.Header),
				}
			}),
			session:    "tok",
			proceeding: sampleProceeding(),
			queue:      sampleQueue(),
			slot:       sampleSlot(),
		},
		{
			name: "unauthorized",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte("unauthorized"))),
				}
			}),
			session:     "tok",
			proceeding:  sampleProceeding(),
			queue:       sampleQueue(),
			slot:        sampleSlot(),
			wantErrStr:  "unauthorized",
			wantErrType: &modelerrors.UnauthorizedError{},
		},
		{
			name: "forbidden",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusForbidden,
					Status:     "403 Forbidden",
					Body:       io.NopCloser(bytes.NewReader([]byte("forbidden"))),
				}
			}),
			session:     "tok",
			proceeding:  sampleProceeding(),
			queue:       sampleQueue(),
			slot:        sampleSlot(),
			wantErrStr:  "forbidden",
			wantErrType: &modelerrors.ForbiddenError{},
		},
		{
			name: "server error",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusTeapot,
					Status:     "418 I'm a teapot",
					Body:       io.NopCloser(bytes.NewReader([]byte("teapot"))),
				}
			}),
			session:    "tok",
			proceeding: sampleProceeding(),
			queue:      sampleQueue(),
			slot:       sampleSlot(),
			wantErrStr: "request for 2025-10-03T11:30:00 failed with status: 418 I'm a teapot",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ReserveDateSlot(tc.client, tc.session, tc.proceeding, tc.queue, tc.slot)

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
			}
		})
	}
}
