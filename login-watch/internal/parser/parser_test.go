package parser

import "testing"

func TestParseFailedPassword(t *testing.T) {
	line := `Jun 30 14:03:04 hartserver sshd[1234]: Failed password for root from 192.168.1.50 port 54321 ssh2`
	e := Parse(line)
	if e == nil {
		t.Fatal("expected event, got nil")
	}
	if e.SrcIP != "192.168.1.50" {
		t.Errorf("expected IP 192.168.1.50, got %s", e.SrcIP)
	}
	if e.Username != "root" {
		t.Errorf("expected username root, got %s", e.Username)
	}
	if e.EventType != "Failed password" {
		t.Errorf("unexpected event type: %s", e.EventType)
	}
}

func TestParseInvalidUser(t *testing.T) {
	line := `Jun 30 14:03:05 hartserver sshd[1235]: Invalid user admin from 10.0.0.5 port 22`
	e := Parse(line)
	if e == nil {
		t.Fatal("expected event, got nil")
	}
	if e.SrcIP != "10.0.0.5" {
		t.Errorf("expected IP 10.0.0.5, got %s", e.SrcIP)
	}
	if e.Username != "admin" {
		t.Errorf("expected username admin, got %s", e.Username)
	}
}

func TestParseFailedPasswordInvalidUser(t *testing.T) {
	line := `Jun 30 14:03:06 hartserver sshd[1236]: Failed password for invalid user deploy from 10.0.0.6 port 22 ssh2`
	e := Parse(line)
	if e == nil {
		t.Fatal("expected event, got nil")
	}
	if e.Username != "deploy" {
		t.Errorf("expected username deploy, got %s", e.Username)
	}
}

func TestParseIgnoresNonFailure(t *testing.T) {
	line := `Jun 30 14:03:06 hartserver sshd[1236]: Accepted publickey for chris from 192.168.1.1 port 54322 ssh2`
	e := Parse(line)
	if e != nil {
		t.Fatal("expected nil for successful login line")
	}
}

func TestParseIgnoresNoIP(t *testing.T) {
	line := `Jun 30 14:03:07 hartserver sshd[1237]: Failed password for root`
	e := Parse(line)
	if e != nil {
		t.Fatal("expected nil when no IP present")
	}
}
