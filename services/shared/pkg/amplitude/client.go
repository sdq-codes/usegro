package amplitude

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

const httpAPIURL = "https://api.eu.amplitude.com/2/httpapi"

// Event represents a single Amplitude event payload.
type Event struct {
	UserID          string                 `json:"user_id"`
	EventType       string                 `json:"event_type"`
	EventProperties map[string]interface{} `json:"event_properties,omitempty"`
	Time            int64                  `json:"time"`
}

type batchPayload struct {
	APIKey string  `json:"api_key"`
	Events []Event `json:"events"`
}

// Client is the Amplitude HTTP client.
type Client struct {
	apiKey      string
	httpClient  *http.Client
	environment string
	appVersion  string
}

var (
	defaultClient *Client
	once          sync.Once
)

// Init initialises the singleton Amplitude client from environment variables.
// Safe to call multiple times — only the first call has effect.
func Init() {
	once.Do(func() {
		defaultClient = &Client{
			apiKey:      os.Getenv("AMPLITUDE_API_KEY"),
			httpClient:  &http.Client{Timeout: 10 * time.Second},
			environment: envOrDefault("ENV", "development"),
			appVersion:  envOrDefault("APP_VERSION", "1.0.0"),
		}
	})
}

// Track fires a single event asynchronously. No-ops when the client is not
// initialised or AMPLITUDE_API_KEY is unset, so it never blocks or panics.
func Track(userID, eventType string, properties map[string]interface{}) {
	if defaultClient == nil || defaultClient.apiKey == "" {
		return
	}
	if properties == nil {
		properties = make(map[string]interface{})
	}
	properties["environment"] = defaultClient.environment
	properties["app_version"] = defaultClient.appVersion

	go defaultClient.sendWithRetry([]Event{{
		UserID:          userID,
		EventType:       eventType,
		EventProperties: properties,
		Time:            time.Now().UnixMilli(),
	}})
}

// sendWithRetry attempts up to 3 POSTs with exponential back-off.
// Stops retrying on any non-5xx response.
func (c *Client) sendWithRetry(events []Event) {
	body, err := json.Marshal(batchPayload{APIKey: c.apiKey, Events: events})
	if err != nil {
		return
	}

	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(math.Pow(2, float64(attempt))) * time.Second)
		}
		resp, err := c.httpClient.Post(httpAPIURL, "application/json", bytes.NewReader(body))
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 500 {
			return
		}
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
