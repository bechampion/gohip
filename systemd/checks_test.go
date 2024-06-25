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

func TestDbConfigIsAlmost7DaysOld(t *testing.T) {
	weekAgo := time.Now().Add(-time.Hour*24*7 + time.Hour)
	details := ClamConfDetails{version: "1.0", sigs: "2.0", DailyCld: weekAgo}

	err := DbConfigAgeCheck(details)

	if err != nil {
		t.Errorf("DailyCld should have been recent enough: \n%v", err)
	}
}

func TestDbConfigIsOver7DaysOld(t *testing.T) {
	weekAgo := time.Now().Add(-time.Hour*24*7 - time.Hour)
	details := ClamConfDetails{version: "1.0", sigs: "2.0", DailyCld: weekAgo}

	err := DbConfigAgeCheck(details)

	if err == nil {
		t.Errorf("DailyCld should have been too old")
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

	return DbFileAgeCheck(clamavDbFile)
}
