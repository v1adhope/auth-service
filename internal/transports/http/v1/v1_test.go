package httpv1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/v1adhope/auth-service/internal/services"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/alert"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/hash"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/repositories"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/tokens"
	"github.com/v1adhope/auth-service/internal/services/infrastructure/validator"
	"github.com/v1adhope/auth-service/internal/testhelpers"
	httpv1 "github.com/v1adhope/auth-service/internal/transports/http/v1"
	"github.com/v1adhope/auth-service/pkg/logger"
	"github.com/v1adhope/auth-service/pkg/postgresql"
)

const (
	_pgMigrationsSourceUrl = "file://../../../../db/migrations"

	_tokensAccessKey  = "secret"
	_tokensAccessTtl  = 3 * time.Minute
	_tokensRefreshkey = "HC2fAkS4Lyfisrt4agCZgRU7eWPpFgbH"
	_tokensIssuer     = "auth-service-test"

	_loggerLevel = "debug"

	_handlerMode = gin.DebugMode
)

var (
	_handlerAllowOrigins = []string{"*"}
	_handlerAllowMethods = []string{"POST", "HEAD", "OPTIONS"}
	_handlerAllowHeaders = []string{"Origin"}
)

type Suite struct {
	suite.Suite
	pgC       *testhelpers.PostgresContainer
	handlerV1 *gin.Engine
	ctx       context.Context
}

func (s *Suite) SetupSuite() {
	t := s.T()

	s.ctx = context.Background()

	pgC, err := testhelpers.BuildContainer(s.ctx, _pgMigrationsSourceUrl)
	if err != nil {
		log.Fatal(err)
	}

	s.pgC = pgC

	driver, err := postgresql.Build(
		s.ctx,
		postgresql.WithConnStr(pgC.ConnStr),
	)
	if err != nil {
		log.Fatal(err)
	}
	t.Cleanup(func() {
		driver.Close()
	})

	if err := s.pgC.MigrateUp(); err != nil {
		log.Fatal(err)
	}

	repos := repositories.New(driver)

	alert := alert.New()

	validator := validator.New()

	hash := hash.New()

	tokenManager := tokens.New(
		tokens.WithAccessKey(_tokensAccessKey),
		tokens.WithAccessTtl(_tokensAccessTtl),
		tokens.WithRefreshKey(_tokensRefreshkey),
		tokens.WithIssuer(_tokensIssuer),
	)

	services := services.New(
		validator,
		tokenManager,
		hash,
		repos,
		alert,
	)

	log := logger.New(
		logger.WithLevel(_loggerLevel),
	)

	handler := httpv1.New(services, log).Handler(
		httpv1.WithAllowOrigins(_handlerAllowOrigins),
		httpv1.WithAllowMethods(_handlerAllowMethods),
		httpv1.WithAllowHeaders(_handlerAllowHeaders),
		httpv1.WithMode(_handlerMode),
	)

	s.handlerV1 = handler
}

func (s *Suite) TearDownSuite() {
	if err := s.pgC.Terminate(s.ctx); err != nil {
		log.Fatalf("httpv1_test: v1_test: TearDownSuite: Terminate: %v", err)
	}
}

type testAuthResp struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func (s *Suite) TestAuth() {
	t := s.T()
	tcs := []struct {
		key   string
		input string
	}{
		{
			key:   "Case 1",
			input: "adb21fec-7892-416a-bbfc-9b2d77e8db4a",
		},
		{
			key:   "Case 2",
			input: "01f20929-dc51-4edb-a472-5672f4678fa2",
		},
		{
			key:   "Case 3",
			input: "4512d372-9de4-4ef3-b528-e4950006660d",
		},
	}
	resp := testAuthResp{}

	for _, tc := range tcs {
		t.Run("get", func(t *testing.T) {
			// INFO: get
			sut := httptest.NewRecorder()
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/tokens/%s", tc.input),
				nil,
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, sut.Code, tc.key)

			err = json.Unmarshal(sut.Body.Bytes(), &resp)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp.Access, tc.key)
			assert.NotEmpty(t, resp.Refresh, tc.key)

			// INFO: refresh
			sut = httptest.NewRecorder()
			jsonData, err := json.Marshal(resp)

			assert.NoError(t, err, tc.key)

			req, err = http.NewRequest(
				"POST",
				"/v1/tokens/refresh",
				strings.NewReader(string(jsonData)),
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, sut.Code, tc.key)

			err = json.Unmarshal(sut.Body.Bytes(), &resp)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp.Access, tc.key)
			assert.NotEmpty(t, resp.Refresh, tc.key)
		})
	}
}

func (s *Suite) TestGetTokenPairNegative() {
	t := s.T()
	tcs := []struct {
		key   string
		input string
	}{
		{
			key:   "Incorrect uuid",
			input: "adb21fec-416a-bbfc-9b2d77e8db4a",
		},
		{
			key:   "Incorrect uuid type",
			input: "11",
		},
	}

	for _, tc := range tcs {
		t.Run("get", func(t *testing.T) {
			sut := httptest.NewRecorder()
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/tokens/%s", tc.input),
				nil,
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusBadRequest, sut.Code, tc.key)
		})
	}
}

