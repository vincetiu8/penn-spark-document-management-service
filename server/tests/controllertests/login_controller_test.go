package controllertests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/tests/util"
)

func TestLogin(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)
	user := util.Users[0]

	testCases := []struct {
		inputUser   models.User
		statusCode  int
		expectedErr error
	}{
		{
			inputUser: models.User{
				Username: user.Username,
				Password: user.Password,
			},
			statusCode:  http.StatusOK,
			expectedErr: nil,
		},
		{
			inputUser: models.User{
				Username: "not a username",
				Password: user.Password,
			},
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrUserNotFound,
		},
	}

	for _, testCase := range testCases {
		inputJSON := util.UserToJSON(testCase.inputUser, "login")
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(inputJSON))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(testServer.Server.Login)
		handler.ServeHTTP(rr, req)

		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				assert.NotEqual(t, "", rr.Body.String())
			case http.StatusBadRequest, http.StatusUnprocessableEntity:
				responseMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}
