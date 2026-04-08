package ll2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

var (
	mongoContainer *mongodb.MongoDBContainer
	testDB         *db.MongoDB
	mockServer     *httptest.Server

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
	testDB, _, err = db.NewMongoDB(uri, "test_db")
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := mongoContainer.Terminate(ctx); err != nil {
		panic(err)
	}
	mockServer.Close()

	os.Exit(code)
}

func newTestLL2Service() *LL2Service {
	return &LL2Service{
		mongoClient:  testDB,
		ll2UrlPrefix: mockServer.URL,
		httpClient:   util.NewHTTPClient(), // Ensure http client is initialized
	}
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

func setupMockLL2API() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		endpoint := parts[2]
		var targetID string
		if len(parts) > 3 && parts[3] != "" {
			targetID = parts[3]
		}

		if filename, ok := endpointFileMap[endpoint]; ok {
			path := filepath.Join("testdata", filename)
			f, err := os.Open(path)
			if err != nil {
				http.Error(w, "Failed to read test data: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer f.Close()

			w.Header().Set("Content-Type", "application/json")

			if targetID == "" {
				_, err = io.Copy(w, f)
				if err != nil {
					// Just ignore error in test mock
				}
				return
			}

			// Handle single item retrieval
			var listResp struct {
				Results []map[string]interface{} `json:"results"`
			}
			if err := json.NewDecoder(f).Decode(&listResp); err != nil {
				http.Error(w, "Failed to decode test data", http.StatusInternalServerError)
				return
			}

			for _, item := range listResp.Results {
				idVal, ok := item["id"]
				if !ok {
					continue
				}
				if fmt.Sprintf("%v", idVal) == targetID {
					json.NewEncoder(w).Encode(item)
					return
				}
			}

			http.NotFound(w, r)
			return
		}

		http.NotFound(w, r)
	}))
}
