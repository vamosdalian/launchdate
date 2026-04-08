package ll2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vamosdalian/launchdate-backend/internal/config"
	"github.com/vamosdalian/launchdate-backend/internal/db"
	"github.com/vamosdalian/launchdate-backend/internal/util"
)

type LL2Service struct {
	mongoClient  *db.MongoDB
	httpClient   *util.HTTPClient
	ll2UrlPrefix string
}

func NewLL2Service(conf *config.Config, db *db.MongoDB, hc *util.HTTPClient) *LL2Service {
	return &LL2Service{
		mongoClient:  db,
		ll2UrlPrefix: conf.LL2URLPrefix,
		httpClient:   hc,
	}
}

func (s *LL2Service) GetDataFromAPI(endpoint string, limit, offset int, payload any, params ...string) error {
	var url strings.Builder
	fmt.Fprintf(&url, "%s/2.3.0/%s?limit=%d&offset=%d&mode=detailed", s.ll2UrlPrefix, endpoint, limit, offset)
	for _, param := range params {
		url.WriteString("&" + param)
	}

	resp, err := s.httpClient.Get(context.Background(), url.String(), nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal LL2 API response: %w", err)
	}

	return nil
}
