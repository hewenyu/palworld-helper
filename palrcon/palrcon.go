package palrcon

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorcon/rcon"
)

type Player struct {
	Name      string
	PlayerUID string // might be int64
	SteamID   string // might be int64
}

type PalRCON interface {
	Info() (string, error)            // Show server information.
	Save() error                      // Save the world data.
	GetPlayers() ([]Player, error)    // ShowPlayers often times out, so ignore the error
	Broadcast(message string) error   // Send message to all player in the server
	KickPlayer(steamID string) error  // Kick player from the server..
	BanPlayer(steamID string) error   // BAN player from the server.
	SetTimeout(timeout time.Duration) //
}

func NewPalRCON(endpoint, password string) PalRCON {
	return &palRCON{
		endpoint: endpoint,
		password: password,
	}
}

type palRCON struct {
	endpoint string
	password string

	timeout time.Duration
}

func (p *palRCON) execute(command string) (string, error) {
	// rcon of palworld in unstable
	// so the connection isn't reused

	rconn, err := rcon.Dial(
		p.endpoint, p.password,
		rcon.SetDialTimeout(p.timeout),
		rcon.SetDeadline(p.timeout),
	)

	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", p.endpoint, err)
	}
	defer rconn.Close()

	result, err := rconn.Execute(command)

	if err != nil {
		return result, fmt.Errorf("failed to execute the command: %w", err)
	}

	if len(result) == 0 {
		return result, nil
	}

	raw := []byte(result)
	i := len(raw)
	for ; i > 0; i-- {
		if raw[i-1] != 0 {
			break
		}
	}

	return string(raw[:i]), nil
}

func (p *palRCON) GetPlayers() ([]Player, error) {
	// ShowPlayers often times out, so ignore the error

	result, err := p.execute("ShowPlayers")

	if len(result) == 0 && err != nil {
		return nil, err
	}

	lines := strings.Split(result, "\n")[1:] // skip header (name,playeruid,steamid)

	players := make([]Player, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		fields := strings.Split(line, ",")

		players = append(players, Player{
			Name:      strings.Join(fields[:len(fields)-2], ","),
			PlayerUID: fields[len(fields)-2],
			SteamID:   fields[len(fields)-1],
		})
	}

	return players, nil
}

func (p *palRCON) Broadcast(message string) error {
	// Send message to all player in the server
	_, err := p.execute(fmt.Sprintf("Broadcast %s", message))

	return err
}

func (p *palRCON) KickPlayer(steamid string) error {
	// Kick player from the server..
	_, err := p.execute(fmt.Sprintf("KickPlayer %s", steamid))

	return err
}

func (p *palRCON) BanPlayer(steamid string) error {
	// BAN player from the server.
	_, err := p.execute(fmt.Sprintf("BanPlayer %s", steamid))

	return err
}

func (p *palRCON) Save() error {
	// BAN player from the server.
	_, err := p.execute("Save")

	return err
}

func (p *palRCON) Info() (message string, err error) {
	// Show server information.
	result, err := p.execute("Info")

	if len(result) == 0 && err != nil {
		return "", err
	}

	lines := strings.Split(result, "\n")

	return lines[0], err
}

func (p *palRCON) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}
