package cigitlab

import "testing"

// configure should set package vars only
// once
func TestConfigure(t *testing.T) {
	host, api, token := "hostA", "apiA", "tokenA"
	// set some values
	Configure(host, api, token)
	// now try to change them
	Configure("", "", "")
	if host != apiHost || api != apiUrl || token != apiToken {
		t.Errorf("expected variable not set properly or changed: %s %s %S", apiHost, apiUrl, apiToken)
		return
	}

}

// trigger should work only with
// defined commands
func TestTrigger(t *testing.T) {
	cmd := "dance"
	_, err := Trigger(cmd, "999", "polka")
	if err == ErrProjID {
		return
	}

	t.Errorf("got %s expected error: %s", err, ErrWrongCMD)
}
