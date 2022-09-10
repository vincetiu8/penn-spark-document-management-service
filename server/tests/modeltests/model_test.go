package modeltests

import (
	"os"
	"testing"

	"github.com/invincibot/penn-spark-server/tests/util"
)

var testServer *util.TestServer

func TestMain(m *testing.M) {
	testServer = util.NewTestServer()

	os.Exit(m.Run())
}
