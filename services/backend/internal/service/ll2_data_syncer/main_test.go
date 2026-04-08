package ll2datasyncer

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/service/core"
	"github.com/vamosdalian/launchdate-backend/internal/service/ll2"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

var (
	mongoContainer *mongodb.MongoDBContainer
	testDB         *db.MongoDB
	mockServer     *httptest.Server
	ll2Service     *ll2.LL2Service
	coreService    *core.MainService

	endpointFileMap = map[string]string{
		"agencies":                        "agency.json",
		"launcher_configurations":         "launcher.json",
		"launcher_configuration_families": "launcher_family.json",
		"locations":                       "location.json",
		"pads":                            "pad.json",
		"launches":                        "launch.json",
	}
)

func TestMain(m *testing.M) {
	// Silence logrus during tests
	logrus.SetOutput(io.Discard)

	ctx := context.Background()

	// Start Global Mock Server
	mockServer = setupMockLL2API()

	// Start MongoDB container
	var err error
	mongoContainer, err = mongodb.Run(ctx, "mongo:6")
	if err != nil {
		panic(err)
	}

	// Get connection string
	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	// Initialize DB connection
	testDB, _, err = db.NewMongoDB(uri, "test_db_syncer")
	if err != nil {
		panic(err)
	}

	// Initialize Services
	httpCli := util.NewHTTPClient()
	cfg := &config.Config{
		LL2URLPrefix: mockServer.URL,
	}
	ll2Service = ll2.NewLL2Service(cfg, testDB, httpCli)
	coreService = core.NewMainService(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	if err := mongoContainer.Terminate(ctx); err != nil {
		panic(err)
	}
	mockServer.Close()

	os.Exit(code)
}

func setupMockLL2API() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		// URL: /2.3.0/endpoint...
		// parts: ["", "2.3.0", "endpoint"]
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		endpoint := parts[2]

		// Remove query params if any
		if strings.Contains(endpoint, "?") {
			endpoint = strings.Split(endpoint, "?")[0]
		}

		if filename, ok := endpointFileMap[endpoint]; ok {
			// Path relative to internal/service/ll2_data_syncer
			path := filepath.Join("..", "ll2", "testdata", filename)
			f, err := os.Open(path)
			if err != nil {
				http.Error(w, "Failed to read test data: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			w.Header().Set("Content-Type", "application/json")
			_, err = io.Copy(w, f)
			if err != nil {
				// Just ignore error in test mock
			}
			return
		}

		http.NotFound(w, r)
	}))
}

// Helper to clear collections
func clearCollections(t *testing.T, collections ...string) {
	for _, coll := range collections {
		err := testDB.Database.Collection(coll).Drop(context.Background())
		if err != nil {
			t.Logf("Failed to drop collection %s: %v", coll, err)
		}
	}
}
