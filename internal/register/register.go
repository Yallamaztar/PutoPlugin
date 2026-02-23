package register

import (
	"fmt"
	"plugin/internal/config"
	"plugin/internal/rcon"
	"strings"
	"sync"
)

type Handler func(
	clientNum uint8,
	playerID int,
	playerName string,
	xuid string,
	level int,
	args []string,
)

type Command struct {
	Name     string
	Aliases  []string
	MinLevel int
	Help     string
	MinArgs  uint8
	Handler  Handler
}

type commands map[string]*Command
type clients map[string]uint8

type Register struct {
	commands commands
	clients  clients

	rc  *rcon.RCON
	cfg config.Config

	mu sync.RWMutex
}

func New(cfg config.Config, rc *rcon.RCON) *Register {
	return &Register{
		commands: make(commands),
		clients:  make(clients),

		rc:  rc,
		cfg: cfg,

		mu: sync.RWMutex{},
	}
}
func (r *Register) RegisterCommand(cmd Command) {
	if cmd.Handler == nil {
		panic("command " + cmd.Name + " registered with nil handler")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	c := cmd
	r.commands[strings.ToLower(c.Name)] = &c
	for _, alias := range c.Aliases {
		r.commands[strings.ToLower(alias)] = &c
	}
}

func (r *Register) Execute(
	clientNum uint8,
	playerID int,
	playerName string,
	xuid string,
	level int,
	command string,
	args []string,
) {
	r.mu.RLock()
	cmd, ok := r.commands[strings.ToLower(command)]
	r.mu.RUnlock()

	if !ok {
		return
	}

	if cmd.Handler == nil {
		return
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

func (r *Register) hasPermission(level, required int) bool {
	return level >= required
}

func (r *Register) tell(clientNum uint8, msg string) {
	if r.rc != nil {
		r.rc.Tell(clientNum, msg)
	}
}
