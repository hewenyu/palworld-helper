package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hewenyu/palworld-helper/cfg"
	"github.com/hewenyu/palworld-helper/palrcon"
)

var (
	interval time.Duration
	timeout  time.Duration

	uconvLatin = os.Getenv("UCONV_LATIN") != "false"

	userList   *cfg.UserWhiteList
	palRCON    palrcon.PalRCON
	cfgManager cfg.ConfigManager
)

func initializeEnvironment() error {
	cfgManager = cfg.NewConfigManager()
	configs := cfgManager.ReadConfig()

	var err error
	timeout, err = time.ParseDuration(configs.Timeout)
	if err != nil {
		return fmt.Errorf("failed to parse timeout: %w", err)
	}

	interval, err = time.ParseDuration(configs.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval: %w", err)
	}

	userList = cfg.NewUserWhiteList()
	if err := userList.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize user list: %w", err)
	}
	if err := userList.StartFileWatcher(); err != nil {
		return fmt.Errorf("failed to start file watcher: %w", err)
	}

	palRCON = palrcon.NewPalRCON(configs.RCONSettings.Endpoint,
		configs.RCONSettings.Password)
	palRCON.SetTimeout(timeout)

	msg, err := palRCON.Info()

	if err != nil {
		return err
	}

	slog.Info(msg)

	return nil
}

// map
func makePlayerMap(players []palrcon.Player) map[string]palrcon.Player {
	m := make(map[string]palrcon.Player)

	for _, player := range players {
		if player.PlayerUID != "00000000" {
			m[player.SteamID] = player
		}
	}

	return m
}

func retriedBroadcast(message string) error {
	message = escapeString(message)
	var err error
	for i := 0; i < 10; i++ {
		err = palRCON.Broadcast(message)
		if err != nil {
			slog.Error("failed to broadcast", "error", err)
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to broadcast: %w", err)
}

func handleLeftPlayer(player palrcon.Player, playersMap map[string]palrcon.Player) {
	if _, ok := playersMap[player.PlayerUID]; !ok {
		slog.Info("Player left", "player", player)

		if err := retriedBroadcast(fmt.Sprintf("Left: %s", player.Name)); err != nil {
			slog.Error("Failed to broadcast", "error", err)
		}
	}
}

func handleNewPlayer(player palrcon.Player, prev map[string]palrcon.Player) error {
	if _, ok := prev[player.SteamID]; !ok {
		if err := retriedBroadcast(fmt.Sprintf("Joined: %s", player.Name)); err != nil {
			slog.Error("Failed to broadcast", "error", err)
		}
		if userList.IsInWhiteList(player.SteamID) {
			slog.Info("Player joined", player.Name, player.SteamID)
		} else {
			slog.Info("The User is not in the white list", player.Name, player.SteamID)
			if err := palRCON.KickPlayer(player.SteamID); err != nil {
				slog.Error("Failed to kick player", "error", err)
			}
			// if err := palRCON.BanPlayer(player.SteamID); err != nil {
			// 	slog.Error("Failed to ban player", "error", err)
			// }
		}
	}
	return nil
}

func runUconvLatin(s string) string {
	var out strings.Builder
	cmd := exec.Command("uconv", "-x", "latin")
	cmd.Stdin = strings.NewReader(s)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		slog.Error("failed to run uconv", "error", err)
		return s
	}

	return out.String()
}

func escapeString(s string) string {
	if uconvLatin {
		s = runUconvLatin(s)
	}
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.TrimSpace(s)

	runes := []rune(s)
	for i := range runes {
		b := []byte(string(runes[i]))

		if len(b) != 1 {
			runes[i] = '*'
		}
	}

	return string(runes)
}

func main() {

	var prev map[string]palrcon.Player

	if err := initializeEnvironment(); err != nil { // improved error handling
		slog.Error(err.Error())
		os.Exit(1)
	}

	for {
		players, err := palRCON.GetPlayers()

		if err != nil {
			slog.Error("Failed to get players", "error", err)
			time.Sleep(interval)
			continue
		}

		playersMap := makePlayerMap(players)

		if prev == nil {
			prev = playersMap
			time.Sleep(interval)
			continue
		}

		for _, player := range players {
			handleNewPlayer(player, prev)
		}

		for _, player := range prev {
			handleLeftPlayer(player, prev)
		}

		prev = playersMap

		time.Sleep(interval)
	}
}
