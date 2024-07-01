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

func TestDatesParsing(t *testing.T) {
	// extra space before day
	if _, err := parseDate("Mon Jul  1 10:40:06 2024"); err != nil {
		t.Errorf("%v", err)
	}

	// no extra space before day
	if _, err := parseDate("Mon Jul 1 10:40:06 2024"); err != nil {
		t.Errorf("%v", err)
	}

	if _, err := parseDate("Mon Jul 11 10:40:06 2024"); err != nil {
		t.Errorf("%v", err)
	}

	if _, err := parseDate("invalid"); err == nil {
		t.Errorf("should have  been an invalid date%s", "invalid")
	}
}

func parseDailyLine(t *testing.T, line string) {
	var results = parseDailyCvdLine(line)
	if len(results) != 4 {
		t.Errorf("should have 4 results: %v", results)
	}

	cd := ClamConfDetails{}
	cd.DailyCld, _ = parseDate(results[3])
	cd.version = results[1]
	cd.sigs = results[2]

	if cd.version != "27323" {
		t.Errorf("version should have been 27323: %v", cd.version)
	}

	if cd.sigs != "2063707" {
		t.Errorf("sigs should have been 2063707: %v", cd.sigs)
	}

	date, _ := parseDate("Mon Jul  1 10:40:06 2024")
	if cd.DailyCld != date {
		t.Errorf("sigs should have been Mon Jul  1 10:40:06 2024: %v", cd.DailyCld)
	}
}

func TestParseDailyCldLine(t *testing.T) {
	var line = "daily.cld: version 27323, sigs: 2063707, built on Mon Jul  1 10:40:06 2024"
	parseDailyLine(t, line)
}

func TestParseDailyCvdLine(t *testing.T) {
	var line = "daily.cvd: version 27323, sigs: 2063707, built on Mon Jul  1 10:40:06 2024"
	parseDailyLine(t, line)
}
