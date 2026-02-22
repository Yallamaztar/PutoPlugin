package register

import (
	"fmt"
	"plugin/internal/config"
	"plugin/internal/rcon"
	"strings"
	"sync"
)

type Level uint8

type Handler func(
	clientNum uint8,
	playerID int,
	playerName string,
	xuid string,
	level Level,
	args []string,
)

type Command struct {
	Name     string
	Aliases  []string
	MinLevel Level
	Help     string
	MinArgs  uint8
	Handler  Handler
}

type commands map[string]Command
type clients map[string]uint8

type Register struct {
	commands commands
	clients  clients
	rc       *rcon.RCON
	cfg      config.Config
	mu       sync.RWMutex
}

func New(cfg config.Config, rc *rcon.RCON) *Register {
	return &Register{
		commands: make(commands),
		clients:  make(clients),
		rc:       rc,
		cfg:      cfg,
	}
}

func (r *Register) RegisterCommand(cmd Command) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.commands[strings.ToLower(cmd.Name)] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[strings.ToLower(alias)] = cmd
	}
}

func (r *Register) Execute(
	clientNum uint8,
	playerID int,
	playerName string,
	xuid string,
	level Level,
	command string,
	args []string,
) {
	r.mu.RLock()
	cmd, ok := r.commands[strings.ToLower(command)]
	r.mu.RUnlock()

	if !ok {
		return // unknown command
	}

	if !r.hasPermission(level, cmd.MinLevel) {
		r.tell(clientNum, fmt.Sprintf(
			"You ^1don't ^7have permission for !%s",
			cmd.Name,
		))
		return
	}

	if len(args) < int(cmd.MinArgs) {
		if cmd.Help != "" {
			r.tell(clientNum, cmd.Help)
		}
		return
	}

	cmd.Handler(clientNum, playerID, playerName, xuid, level, args)
}

func (r *Register) SetClientNum(xuid string, clientNum uint8) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[xuid] = clientNum
}

func (r *Register) GetClientNum(xuid string) (uint8, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cn, ok := r.clients[xuid]
	return cn, ok
}

func (r *Register) RemoveClientNum(xuid string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, xuid)
}

func (r *Register) hasPermission(level, required Level) bool {
	return level >= required
}

func (r *Register) tell(clientNum uint8, msg string) {
	if r.rc != nil {
		r.rc.Tell(clientNum, msg)
	}
}
