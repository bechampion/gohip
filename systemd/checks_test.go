package systemd

import (
	"os"
	"testing"
	"time"
)

func TestFileIs7DaysOld(t *testing.T) {
	clamavError := fileErrorCheck(t, time.Now().Add(-time.Hour*24*7))

	if clamavError != nil {
		t.Errorf("should not be old enough: \n\t%v", clamavError)
	}
}

func TestFileIsAlmost7DaysOld(t *testing.T) {
	clamavError := fileErrorCheck(t, time.Now().Add(-time.Hour*24*7+time.Hour))

	if clamavError != nil {
		t.Errorf("should not be old enough: \n\t%v", clamavError)
	}
}

func TestFileIsMoreThan7DaysOld(t *testing.T) {
	clamavError := fileErrorCheck(t, time.Now().Add(-time.Hour*24*7-time.Hour))

	if clamavError == nil {
		t.Errorf("should not be old enough: \n\t%v", clamavError)
	}
}

func fileErrorCheck(t *testing.T, when time.Time) error {
	nowFilePath := "/tmp/now"

	t.Cleanup(func() {
		os.Remove(nowFilePath)
	})

	_, err := os.Create(nowFilePath)
	if err != nil {
		t.Errorf("%v", err)
	}

	os.Chtimes(nowFilePath, when, when)

	clamavDbFile := ClamavDbFile{path: nowFilePath}

	return DbAgeCheck(clamavDbFile)
}
