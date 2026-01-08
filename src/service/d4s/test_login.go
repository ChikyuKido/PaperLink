package d4s

import (
	"os/exec"
	"paperlink/db/entity"
	"strings"
)

func TestLogin(acc *entity.Digi4SchoolAccount) bool {
	cmd := exec.Command("./integrations/d4s", "test-login", acc.Username, acc.Password)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to execute test-login for user %s: %v, output: %s", acc.Username, err, string(output))
		return false
	}

	outStr := string(output)
	outStr = strings.TrimSpace(outStr)
	if outStr == "1" {
		return true
	} else if outStr == "0" {
		return false
	}

	log.Printf("unexpected output from test-login for user %s: %q", acc.Username, outStr)
	return false
}
