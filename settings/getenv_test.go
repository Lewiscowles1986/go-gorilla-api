package settings

import (
    "testing"
)

func TestGetEnvPathVar(t *testing.T) {
    result := Getenv("PATH", "")
    if len(result) < 0 {
        t.Errorf("Expected OS PATH variable to be set always, contact a sysadmin")
    }
}

func TestGetEnvNonExistantVar(t *testing.T) {
    result := Getenv("HSDUIHFDSKAJAKCNCKBCK", "Hooray")
    if result != "Hooray" {
        t.Errorf("Expected you wouldn't have gibberish in your environment. "+
          "Please ensure 'HSDUIHFDSKAJAKCNCKBCK' is not a valid entry before "+
          "reporting a bug")
    }
}
