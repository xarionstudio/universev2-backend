package solutionx100c

import (
	"testing"
	"time"
)

func TestParseAttLogResponse(t *testing.T) {
	xmlData := `<?xml version="1.0"?>
<GetAttLogResponse>
  <Row>
    <PIN>503264133</PIN>
    <DateTime>2026-07-29 05:45:12</DateTime>
    <Verified>1</Verified>
    <Status>0</Status>
    <WorkCode>0</WorkCode>
  </Row>
  <Row>
    <PIN>503264134</PIN>
    <DateTime>2026-07-29 17:15:00</DateTime>
    <Verified>1</Verified>
    <Status>1</Status>
    <WorkCode>0</WorkCode>
  </Row>
</GetAttLogResponse>`

	records, err := ParseAttLogResponse(xmlData)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(records))
	}

	if records[0].NIK != "503264133" {
		t.Errorf("Expected NIK 503264133, got %s", records[0].NIK)
	}

	expectedTime := time.Date(2026, 7, 29, 5, 45, 12, 0, time.UTC)
	if !records[0].Timestamp.Equal(expectedTime) {
		t.Errorf("Expected timestamp %v, got %v", expectedTime, records[0].Timestamp)
	}

	if records[1].NIK != "503264134" {
		t.Errorf("Expected NIK 503264134, got %s", records[1].NIK)
	}
}
