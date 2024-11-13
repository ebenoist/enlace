package scrapper

import (
	"bytes"
	"os/exec"
)

func Scrape(u string) ([]byte, error) {
	cmd := exec.Command("monolith", u)

	var buf bytes.Buffer
	cmd.Stdout = &buf

	err := cmd.Run()
	return buf.Bytes(), err
}
