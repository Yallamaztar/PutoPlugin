package main

import (
	"log"
	"plugin/internal/commands"
	"plugin/internal/config"
	"plugin/internal/database"
	"plugin/internal/events"
	"plugin/internal/rcon"
	"plugin/internal/register"
	"sync"

	br "plugin/internal/repository/bank"
	pr "plugin/internal/repository/player"
	sr "plugin/internal/repository/stats"
	wr "plugin/internal/repository/wallet"

	bs "plugin/internal/service/bank"
	ps "plugin/internal/service/player"
	ss "plugin/internal/service/stats"
	ws "plugin/internal/service/wallet"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("failed database migration")
	}

	player := ps.New(pr.New(db))
	wallet := ws.New(wr.New(db))
	bank := bs.New(br.New(db))
	playerStats := ss.NewPlayerStats(sr.NewPlayerStats(db))
	gambleStats := ss.NewGamblingStats(sr.NewGamblingStats(db))

	log.Println("Loading config.yaml")
	cfg, err := config.Setup()
	if err != nil {
		log.Fatal("config setup failed:", err)
	}

	var wg sync.WaitGroup
	for i, s := range cfg.Server {
		log.Printf("Connecting to RCON on %s\n", s.Host)
		rc, err := rcon.New(s.Host, s.Password, *cfg)
		if err != nil {
			continue
		}

		reg := register.New(*cfg, rc)
		commands.RegisterClientCommands(*cfg, rc, reg, player, wallet, bank, playerStats, gambleStats)

		log.Println("Starting Plugin")
		wg.Add(1)
		go func(index int, rc *rcon.RCON) {
			defer wg.Done()
			defer rc.Close()
			events.RunEventTailLoop(index, *cfg, rc, reg, player, wallet)
		}(i, rc)
	}

	wg.Wait()

}
