package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBoring(t *testing.T) {
	jwt1, err := MakeJWT(uuid.New(), "secretSauce", time.Minute)
	if err != nil {
		t.Errorf("making boring jwt failed %v", err)
		return
	}
	t.Logf("boring jwt is %s \n", jwt1)
	_, err = ValidateJWT(jwt1, "secretSauce")
	if err != nil {
		t.Errorf("validating boring jwt failed %v", err)
		return
	}
}

func TestStale(t *testing.T) {
	staleJwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI5ZmJjZGVkOC0xOGIyLTRkODktYTA3NS1iMjk3N2JiNjRjOGIiLCJleHAiOjE3NTQ4MDEzMDMsImlhdCI6MTc1NDgwMTI0M30.g3BotPB4j46_jJR9iSdircAb3mNeoXGQuVaZqTmQGyc"
	_, err := ValidateJWT(staleJwt, "secretSauce")
	if err == nil {
		t.Error("stale jwt should have failed")
		return
	}
	t.Logf("stale jwt err %v \n", err)
}

func TestWrongKey(t *testing.T) {
	jwt1, err := MakeJWT(uuid.New(), "secretSauce", time.Minute)
	if err != nil {
		t.Errorf("making boring jwt failed %v", err)
		return
	}
	t.Logf("boring jwt is %s \n", jwt1)
	_, err = ValidateJWT(jwt1, "secretSauce")
	if err != nil {
		t.Errorf("validating boring jwt failed %v", err)
		return
	}
	_, err = ValidateJWT(jwt1, "guess")
	if err == nil {
		t.Error("swrong key jwt should have failed")
		return
	}
	t.Logf("wrong key jwt err %v \n", err)
}
