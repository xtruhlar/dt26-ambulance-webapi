package ambulance_ufe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xtruhlar/dt26-ambulance-webapi/internal/db_service"
)

type ConsultationSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Ambulance]
}

func TestConsultationSuite(t *testing.T) {
	suite.Run(t, new(ConsultationSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *ConsultationSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Ambulance]{}

	// Compile-time check that mock satisfies the interface
	var _ db_service.DbService[Ambulance] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Ambulance{
				Id: "test-ambulance",
				ConsultationEntries: []ConsultationEntryFull{
					{
						ConsultationEntry: ConsultationEntry{
							Id:          "test-entry",
							PatientId:   "test-patient",
							PatientName: "Test Patient",
							Condition:   "Hypertenzia",
							Status:      "active",
							CreatedAt:   time.Now(),
						},
					},
				},
			},
			nil,
		)
}

func (suite *ConsultationSuite) Test_UpdateConsultation_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	body := `{"status":"closed","condition":"Hypertenzia"}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("PUT",
		"/api/remote-consultation/test-ambulance/entries/test-entry",
		strings.NewReader(body))

	sut := implAmbulanceRemoteConsultationAPI{}

	// ACT
	sut.UpdateConsultationEntry(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
}

func (suite *ConsultationSuite) Test_GetConsultationEntries_ReturnsEntries() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
	}
	ctx.Request = httptest.NewRequest("GET",
		"/api/remote-consultation/test-ambulance/entries", nil)

	sut := implAmbulanceRemoteConsultationAPI{}

	// ACT
	sut.GetConsultationEntries(ctx)

	// ASSERT
	suite.Equal(http.StatusOK, recorder.Code)
	suite.dbServiceMock.AssertCalled(suite.T(), "FindDocument", mock.Anything, "test-ambulance")
}

func (suite *ConsultationSuite) Test_DeleteConsultationEntry_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "ambulanceId", Value: "test-ambulance"},
		{Key: "entryId", Value: "test-entry"},
	}
	ctx.Request = httptest.NewRequest("DELETE",
		"/api/remote-consultation/test-ambulance/entries/test-entry", nil)

	sut := implAmbulanceRemoteConsultationAPI{}

	// ACT
	sut.DeleteConsultationEntry(ctx)

	// ASSERT
	suite.Equal(http.StatusNoContent, recorder.Code)
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "test-ambulance", mock.Anything)
}
