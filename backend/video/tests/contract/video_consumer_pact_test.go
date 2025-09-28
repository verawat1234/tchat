package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

// TestVideoConsumerPact runs consumer tests against the video service
func TestVideoConsumerPact(t *testing.T) {
	// Create Pact consumer
	mockProvider, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "video-web-client",
		Provider: "video-service",
		Host:     "localhost",
		Port:     9000,
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	t.Run("GetVideos", func(t *testing.T) {
		// Create interaction
		err := mockProvider.
			AddInteraction().
			Given("videos exist").
			UponReceiving("a request for videos").
			WithRequest("GET", "/api/v1/videos?page=1&per_page=10").
			WillRespondWith(consumer.Response{
				Status: 200,
				Headers: matchers.MapMatcher{
					"Content-Type": matchers.String("application/json; charset=utf-8"),
				},
				Body: matchers.Map{
					"success": matchers.Bool(true),
					"data": matchers.Map{
						"videos": matchers.EachLike(matchers.Map{
							"id":          matchers.String("550e8400-e29b-41d4-a716-446655440000"),
							"title":       matchers.String("Sample Video"),
							"description": matchers.String("Sample video description"),
							"thumbnail":   matchers.String("https://example.com/thumbnail.jpg"),
							"videoUrl":    matchers.String("https://example.com/video.mp4"),
							"duration":    matchers.String("1:30"),
							"views":       matchers.Like(1000),
							"likes":       matchers.Like(50),
							"category":    matchers.String("entertainment"),
							"tags":        matchers.EachLike("tag"),
							"type":        matchers.String("short"),
							"status":      matchers.String("active"),
							"channelId":   matchers.String("550e8400-e29b-41d4-a716-446655440001"),
							"channel": matchers.Map{
								"id":          matchers.String("550e8400-e29b-41d4-a716-446655440001"),
								"name":        matchers.String("Sample Channel"),
								"avatar":      matchers.String("https://example.com/avatar.jpg"),
								"subscribers": matchers.Like(1000),
								"verified":    matchers.Like(true),
							},
							"createdAt": matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
							"updatedAt": matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
						}, 1),
						"page":     matchers.Like(1),
						"per_page": matchers.Like(10),
						"total":    matchers.Like(1),
						"has_more": matchers.Like(false),
					},
				},
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				// Make HTTP request to mock server
				url := fmt.Sprintf("http://%s:%d/api/v1/videos?page=1&per_page=10", config.Host, config.Port)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode)
				return nil
			})

		assert.NoError(t, err)
	})

	t.Run("GetVideo", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("video exists").
			UponReceiving("a request for a specific video").
			WithRequest("GET", "/api/v1/videos/550e8400-e29b-41d4-a716-446655440000").
			WithHeaders(matchers.MapMatcher{
				"Content-Type": matchers.String("application/json"),
			}).
			WillRespondWith(consumer.Response{
				Status: 200,
				Headers: matchers.MapMatcher{
					"Content-Type": matchers.String("application/json; charset=utf-8"),
				},
				Body: matchers.Map{
					"success": matchers.Bool(true),
					"data": matchers.Map{
						"id":          matchers.String("550e8400-e29b-41d4-a716-446655440000"),
						"title":       matchers.String("Sample Video"),
						"description": matchers.String("Sample video description"),
						"thumbnail":   matchers.String("https://example.com/thumbnail.jpg"),
						"videoUrl":    matchers.String("https://example.com/video.mp4"),
						"duration":    matchers.String("1:30"),
						"views":       matchers.Like(1000),
						"likes":       matchers.Like(50),
						"category":    matchers.String("entertainment"),
						"tags":        matchers.EachLike("tag"),
						"type":        matchers.String("short"),
						"status":      matchers.String("active"),
						"channelId":   matchers.String("550e8400-e29b-41d4-a716-446655440001"),
						"channel": matchers.Map{
							"id":          matchers.String("550e8400-e29b-41d4-a716-446655440001"),
							"name":        matchers.String("Sample Channel"),
							"avatar":      matchers.String("https://example.com/avatar.jpg"),
							"subscribers": matchers.Like(1000),
							"verified":    matchers.Like(true),
						},
						"createdAt": matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
						"updatedAt": matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
					},
				},
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/videos/550e8400-e29b-41d4-a716-446655440000", config.Host, config.Port)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode)
				return nil
			})

		assert.NoError(t, err)
	})

	t.Run("GetVideoNotFound", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("video does not exist").
			UponReceiving("a request for non-existent video").
			WithRequest("GET", "/api/v1/videos/non-existent-id").
			WithHeaders(matchers.MapMatcher{
				"Content-Type": matchers.String("application/json"),
			}).
			WillRespondWith(consumer.Response{
				Status: 404,
				Headers: matchers.MapMatcher{
					"Content-Type": matchers.String("application/json; charset=utf-8"),
				},
				Body: matchers.Map{
					"success": matchers.Bool(false),
					"error":   matchers.String("Video not found"),
				},
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/videos/non-existent-id", config.Host, config.Port)
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, 404, resp.StatusCode)
				return nil
			})

		assert.NoError(t, err)
	})

	t.Run("CreateVideo", func(t *testing.T) {
		err := mockProvider.
			AddInteraction().
			Given("valid video data").
			UponReceiving("a request to create video").
			WithRequest("POST", "/api/v1/videos").
			WithHeaders(matchers.MapMatcher{
				"Content-Type": matchers.String("application/json"),
			}).
			WithJSONBody(matchers.Map{
				"title":       matchers.String("New Video"),
				"description": matchers.String("New video description"),
				"videoUrl":    matchers.String("https://example.com/new-video.mp4"),
				"category":    matchers.String("entertainment"),
				"type":        matchers.String("short"),
				"channelId":   matchers.String("550e8400-e29b-41d4-a716-446655440001"),
			}).
			WillRespondWith(consumer.Response{
				Status: 200,
				Headers: matchers.MapMatcher{
					"Content-Type": matchers.String("application/json; charset=utf-8"),
				},
				Body: matchers.Map{
					"success": matchers.Bool(true),
					"message": matchers.String("Video created successfully"),
					"data": matchers.Map{
						"id":          matchers.String("550e8400-e29b-41d4-a716-446655440002"),
						"title":       matchers.String("New Video"),
						"description": matchers.String("New video description"),
						"videoUrl":    matchers.String("https://example.com/new-video.mp4"),
						"category":    matchers.String("entertainment"),
						"type":        matchers.String("short"),
						"status":      matchers.String("active"),
						"channelId":   matchers.String("550e8400-e29b-41d4-a716-446655440001"),
						"views":       matchers.Like(0),
						"likes":       matchers.Like(0),
						"createdAt":   matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
						"updatedAt":   matchers.Regex("2000-02-01T12:30:00Z", `^[\+\-]?\d{4}(?!\d{2}\b)((-?)((0[1-9]|1[0-2])(\3([12]\d|0[1-9]|3[01]))?|W([0-4]\d|5[0-2])(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+(?!:))?)?(\17[0-5]\d([\.,]\d+)?)?([zZ]|([\+\-])([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`),
					},
				},
			}).
			ExecuteTest(t, func(config consumer.MockServerConfig) error {
				url := fmt.Sprintf("http://%s:%d/api/v1/videos", config.Host, config.Port)

				videoData := map[string]interface{}{
					"title":       "New Video",
					"description": "New video description",
					"videoUrl":    "https://example.com/new-video.mp4",
					"category":    "entertainment",
					"type":        "short",
					"channelId":   "550e8400-e29b-41d4-a716-446655440001",
				}

				jsonData, err := json.Marshal(videoData)
				if err != nil {
					return err
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
				if err != nil {
					return err
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				assert.Equal(t, 200, resp.StatusCode)
				return nil
			})

		assert.NoError(t, err)
	})
}