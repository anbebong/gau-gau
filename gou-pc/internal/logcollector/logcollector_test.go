package logcollector

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func createTestLogFile(t *testing.T, path string, entries []ArchiveLogEntry) {
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test log file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			t.Fatalf("failed to encode entry: %v", err)
		}
	}
}

func TestGetLogsPaged(t *testing.T) {
	tmp := "test_archive.log"
	defer os.Remove(tmp)
	entries := []ArchiveLogEntry{
		{Time: "1", AgentID: "a", Message: "msg1"},
		{Time: "2", AgentID: "b", Message: "msg2"},
		{Time: "3", AgentID: "c", Message: "msg3"},
	}
	createTestLogFile(t, tmp, entries)

	logs, total, err := GetLogsPaged(tmp, 1, 2)
	if err != nil {
		t.Fatalf("GetLogsPaged error: %v", err)
	}
	if total != 3 || len(logs) != 2 {
		t.Errorf("unexpected result: total=%d, len(logs)=%d", total, len(logs))
	}
}

func TestRotateLog(t *testing.T) {
	tmp := "test_archive.log"
	old := tmp + ".old"
	defer os.Remove(tmp)
	defer os.Remove(old)
	entries := []ArchiveLogEntry{
		{Time: "1", AgentID: "a", Message: "msg1"},
	}
	createTestLogFile(t, tmp, entries)

	if err := RotateLog(tmp); err != nil {
		t.Fatalf("RotateLog error: %v", err)
	}
	// Check old file exists and new file is empty
	if _, err := os.Stat(old); err != nil {
		t.Errorf("old log file not found after rotate")
	}
	info, err := os.Stat(tmp)
	if err != nil || info.Size() != 0 {
		t.Errorf("new log file not empty after rotate")
	}
}

func TestGetLogsPaged_RealFileCopy(t *testing.T) {
	src := "../../etc/archive.log"
	dst := "test_archive_real_copy.log"
	// Copy file thực tế sang file tạm
	in, err := os.Open(src)
	if err != nil {
		t.Skipf("Không tìm thấy file log thực tế: %v", err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		t.Fatalf("Không tạo được file copy: %v", err)
	}
	_, err = io.Copy(out, in)
	if err != nil {
		out.Close()
		t.Fatalf("Copy file thất bại: %v", err)
	}
	out.Close()
	defer os.Remove(dst)

	logs, total, err := GetLogsPaged(dst, 4, 10)
	if err != nil {
		t.Fatalf("GetLogsPaged error: %v", err)
	}
	t.Logf("Total logs: %d, First page count: %d", total, len(logs))
	for i, l := range logs {
		t.Logf("Log %d: %+v", i, l)
	}
}
