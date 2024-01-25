package cfg

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slog"
)

const (
	whiteListFile = "user_allow_list.txt"
)

// The WhiteList interface defines the methods for adding users, deleting users, and initializing the list
type WhiteList interface {
	AddUser(userID string) error    // Adding user to the whitelist
	RemoveUser(userID string) error // Deleting user from the whitelist
	Reloaded() error
	Initialize() error                // Initializing whitelist
	GetWhiteList() []string           // retrieves the whitelist
	IsInWhiteList(userID string) bool // checks if a user is in the whitelist
	StartFileWatcher() error          // monitor
}

// The UserWhiteList struct that implements the WhiteList interface
type UserWhiteList struct {
	steamID  []string
	filename string       // The name of the whitelist file
	fileLock sync.RWMutex // A Mutex to protect against concurrent writes
}

func NewUserWhiteList() *UserWhiteList {
	return &UserWhiteList{
		steamID:  make([]string, 0),
		filename: whiteListFile,
	}
}

// The AddUser method adds a user to the whitelist
func (u *UserWhiteList) AddUser(userID string) error {
	u.fileLock.Lock()
	defer u.fileLock.Unlock()

	// Open the file in a read mode and check if user already exists
	file, err := os.Open(u.filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == userID {
			file.Close()
			// User already exists
			slog.Info("Player Exit", "Player", userID)
			return nil
		}
	}
	file.Close()

	// Open the file in an append mode, or create it if it doesn't exist
	file, err = os.OpenFile(u.filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the user ID to the file
	_, err = file.WriteString(userID + "\n")
	if err != nil {
		return err
	}

	return nil
}

// The RemoveUser method deletes a user from the whitelist
func (u *UserWhiteList) RemoveUser(userID string) error {
	u.fileLock.Lock()
	defer u.fileLock.Unlock()

	// Open the file
	file, err := os.Open(u.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Initialize a reader and a slice of lines
	reader := bufio.NewReader(file)
	lines := []string{}

	// Iterate over the lines in the file
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// If the line content is not the user ID we want to delete, add it to the lines slice
		if strings.Trim(line, "\n") != userID {
			lines = append(lines, line)
		}
	}

	// truncate and reopen the file in write mode
	file, err = os.OpenFile(u.filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write back the lines
	for _, line := range lines {
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// The GetWhiteList method retrieves the whitelist
func (u *UserWhiteList) GetWhiteList() []string {
	u.fileLock.RLock()
	defer u.fileLock.RUnlock()

	// Make a copy of the users
	usersCopy := make([]string, len(u.steamID))
	copy(usersCopy, u.steamID)

	return usersCopy
}

// The IsInWhiteList method checks if a user is in the whitelist
func (u *UserWhiteList) IsInWhiteList(userID string) bool {
	u.fileLock.RLock()
	defer u.fileLock.RUnlock()

	// Check if the user is in the list
	for _, user := range u.steamID {
		if user == userID {
			return true
		}
	}
	return false
}

func (u *UserWhiteList) Reloaded() error {
	u.fileLock.Lock()
	defer u.fileLock.Unlock()

	// Open the file
	file, err := os.Open(u.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// cleared users, and refresh data from file
	u.steamID = []string{}
	for scanner.Scan() {
		u.steamID = append(u.steamID, scanner.Text())
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return scanErr
	}

	return nil
}

// The Initialize method initializes the whitelist
func (u *UserWhiteList) Initialize() error {
	// Check if the file exists
	if _, err := os.Stat(u.filename); os.IsNotExist(err) {
		// The file does not exist, so create it
		file, err := os.Create(u.filename)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}

func (u *UserWhiteList) StartFileWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					err := u.Reloaded()
					if err != nil {
						slog.Error("Failed to reload whitelist", "Error", err)
					} else {
						slog.Info("Reloaded whitelist", "steamID", u.steamID)
					}
				}
			case err := <-watcher.Errors:
				slog.Error("Failed to watch whitelist:", "Error:", err)
			}
		}
	}()

	err = watcher.Add(u.filename)
	if err != nil {
		return err
	}

	return nil
}
