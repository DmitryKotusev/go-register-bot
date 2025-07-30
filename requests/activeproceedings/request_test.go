package activeproceedings

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

func TestGetActiveProceedings(t *testing.T) {
	globalvars.GetActiveProceedingsRequestUrl = "https://fake/active"
	globalvars.HomePageUrl = "https://fake/home"

	testCases := []struct {
		name         string
		client       *http.Client
		sessionToken string
		wantList     []models.ActiveProceeding
		wantErrStr   string
		wantErrType  any
	}{
		{
			name:       "nil client",
			client:     nil,
			wantErrStr: "HTTP client is nil",
		},
		{
			name: "unauthorized",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Status:     "401 Unauthorized",
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
			}),
			sessionToken: "invalid-token",
			wantErrStr:   "GetActiveProceedings failed because of unauthorized status code: 401 Unauthorized",
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
			wantErrStr:   "JSON parcing error",
		},
		{
			name: "success",
			client: test_utils.NewTestClient(func(req *http.Request) *http.Response {
				data := []models.ActiveProceeding{
					{ProceedingsID: "abc123", ForeignerFullName: "Case 1"},
					{ProceedingsID: "def456", ForeignerFullName: "Case 2"},
				}
				respBytes, _ := json.Marshal(data)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(respBytes)),
				}
			}),
			sessionToken: "valid-token",
			wantList: []models.ActiveProceeding{
				{ProceedingsID: "abc123", ForeignerFullName: "Case 1"},
				{ProceedingsID: "def456", ForeignerFullName: "Case 2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetActiveProceedings(tc.client, tc.sessionToken)

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
