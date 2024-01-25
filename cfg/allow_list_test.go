package cfg

import (
	"bufio"
	"os"
	"testing"
)

func TestWhiteList(t *testing.T) {
	wl := NewUserWhiteList()

	// Initial run should fail if there is no file, unless we run Initialize()
	_, err := os.Stat(wl.filename)
	if err == nil {
		t.Errorf("File %v should not exist yet", wl.filename)
	}

	err = wl.Initialize()
	if err != nil {
		t.Errorf("Initialize() failed: %v", err)
	}

	_, err = os.Stat(wl.filename)
	if err != nil {
		t.Errorf("File %v should exist after Initialize(), but got error %v", wl.filename, err)
	}

	err = wl.AddUser("testuser")
	if err != nil {
		t.Errorf("AddUser() failed: %v", err)
	}

	// Check if "testuser" was added
	file, err := os.Open(whiteListFile)
	if err != nil {
		t.Errorf("OpenFile() failed: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("ReadString() failed: %v", err)
	}
	if line != "testuser\n" {
		t.Errorf("AddUser() did not add the user correctly: got %s want %s", line, "testuser\n")
	}

	// Try to remove an user that does not exist
	err = wl.RemoveUser("nonexistinguser")
	if err != nil {
		t.Errorf("RemoveUser() failed: %v", err)
	}

	// Test removing the added user
	err = wl.RemoveUser("testuser")
	if err != nil {
		t.Errorf("RemoveUser() failed: %v", err)
	}

	// Check if "testuser" was removed
	file, err = os.Open(whiteListFile)
	if err != nil {
		t.Errorf("OpenFile() failed: %v", err)
	}
	defer file.Close()

	reader = bufio.NewReader(file)
	byteLine, _, err := reader.ReadLine()
	if err != nil && err.Error() != "EOF" {
		t.Errorf("ReadLine() failed: %v", err)
	}
	line = string(byteLine)
	if len(line) != 0 {
		t.Errorf("RemoveUser() did not remove the user correctly: got %s want '' ", line)
	}
}

func TestAddUser(t *testing.T) {
	wl := &UserWhiteList{filename: whiteListFile}

	// Initialize
	err := wl.Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Test AddUser
	err = wl.AddUser("testuser")
	if err != nil {
		t.Fatalf("AddUser() failed: %v", err)
	}

	// Check if "testuser" was added
	file, err := os.Open(whiteListFile)
	if err != nil {
		t.Fatalf("OpenFile() failed: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("ReadString() failed: %v", err)
	}
	if line != "testuser\n" {
		t.Errorf("AddUser() did not add the user correctly: got %s want %s", line, "testuser\n")
	}
}

func TestRemoveUser(t *testing.T) {
	wl := &UserWhiteList{filename: whiteListFile}

	// Initialize
	err := wl.Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Add a user to remove
	err = wl.AddUser("testuser")
	if err != nil {
		t.Fatalf("AddUser() failed: %v", err)
	}

	// Test removing the added user
	err = wl.RemoveUser("testuser")
	if err != nil {
		t.Fatalf("RemoveUser() failed: %v", err)
	}

	// Check if "testuser" was removed
	file, err := os.Open(whiteListFile)
	if err != nil {
		t.Fatalf("OpenFile() failed: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	byteLine, _, err := reader.ReadLine()
	if err != nil && err.Error() != "EOF" {
		t.Fatalf("ReadLine() failed: %v", err)
	}
	line := string(byteLine)
	if len(line) != 0 {
		t.Errorf("RemoveUser() did not remove the user correctly: got %s want '' ", line)
	}
}
