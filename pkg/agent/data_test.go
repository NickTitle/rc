package agent

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/nicktitle/rc/pkg/fileformat"
	"github.com/stretchr/testify/require"
)

func TestSendData(t *testing.T) {
	rand.Seed(time.Now().Unix())
	// this can definitely conflict but should be relatively safe due to the range
	sourcePort := rand.Intn(20000) + 2000
	byteCount := rand.Intn(200)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)

		defer req.Body.Close()

		require.Equal(t, byteCount, len(body))
	}))
	defer srv.Close()

	srvURL, err := url.Parse(srv.URL)
	require.NoError(t, err)

	srvPort, err := strconv.Atoi(srvURL.Port())
	require.NoError(t, err)

	settings := fileformat.ConnectionSettings{
		Source:      fileformat.ConnectionAddress{Host: "localhost", Port: sourcePort},
		Destination: fileformat.ConnectionAddress{Host: srvURL.Hostname(), Port: srvPort},
		PayloadSize: byteCount,
		Scheme:      "http",
	}

	event, err := sendData(settings)
	require.NoError(t, err)

	require.Equal(t, event["bytes_sent"], byteCount)
	require.Equal(t, event["source_port"], sourcePort)
}
