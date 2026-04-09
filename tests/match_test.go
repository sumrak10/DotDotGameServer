package tests

import (
	"OnlineGame/database"
	"OnlineGame/manager"
	"OnlineGame/server/api/auth"
	wsAPI "OnlineGame/server/api/ws"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var testServer *httptest.Server
var testUserID1 uint
var testToken1 string
var testUserID2 uint
var testToken2 string

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	repo := database.NewUserRepository(database.GetDB())

	created1, _ := repo.Create("test1", "test1", "test1")
	testUserID1 = created1.ID
	testToken1, _ = auth.GenerateToken(created1.ID)
	created2, _ := repo.Create("test2", "test2", "test2")
	testUserID2 = created2.ID
	testToken2, _ = auth.GenerateToken(created2.ID)
}

func setupTestServer(t *testing.T) *httptest.Server {
	// Server
	ts := httptest.NewServer(http.HandlerFunc(wsAPI.HandleWS))
	t.Cleanup(func() { ts.Close() })

	return ts
}

func TestWebSocketConnect(t *testing.T) {
	ts := setupTestServer(t)
	// User1 connection
	conn1, _, err := websocket.DefaultDialer.Dial("ws"+ts.URL[4:]+"/ws?token="+testToken1, nil)
	assert.NoError(t, err)
	defer conn1.Close()
	// User2 connection
	conn2, _, err := websocket.DefaultDialer.Dial("ws"+ts.URL[4:]+"/ws?token="+testToken2, nil)
	assert.NoError(t, err)
	defer conn2.Close()
}

func TestWebSocketMessage(t *testing.T) {
	ts := setupTestServer(t)
	// User1 connection
	conn1, _, err := websocket.DefaultDialer.Dial("ws"+ts.URL[4:]+"/ws?token="+testToken1, nil)
	assert.NoError(t, err)
	defer conn1.Close()
	// User2 connection
	conn2, _, err := websocket.DefaultDialer.Dial("ws"+ts.URL[4:]+"/ws?token="+testToken2, nil)
	assert.NoError(t, err)
	defer func(conn2 *websocket.Conn) {
		err := conn2.Close()
		if err != nil {
			panic(err)
		}
	}(conn2)

	// User1 creating a match
	fmt.Println("User1 creating match...")
	_payload, err := json.Marshal(manager.CreateMatchEventMessage{
		WorldName: "default",
	})
	data, _ := manager.NewInputMessage(
		"match",
		manager.MatchEventMessage{
			Type:    "create",
			Payload: _payload,
		},
	)
	err = conn1.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	// Receive MatchID
	var responseCreateMatch map[string]interface{}
	err = conn1.ReadJSON(&responseCreateMatch)
	assert.NoError(t, err)

	var matchID uint
	if responseCreateMatch["status"] == "success" {
		matchID = uint(responseCreateMatch["payload"].(float64))
		fmt.Printf("Match created! ID: %d\n", matchID)
	} else {
		fmt.Println("error:", responseCreateMatch["error"])
		return
	}

	// User2 joining to match
	fmt.Println("User2 joining to match...")
	_payload, err = json.Marshal(manager.JoinMatchEventMessage{
		MatchID: matchID,
	})
	data, _ = manager.NewInputMessage(
		"match",
		manager.MatchEventMessage{
			Type:    "join",
			Payload: _payload,
		},
	)
	err = conn2.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	var responseJoinMatch map[string]interface{}
	err = conn2.ReadJSON(&responseJoinMatch)
	assert.NoError(t, err)
	if responseJoinMatch["status"] == "success" {
		fmt.Println("User2 joined match!")
	} else {
		fmt.Println("error:", responseJoinMatch["error"])
		return
	}

	// User1 starting a match
	fmt.Println("User1 starting match...")
	_payload, err = json.Marshal(manager.StartMatchEventMessage{
		MatchID: matchID,
	})
	data, _ = manager.NewInputMessage(
		"match",
		manager.MatchEventMessage{
			Type:    "start",
			Payload: _payload,
		},
	)
	err = conn1.WriteMessage(websocket.TextMessage, data)
	assert.NoError(t, err)

	var responseStartMatch map[string]interface{}
	err = conn1.ReadJSON(&responseStartMatch)
	assert.NoError(t, err)
	if responseStartMatch["status"] == "success" {
		fmt.Println("User1 started match!")
	} else {
		fmt.Println("error:", responseStartMatch["error"])
		return
	}

	var responseMatchTick map[string]interface{}
	for {
		err = conn1.ReadJSON(&responseMatchTick)
		assert.NoError(t, err)
		fmt.Println(responseMatchTick)
	}
}
