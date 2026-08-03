package web_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/erd"
	"dbmw/core/explorer"
	"dbmw/core/project"
	"dbmw/core/query"
	"dbmw/storage"
	"dbmw/web"
)

func TestWebServer(t *testing.T) {
	tempDir := t.TempDir()
	cfgStore, _ := storage.NewConfigStore(tempDir + "/config.json")
	connStore, _ := storage.NewConnectionStore(tempDir + "/connections.json")
	histStore, _ := storage.NewHistoryStore(tempDir + "/history.db")
	defer histStore.Close()

	connSvc := connection.NewService(connStore, connector.DefaultFactory)
	defer connSvc.CloseAll()
	expSvc := explorer.NewService()
	qSvc := query.NewService(histStore)
	dataSvc := data.NewService()
	erdSvc := erd.NewService()
	projSvc := project.NewService()

	srv, err := web.NewServer(8085, connSvc, expSvc, qSvc, dataSvc, erdSvc, projSvc, cfgStore)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	app := srv.App()

	t.Run("GET /api/health", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/health", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/connections/test with valid sqlite memory", func(t *testing.T) {
		body, _ := json.Marshal(connection.ConnectionConfig{
			Name:     "Test SQLite Mem",
			Driver:   connection.DriverSQLite,
			FilePath: ":memory:",
		})
		req := httptest.NewRequest("POST", "/api/connections/test", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/connections and activate", func(t *testing.T) {
		body, _ := json.Marshal(connection.ConnectionConfig{
			Name:     "Local SQLite Test",
			Driver:   connection.DriverSQLite,
			FilePath: tempDir + "/test_web.db",
		})
		req := httptest.NewRequest("POST", "/api/connections", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		resBody, _ := io.ReadAll(resp.Body)
		var saved connection.ConnectionConfig
		json.Unmarshal(resBody, &saved)

		// Activate connection
		actBody, _ := json.Marshal(map[string]string{"id": saved.ID})
		actReq := httptest.NewRequest("POST", "/api/connections/active", bytes.NewReader(actBody))
		actReq.Header.Set("Content-Type", "application/json")
		actResp, err := app.Test(actReq)
		if err != nil || actResp.StatusCode != http.StatusOK {
			t.Fatalf("failed to activate: %v, status: %d", err, actResp.StatusCode)
		}

		// Run query on this connection
		qBody, _ := json.Marshal(map[string]string{
			"query": "CREATE TABLE demo (id INTEGER PRIMARY KEY, msg TEXT); INSERT INTO demo (msg) VALUES ('Hello Fiber');",
		})
		qReq := httptest.NewRequest("POST", "/api/query/execute", bytes.NewReader(qBody))
		qReq.Header.Set("Content-Type", "application/json")
		qResp, err := app.Test(qReq)
		if err != nil || qResp.StatusCode != http.StatusOK {
			t.Fatalf("failed query exec: %v", err)
		}

		// Browse data
		bReq := httptest.NewRequest("GET", "/api/data/browse/demo", nil)
		bResp, err := app.Test(bReq)
		if err != nil || bResp.StatusCode != http.StatusOK {
			t.Fatalf("failed to browse: %v", err)
		}
	})
}