func (s *Suite) TestRefreshTokenPairNegative() {
	t := s.T()
	tcs := []struct {
		key   string
		input testAuthResp
	}{
		{
			key: "Incorrect data",
			input: testAuthResp{
				"some access",
				"some refresh",
			},
		},
		{
			key: "Expired or not valid",
			input: testAuthResp{
				"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpcCI6IjE5Mi4xNjguNjUuMSIsImlzcyI6ImF1dGgtc2VydmljZSIsInN1YiI6IjUwNGJhY2Q3LTU3MWMtNGY0OS05ZTBlLTg3ZTg4M2Y5NjU4YSIsImV4cCI6MTcyNDQzMjI0OCwiaWF0IjoxNzI0NDMxMDQ4LCJqdGkiOiIwMWVmNjE2ZC1mYzcxLTYwODItOWFlYy0wMjQyYWMxMjAwMDMifQ.q3tidVjk6SdpIauJ_HBFBCf8QnRHhb6oNfY_qsMwNJjVEeiuC9IEFakDP2w1Rx4NEy9-KXfIzh3eXRFB9hp3_w",
				"LuITcRAwHGcsWsOt9wbzulpxzikFZHMvtYUKzaNGN+9m6qdQxk4F/XcI7eXylBR9QcpIeX6E+6XE66UzfFYoWQ==",
			},
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			sut := httptest.NewRecorder()
			jsonData, err := json.Marshal(tc.input)

			assert.NoError(t, err, tc.key)

			req, err := http.NewRequest(
				"POST",
				"/v1/tokens/refresh",
				strings.NewReader(string(jsonData)),
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusBadRequest, sut.Code, tc.key)
		})
	}
}

func (s *Suite) TestRefreshTokenCanNotBeUsingTwice() {
	t := s.T()
	tcs := []struct {
		key   string
		input string
	}{
		{
			key:   "Case 1",
			input: "adb21fec-7892-416a-bbfc-9b2d77e8db4a",
		},
		{
			key:   "Case 2",
			input: "01f20929-dc51-4edb-a472-5672f4678fa2",
		},
		{
			key:   "Case 3",
			input: "4512d372-9de4-4ef3-b528-e4950006660d",
		},
	}
	resp := testAuthResp{}
	jsonData := []byte{}

	for _, tc := range tcs {
		t.Run("get", func(t *testing.T) {
			// INFO: get
			w := httptest.NewRecorder()
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/tokens/%s", tc.input),
				nil,
			)
			s.handlerV1.ServeHTTP(w, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, w.Code, tc.key)

			err = json.Unmarshal(w.Body.Bytes(), &resp)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp.Access, tc.key)
			assert.NotEmpty(t, resp.Refresh, tc.key)

			// INFO: first refresh
			w = httptest.NewRecorder()
			jsonData, err = json.Marshal(resp)

			assert.NoError(t, err, tc.key)

			req, err = http.NewRequest(
				"POST",
				"/v1/tokens/refresh",
				strings.NewReader(string(jsonData)),
			)
			s.handlerV1.ServeHTTP(w, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, w.Code, tc.key)

			err = json.Unmarshal(w.Body.Bytes(), &resp)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp.Access, tc.key)
			assert.NotEmpty(t, resp.Refresh, tc.key)

			// INFO: sut
			sut := httptest.NewRecorder()
			req, err = http.NewRequest(
				"POST",
				"/v1/tokens/refresh",
				strings.NewReader(string(jsonData)),
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusBadRequest, sut.Code, tc.key)
		})
	}
}

type inputDuoID struct {
	firstId, secondId string
}

func (s *Suite) TestTokensFromDifferentPairs() {
	t := s.T()
	tcs := []struct {
		key   string
		input inputDuoID
	}{
		{
			key: "Case 1",
			input: inputDuoID{
				"adb21fec-7892-416a-bbfc-9b2d77e8db4a",
				"01f20929-dc51-4edb-a472-5672f4678fa2",
			},
		},
		{
			key: "Case 2",
			input: inputDuoID{
				"4512d372-9de4-4ef3-b528-e4950006660d",
				"27588644-cbae-47ed-b433-d44fa040c133",
			},
		},
		{
			key: "Case 3",
			input: inputDuoID{
				"dfe440aa-2ab3-4ccb-a395-8b9eec475d4f",
				"17197b42-9777-41b7-b2c0-28378fd52823",
			},
		},
	}

	for _, tc := range tcs {
		t.Run("get", func(t *testing.T) {
			// INFO: get first
			w := httptest.NewRecorder()
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/tokens/%s", tc.input.firstId),
				nil,
			)
			s.handlerV1.ServeHTTP(w, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, w.Code, tc.key)

			resp1 := testAuthResp{}
			err = json.Unmarshal(w.Body.Bytes(), &resp1)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp1.Access, tc.key)
			assert.NotEmpty(t, resp1.Refresh, tc.key)

			// INFO: get second
			w = httptest.NewRecorder()
			req, err = http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/tokens/%s", tc.input.secondId),
				nil,
			)
			s.handlerV1.ServeHTTP(w, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusCreated, w.Code, tc.key)

			resp2 := testAuthResp{}
			err = json.Unmarshal(w.Body.Bytes(), &resp2)

			assert.NoError(t, err, tc.key)
			assert.NotEmpty(t, resp2.Access, tc.key)
			assert.NotEmpty(t, resp2.Refresh, tc.key)

			// INFO: sut
			sut := httptest.NewRecorder()
			jsonData, err := json.Marshal(testAuthResp{
				resp1.Access,
				resp2.Refresh,
			})

			assert.NoError(t, err, tc.key)

			req, err = http.NewRequest(
				"POST",
				"/v1/tokens/refresh",
				strings.NewReader(string(jsonData)),
			)
			s.handlerV1.ServeHTTP(sut, req)

			assert.NoError(t, err, tc.key)
			assert.Equal(t, http.StatusBadRequest, sut.Code, tc.key)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
