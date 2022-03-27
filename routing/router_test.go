package routing

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// add handler dependencies if needed

	router := gin.Default()

	// add routes and handlers as needed
	router.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, mockUser{Users: "mock response"})
	})

	return router
}

type routeTest struct {
	httpMethod   string
	path         string
	body         io.Reader
	wantErr      bool
	wantRespCode int
	want         any
	got          any
}

func doRouteTest(t *testing.T, options routeTest) {
	testRouter := setupTestRouter()

	req, err := http.NewRequest(options.httpMethod, options.path, options.body)
	if err != nil {
		t.Fatalf("couldn't create request: %v", err)
	}

	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	if options.wantRespCode != resp.Code {
		t.Errorf("expected response code %v got %v", options.wantRespCode, resp.Code)
	}

	respBytes, err := io.ReadAll(resp.Result().Body)
	if err != nil {
		t.Errorf("error reading response body bytes: %v", err)
	}

	// unmarshal and check result
	if options.wantErr {
		// unmarshal into custom error type, if any, and check against expected properties
		return
	}

	err = json.Unmarshal(respBytes, &options.got)
	if err != nil {
		t.Errorf("error unmarshaling response: %v", err)
	}

	if respDiff := cmp.Diff(options.want, options.got); respDiff != "" {
		t.Errorf("path %v response mismatch (-want +got)\n%v", options.path, respDiff)
	}
}

type mockUser struct {
	Users string `json:"users,omitempty"`
}

func TestRoute(t *testing.T) {
	rt := routeTest{
		httpMethod:   http.MethodGet,
		path:         "/users",
		body:         nil,
		wantRespCode: http.StatusOK,
		// use pointers here or go can't unmarshal to specific struct
		want: &mockUser{Users: "mock response"},
		got:  &mockUser{},
	}

	doRouteTest(t, rt)
}
