package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// T035: Integration test multi-region data sync
// Tests data synchronization across Southeast Asian regions
type RegionSyncTestSuite struct {
	suite.Suite
	router  *gin.Engine
	regions map[string]map[string]interface{} // region -> data
}

func TestRegionSyncSuite(t *testing.T) {
	suite.Run(t, new(RegionSyncTestSuite))
}

func (suite *RegionSyncTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()
	suite.regions = map[string]map[string]interface{}{
		"ap-southeast-1": make(map[string]interface{}), // Singapore
		"ap-southeast-2": make(map[string]interface{}), // Thailand
		"ap-southeast-3": make(map[string]interface{}), // Indonesia
	}

	suite.setupRegionSyncEndpoints()
}

func (suite *RegionSyncTestSuite) setupRegionSyncEndpoints() {
	// Sync data across regions
	suite.router.POST("/sync/regions", func(c *gin.Context) {
		var req map[string]interface{}
		c.ShouldBindJSON(&req)

		sourceRegion := req["source_region"].(string)
		targetRegions := req["target_regions"].([]interface{})
		data := req["data"].(map[string]interface{})

		// Store in source region
		suite.regions[sourceRegion] = data

		// Replicate to target regions
		for _, region := range targetRegions {
			regionStr := region.(string)
			suite.regions[regionStr] = data
		}

		c.JSON(http.StatusOK, gin.H{
			"synced_regions": len(targetRegions) + 1,
			"sync_time":      time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Get regional data
	suite.router.GET("/regions/:region/data/:key", func(c *gin.Context) {
		region := c.Param("region")
		key := c.Param("key")

		regionData, exists := suite.regions[region]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "region_not_found"})
			return
		}

		data, exists := regionData[key]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "data_not_found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"region": region,
			"key":    key,
			"data":   data,
		})
	})
}

func (suite *RegionSyncTestSuite) TestMultiRegionSync() {
	// Test data sync across Southeast Asian regions
	syncData := map[string]interface{}{
		"source_region":  "ap-southeast-1",
		"target_regions": []string{"ap-southeast-2", "ap-southeast-3"},
		"data": map[string]interface{}{
			"user_123": map[string]interface{}{
				"phone":   "+65 9123 4567",
				"country": "SG",
				"balance": 1000.0,
			},
		},
	}

	jsonData, _ := json.Marshal(syncData)
	req := httptest.NewRequest("POST", "/sync/regions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	start := time.Now()
	suite.router.ServeHTTP(w, req)
	syncDuration := time.Since(start)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.True(suite.T(), syncDuration < 200*time.Millisecond, "Region sync should complete in <200ms")

	// Verify data exists in all regions
	regions := []string{"ap-southeast-1", "ap-southeast-2", "ap-southeast-3"}
	for _, region := range regions {
		req = httptest.NewRequest("GET", "/regions/"+region+"/data/user_123", nil)
		w = httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(suite.T(), region, response["region"])

		userData := response["data"].(map[string]interface{})
		assert.Equal(suite.T(), "SG", userData["country"])
	}
}

func (suite *RegionSyncTestSuite) TestRegionLatency() {
	// Test acceptable latency for Southeast Asian regions
	regions := []string{"ap-southeast-1", "ap-southeast-2", "ap-southeast-3"}

	for _, region := range regions {
		suite.T().Logf("Testing latency for region: %s", region)

		req := httptest.NewRequest("GET", "/regions/"+region+"/data/test", nil)
		w := httptest.NewRecorder()

		start := time.Now()
		suite.router.ServeHTTP(w, req)
		latency := time.Since(start)

		// Should respond quickly even for missing data
		assert.True(suite.T(), latency < 100*time.Millisecond, "Region %s latency should be <100ms", region)
	}
}