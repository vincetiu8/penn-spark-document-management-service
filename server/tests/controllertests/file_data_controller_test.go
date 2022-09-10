package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/invincibot/penn-spark-server/api/controllers"
	"github.com/invincibot/penn-spark-server/api/models"
	"github.com/invincibot/penn-spark-server/tests/util"
)

func TestCreateFileData(t *testing.T) {
	testServer.RefreshFileSystem()

	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	file := testServer.Data.Files[0]

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", file.Name)
	require.NoError(t, err)
	_, err = io.Copy(fw, bytes.NewBufferString("text"))
	require.NoError(t, err)
	err = w.Close()
	require.NoError(t, err)

	testCases := []struct {
		id          uint
		user        models.User
		statusCode  int
		expectedErr error
	}{
		{
			id:         file.ID,
			user:       users[0],
			statusCode: http.StatusOK,
		},
		{
			id:          file.ID,
			user:        users[1],
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			id:          999,
			user:        users[0],
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileNotFound,
		},
	}

	for i, testCase := range testCases {
		if i == len(testCases)-1 {
			err = models.DeleteAccessRole(testServer.Server.DB, testServer.Data.AccessRoles[0].ID)
			require.NoError(t, err)
		}

		data := b
		req, err := http.NewRequest("PUT", "/file-data", &data)
		require.NoError(t, err)
		req.Header.Set("Content-Type", w.FormDataContentType())
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.CreateFileData(rr, req, testCase.user)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				util.CheckFilesEqual(t, file, responseMap)
				returnedData, err := afero.ReadFile(testServer.Server.FileSystem,
					testServer.Server.FileSystem.FilePath+"/"+strconv.Itoa(int(file.ID)))
				require.NoError(t, err)
				assert.Equal(t, []byte("text"), returnedData)
			case http.StatusBadRequest:
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}

func TestGetFileData(t *testing.T) {
	testServer.RefreshFileSystem()

	err := testServer.SeedData()
	require.NoError(t, err)

	users := testServer.Data.Users
	file := testServer.Data.Files[1]

	_, err = models.UpdateAccessRole(testServer.Server.DB, models.AccessRole{
		ID:          users[0].UserRoles[0].AccessRoles[1].ID,
		AccessLevel: models.Publisher,
	})

	err = afero.WriteFile(
		testServer.Server.FileSystem,
		testServer.Server.FileSystem.FilePath+"/"+strconv.Itoa(int(file.ID)),
		[]byte("text"),
		0666,
	)
	require.NoError(t, err)

	testCases := []struct {
		id          uint
		user        models.User
		statusCode  int
		expectedErr error
	}{
		{
			id:         file.ID,
			user:       users[0],
			statusCode: http.StatusOK,
		},
		{
			id:          file.ID,
			user:        users[2],
			statusCode:  http.StatusForbidden,
			expectedErr: controllers.ErrUserForbidden,
		},
		{
			id:          999,
			user:        users[1],
			statusCode:  http.StatusBadRequest,
			expectedErr: models.ErrFileNotFound,
		},
		{
			id:         file.ID,
			user:       users[1],
			statusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		req, err := http.NewRequest("GET", "/file-data", nil)
		require.NoError(t, err)
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprint(testCase.id)})
		rr := httptest.NewRecorder()
		testServer.Server.GetFileData(rr, req, testCase.user)

		require.NoError(t, err)
		if assert.Equal(t, testCase.statusCode, rr.Code) {
			switch testCase.statusCode {
			case http.StatusOK:
				assert.Equal(t, []byte("text"), rr.Body.Bytes())
			case http.StatusBadRequest:
				responseMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
				assert.Equal(t, testCase.expectedErr.Error(), responseMap["error"])
			}
		}
	}
}
