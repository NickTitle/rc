package agent

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nicktitle/rc/pkg/fileformat"
)

func TestStartProcess(t *testing.T) {
	rand.Seed(time.Now().Unix())

	// generating random strings is a pain in go, so we'll generate a timestamped tempfile instead
	tmpFileName := fmt.Sprintf("testFile-%s", strconv.Itoa(rand.Intn(20000)))
	tmpFilePath := path.Join(os.TempDir(), tmpFileName)

	// note that this will fail on windows, since there's no touch command
	settings := fileformat.ProcessSettings{
		Name: "touch",
		Wait: true,
		Args: []string{tmpFilePath},
	}

	event, err := startProcess(settings)
	require.NoError(t, err)

	// check if the file exists
	_, err = os.Stat(tmpFilePath)
	require.NoError(t, err)

	require.Equal(t, event["started_process_name"], settings.Name)
}
