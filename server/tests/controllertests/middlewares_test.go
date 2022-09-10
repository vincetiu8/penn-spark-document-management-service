package controllertests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/auth"
	"github.com/invincibot/penn-spark-server/api/controllers"
	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/api/responses"
)

func EmptyControllerFunc(w http.ResponseWriter, _ *http.Request, _ models.User) {
	responses.JSON(w, http.StatusNoContent, "")
}

func TestSetMiddlewareAuthentication(t *testing.T) {
	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users

	tokens := []string{}
	for _, user := range users {
		token, err := auth.CreateToken(user.ID, time.Minute*15)
		require.NoError(t, err)
		token = fmt.Sprintf("Bearer %v", token)
		tokens = append(tokens, token)
	}

	testCases := []struct {
		id          uint
		token       string
		isAdmin     bool
		statusCode  int
		expectedErr error
	}{
		{
			token:       "",
			statusCode:  http.StatusUnauthorized,
			expectedErr: controllers.ErrUserUnauthorized,
		},
		{
			token:       tokens[1],
			isAdmin:     true,
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			token:      tokens[0],
			isAdmin:    true,
			statusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)
		if testCase.id > 0 {
			req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		}
		req.Header.Set("Authorization", testCase.token)
		rr := httptest.NewRecorder()
		middlewareFunc := controllers.SetMiddlewareAuthentication(EmptyControllerFunc, &testServer.Server, testCase.isAdmin)
		middlewareFunc(rr, req)

		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusUnauthorized, http.StatusForbidden:
				responseMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}
