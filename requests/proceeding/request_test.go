package proceeding

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

func TestGetProceedingData(t *testing.T) {
	globalvars.GetProceedingRequestUrl = "https://fake/proceeding/%s"
	globalvars.HomePageUrl = "https://fake/home"

	sampleProceeding := models.ActiveProceeding{
		ProceedingsID:     "abc123",
		ForeignerFullName: "Case 1",
	}

	testCases := []struct {
		name         string
		client       *http.Client
		sessionToken string
		proceeding   models.ActiveProceeding
		wantData     *models.DetailedProceedingData
		wantErrStr   string
		wantErrType  any
	}{
		{
			name:       "nil client",
			client:     nil,
			proceeding: sampleProceeding,
			wantErrStr: "HTTP client is nil",
		},
		{
			name: "unauthorized",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, "https://fake/proceeding/abc123", req.URL.String())
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
			}),
			sessionToken: "invalid-token",
			proceeding:   sampleProceeding,
			wantErrStr:   "GetProceedingData failed because of unauthorized status code: 401 Unauthorized",
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
			wantErrStr:   "JSON parcing error",
		},
		{
			name: "success",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				data := models.DetailedProceedingData{
					ID:               "b0200708-0823-46c5-a4c2-839df2141c98",
					CircumstanceText: "Wykonywanie pracy",
					Status:           "1",
				}
				respBytes, _ := json.Marshal(data)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
				}
			}),
			sessionToken: "valid-token",
			proceeding:   sampleProceeding,
			wantData: &models.DetailedProceedingData{
				ID:               "b0200708-0823-46c5-a4c2-839df2141c98",
				CircumstanceText: "Wykonywanie pracy",
				Status:           "1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetProceedingData(tc.client, tc.sessionToken, tc.proceeding)

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
				assert.Equal(t, tc.wantData, result)
			}
		})
	}
}
