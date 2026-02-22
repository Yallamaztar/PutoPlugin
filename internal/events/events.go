package events

import (
	"context"
	"fmt"
	"log"
	"plugin/internal/config"
	"plugin/internal/rcon"
	"plugin/internal/register"

	ps "plugin/internal/service/player"
	ws "plugin/internal/service/wallet"

	"github.com/Yallamaztar/eventsv2/events"
)

func RunEventTailLoop(index int, cfg config.Config, rcon *rcon.RCON, reg *register.Register, playerService *ps.Service, walletService *ws.Service) {
	eventsCh := make(chan events.Event, 128)
	go func() {
		if err := events.TailFileContext(context.Background(), cfg.Server[index].LogPath, true, eventsCh); err != nil {
			log.Fatalf("Failed to tail file: %s", cfg.Server[index].LogPath)
			return
		}
		close(eventsCh)
	}()

	for e := range eventsCh {
		if !cfg.Gambling.Enabled {
			continue
		}

		switch event := e.(type) {
		case *events.PlayerEvent:
			switch event.Command {
			case "J":
				reg.SetClientNum(event.XUID, event.ClientNum)

				guid := rcon.PlayerGUIDByClientNum(event.ClientNum)
				if guid == "" {
					continue
				}

				exists, err := playerService.ExistsByXUID(event.XUID)
				if err != nil {
					continue
				}

				if !exists {
					id, err := playerService.CreatePlayer(event.Name, event.XUID, guid, 0)
					if err != nil {
						continue
					}

					err = walletService.CreateWallet(int(id), int(cfg.Economy.FirstTimeReward))
					if err != nil {
						continue
					}

					rcon.Tell(
						event.ClientNum,
						fmt.Sprintf(
							"^7Created a wallet with ^6%s%d balance",
							cfg.Gambling.Currency,
							cfg.Economy.FirstTimeReward,
						),
					)
				}

			case "Q":
				reg.RemoveClientNum(event.XUID)

			case "say", "sayteam":
			}

		case *events.KillEvent:
		}
	}
}
