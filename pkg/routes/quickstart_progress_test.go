package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedHatInsights/quickstarts/pkg/database"
	"github.com/RedHatInsights/quickstarts/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func mockQuickstartProgress(id uint) *models.QuickstartProgress {
	var quickstartProgress models.QuickstartProgress

	quickstartProgress.ID = id
	quickstartProgress.QuickstartID = 1234 + id
	quickstartProgress.AccountId = 4321

	database.DB.Create(&quickstartProgress)

	return &quickstartProgress
}

func setupQuickstartProgressRouter() *gin.Engine {
	r := gin.Default()
	r.Use(QuickstartEntityContext())
	r.GET("/", getAllQuickstartsProgress)
	r.POST("/:quickstartId", createQuickstartProgress)
	r.GET("/get", getQuickstartProgress)
	return r
}

func TestGetAllQuickstartProgresses(t *testing.T) {
	router := setupQuickstartProgressRouter()

	qp1 := mockQuickstartProgress(1)
	qp2 := mockQuickstartProgress(2)

	t.Run("returns GET all quickstarts successfully", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		type responsePayload struct {
			Data []models.QuickstartProgress
		}

		var payload *responsePayload
		json.NewDecoder(response.Body).Decode(&payload)
		assert.Equal(t, 200, response.Code)
		assert.Equal(t, 2, len(payload.Data))
		assert.Equal(t, qp1.AccountId, payload.Data[0].AccountId)
		assert.Equal(t, qp1.QuickstartID, payload.Data[0].QuickstartID)
		assert.Equal(t, qp2.AccountId, payload.Data[1].AccountId)
		assert.Equal(t, qp2.QuickstartID, payload.Data[1].QuickstartID)
	})
}
