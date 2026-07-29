package solutionx100c

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AttendanceRecord struct {
	NIK       string    `json:"nik"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
}

type Client struct {
	IPAddress string
	Port      int
	ComKey    int
	Timeout   time.Duration
}

func NewClient(ip string, port int, comKey int) *Client {
	if port <= 0 {
		port = 80
	}
	return &Client{
		IPAddress: ip,
		Port:      port,
		ComKey:    comKey,
		Timeout:   3 * time.Second,
	}
}

// FetchAttLog connects to Solution X100C via SOAP (/iWsService) and retrieves attendance logs.
func (c *Client) FetchAttLog() ([]AttendanceRecord, error) {
	url := fmt.Sprintf("http://%s:%d/iWsService", c.IPAddress, c.Port)
	body := fmt.Sprintf(`<GetAttLog>
  <ArgComKey xsi:type="xsd:integer">%d</ArgComKey>
  <Arg><PIN xsi:type="xsd:integer">All</PIN></Arg>
</GetAttLog>`, c.ComKey)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml")

	client := &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: c.Timeout,
			}).DialContext,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed to %s:%d: %w", c.IPAddress, c.Port, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return ParseAttLogResponse(string(respBytes))
}

// ClearAttLog sends <ClearData> command to flush attendance logs from machine memory.
func (c *Client) ClearAttLog() error {
	url := fmt.Sprintf("http://%s:%d/iWsService", c.IPAddress, c.Port)
	body := fmt.Sprintf(`<ClearData>
  <ArgComKey xsi:type="xsd:integer">%d</ArgComKey>
  <Arg><Value xsi:type="xsd:integer">3</Value></Arg>
</ClearData>`, c.ComKey)

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("failed to create clear request: %w", err)
	}

	req.Header.Set("Content-Type", "text/xml")

	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send clear command to %s:%d: %w", c.IPAddress, c.Port, err)
	}
	defer resp.Body.Close()

	return nil
}

// ParseAttLogResponse parses the XML response from <GetAttLogResponse>.
func ParseAttLogResponse(xmlData string) ([]AttendanceRecord, error) {
	var records []AttendanceRecord

	// Extract content inside <GetAttLogResponse> ... </GetAttLogResponse>
	startTag := "<GetAttLogResponse>"
	endTag := "</GetAttLogResponse>"

	idxStart := strings.Index(xmlData, startTag)
	if idxStart == -1 {
		// Response might not contain logs or might be empty
		return records, nil
	}

	content := xmlData[idxStart+len(startTag):]
	idxEnd := strings.Index(content, endTag)
	if idxEnd != -1 {
		content = content[:idxEnd]
	}

	// Split by <Row>
	rows := strings.Split(content, "<Row>")
	for _, row := range rows {
		if !strings.Contains(row, "</Row>") {
			continue
		}

		nik := parseXMLTag(row, "PIN")
		dateTimeStr := parseXMLTag(row, "DateTime")
		statusStr := parseXMLTag(row, "Status")

		if nik == "" || dateTimeStr == "" {
			continue
		}

		// Parse timestamp e.g. "2026-07-29 05:45:12"
		ts, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
		if err != nil {
			// Try alternative date format "2006-01-02 15:04"
			ts, err = time.Parse("2006-01-02 15:04", dateTimeStr)
			if err != nil {
				continue
			}
		}

		status, _ := strconv.Atoi(statusStr)

		records = append(records, AttendanceRecord{
			NIK:       strings.TrimSpace(nik),
			Timestamp: ts,
			Status:    status,
		})
	}

	return records, nil
}

func parseXMLTag(data, tag string) string {
	openTag := "<" + tag + ">"
	closeTag := "</" + tag + ">"

	start := strings.Index(data, openTag)
	if start == -1 {
		return ""
	}

	end := strings.Index(data, closeTag)
	if end == -1 || end < start {
		return ""
	}

	return data[start+len(openTag) : end]
}
