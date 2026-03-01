package commands

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"plugin/internal/commands/gamble"
	"plugin/internal/commands/pay"
	"plugin/internal/config"
	"plugin/internal/discord/webhook"
	"plugin/internal/helpers"
	"plugin/internal/rcon"
	"plugin/internal/register"
	rp "plugin/internal/repository/player"
	"plugin/internal/service/bank"
	"plugin/internal/service/link"
	"plugin/internal/service/player"
	"plugin/internal/service/stats"
	"plugin/internal/service/wallet"
	"strings"
)

const (
	LevelUser = iota
	levelAdmin
	LevelOwner
)

const unknownErr = "an ^1error ^7occurred, please ^1try again ^7later"

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

func generateCode() (string, error) {
	code := make([]byte, codeLength)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}

	return string(code), nil
}

func resolveClientNum(
	rcon *rcon.RCON,
	reg *register.Register,
	clientNum uint8,
	args []string,
) (int, error) {
	if len(args) == 0 {
		return int(clientNum), nil
	}

	query := strings.TrimSpace(strings.Join(args, " "))
	target := reg.FindPlayer(query)
	if target == nil {
		return -1, fmt.Errorf("Couldnt find ^6%s", query)
	}

	cn := rcon.ClientNumByGUID(target.GUID)
	if cn == -1 {
		return -1, fmt.Errorf("Couldnt resolve ^6%s", query)
	}

	return cn, nil
}

func resolveClientNums(
	rcon *rcon.RCON,
	reg *register.Register,
	clientNum uint8,
	args []string,
) (int, int, error) {
	if len(args) == 0 {
		return -1, -1, fmt.Errorf("At least one target is required")
	}

	resolve := func(query string) (int, error) {
		target := reg.FindPlayer(query)
		if target == nil {
			return -1, fmt.Errorf("Couldn't find ^6%s", query)
		}

		cn := rcon.ClientNumByGUID(target.GUID)
		if cn == -1 {
			return -1, fmt.Errorf("Couldn't resolve ^6%s", query)
		}

		return cn, nil
	}

	if len(args) == 1 {
		cn2, err := resolve(strings.TrimSpace(args[0]))
		if err != nil {
			return -1, -1, err
		}
		return int(clientNum), cn2, nil
	}

	cn1, err := resolve(strings.TrimSpace(args[0]))
	if err != nil {
		return -1, -1, err
	}

	cn2, err := resolve(strings.TrimSpace(args[1]))
	if err != nil {
		return -1, -1, err
	}

	if cn1 == cn2 {
		return -1, -1, fmt.Errorf("Targets must be different players")
	}

	return cn1, cn2, nil
}

func RegisterAdminCommands(
	cfg *config.Config,
	rcon *rcon.RCON,
	reg *register.Register,

	player *player.Service,
	wallet *wallet.Service,
	bank *bank.Service,
	link *link.Service,
) {
	reg.RegisterCommand(register.Command{
		Name:     "freeze",
		Aliases:  []string{"fz", "freez"},
		MinLevel: levelAdmin,
		MinArgs:  0,
		Help:     "Usage: ^6!freeze ^7<player>",
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			cn, err := resolveClientNum(rcon, reg, clientNum, args)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}
			rcon.SetInDvar(fmt.Sprintf("freeze %d", cn))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "dropgun",
		Aliases:  []string{"dg", "drop"},
		MinLevel: levelAdmin,
		MinArgs:  0,
		Help:     "Usage: ^6!dropgun ^7<player>",
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			cn, err := resolveClientNum(rcon, reg, clientNum, args)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.SetInDvar(fmt.Sprintf("dropgun %d", cn))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "setspeed",
		Aliases:  []string{"ss", "sets", "sspeed"},
		MinLevel: levelAdmin,
		MinArgs:  0,
		Help:     "Usage ^6!setspeed ^7<player> <amount>",
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			cn, err := resolveClientNum(rcon, reg, clientNum, args)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.SetInDvar(fmt.Sprintf("setspeed %d", cn))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "killplayer",
		Aliases:  []string{"kpl", "kplayer", "killp"},
		MinLevel: levelAdmin,
		MinArgs:  0,
		Help:     "Usage: ^6!killplayer ^7<player>",
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			cn, err := resolveClientNum(rcon, reg, clientNum, args)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.SetInDvar(fmt.Sprintf("killplayer %d", cn))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "hide",
		Aliases:  []string{"hd", "hid", "invisible", "invis"},
		MinLevel: levelAdmin,
		MinArgs:  0,
		Help:     "Usage: ^6!hide ^7<player>",
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			cn, cn2, err := resolveClientNums(rcon, reg, clientNum, args)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.SetInDvar(fmt.Sprintf("swap %d %d", cn, cn2))
		},
	})

}

func RegisterClientCommands(
	cfg *config.Config,
	rcon *rcon.RCON,
	reg *register.Register,

	player *player.Service,
	wallet *wallet.Service,
	bank *bank.Service,
	link *link.Service,

	playerStats *stats.PlayeStatsService,
	gambleStats *stats.GamblingStatsService,
	walletStats *stats.WalletStatsService,

	webhook *webhook.Webhook,
) {
	reg.RegisterCommand(register.Command{
		Name:     "link",
		Aliases:  []string{"lnk", "lk", "linkdc"},
		MinLevel: LevelUser,
		MinArgs:  0,
		Help:     "Usage: ^6!link",
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			discordID, err := player.GetDiscordIDByID(id)
			if err != nil {
				rcon.Tell(clientNum, "^1Error ^7checking your account, try again later")
				return
			}

			if discordID != "" {
				rcon.Tell(clientNum, "You have ^6already linked ^7your account")
				return
			}

			code, err := generateCode()
			if err != nil {
				rcon.Tell(clientNum, "^1Failed ^7to generate link code, ^6try again ^7later")
				return
			}

			if err = link.CreateLink(id, code); err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			rcon.Tell(clientNum, fmt.Sprintf("Your code is: ^6%s", code))
			rcon.Tell(clientNum, "use ^6/link <code> ^7in discord to link your account")
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "gamble",
		Aliases:  []string{"g", "cf", "coinflip"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!gamble ^7<amount>",
		MinArgs:  1,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			balance, err := wallet.GetBalance(id)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			amount, err := helpers.ParseAmountArg(args[0], int64(balance))
			if err != nil {
				rcon.Tell(clientNum, fmt.Sprintf("%s (%q)", err, args[0]))
				return
			}

			res, err := gamble.Gamble(id, name, amount, cfg, player, wallet, bank, playerStats, gambleStats, walletStats, webhook)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.Tell(clientNum, res.Message)
			if res.Won {
				rcon.Say(fmt.Sprintf("%s just ^6won ^7%s%d!", name, cfg.Gambling.Currency, res.Amount))
			} else {
				rcon.Say(fmt.Sprintf("%s just ^6lost ^7%s%d!", name, cfg.Gambling.Currency, res.Amount))
			}

		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "pay",
		Aliases:  []string{"pp", "payplayer", "transfer"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!pay <player> <amount>",
		MinArgs:  2,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			balance, err := wallet.GetBalance(id)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			amount, err := helpers.ParseAmountArg(args[0], int64(balance))
			if err != nil {
				rcon.Tell(clientNum, fmt.Sprintf("%s ^6(%q)", err, args[0]))
				return
			}

			t := reg.FindPlayer(args[0])
			if t == nil {
				rcon.Tell(clientNum, fmt.Sprintf("player ^6%s ^7couldnt be found", args[0]))
				return
			}

			target, err := player.GetPlayerByGUID(t.GUID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			res, err := pay.Pay(id, target.ID, amount, cfg, player, wallet, walletStats, webhook)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.Tell(clientNum, res.Message)
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "balance",
		Aliases:  []string{"bal", "balanc", "money"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!balance <player>",
		MinArgs:  0,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			targetGUID := ""
			targetName := name

			if len(args) > 0 {
				query := strings.TrimSpace(strings.Join(args, " "))
				target := reg.FindPlayer(query)
				if target == nil {
					rcon.Tell(clientNum, fmt.Sprintf("Couldnt find player ^6%s", query))
					return
				}

				targetGUID = target.GUID
				targetName = target.Name
			}

			var p *rp.Player
			var err error

			if targetGUID != "" {
				p, err = player.GetPlayerByGUID(targetGUID)
			} else {
				p, err = player.GetPlayerByID(id)
			}

			if err != nil || p == nil {
				rcon.Tell(clientNum, "Player account not found")
				return
			}

			bal, err := wallet.GetBalance(p.ID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			if targetGUID != "" {
				rcon.Tell(clientNum, fmt.Sprintf("%s's balance: ^6%s%d", targetName, cfg.Gambling.Currency, bal))
			} else {
				rcon.Tell(clientNum, fmt.Sprintf("Your balance: ^6%s%d", cfg.Gambling.Currency, bal))
			}
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "bankbalance",
		Aliases:  []string{"bb", "bank", "bankbal"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!bankbalance",
		MinArgs:  1,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			bal, err := bank.GetBalance()
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			rcon.Tell(clientNum, fmt.Sprintf("bank balance is ^6%s%d", cfg.Gambling.Currency, bal))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "discord",
		Aliases:  []string{"dc", "disc"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!discord",
		MinArgs:  1,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			if cfg.Discord.Enabled {
				rcon.Tell(clientNum, fmt.Sprintf("^6%s", cfg.Discord.InviteLink))
			}
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "richest",
		Aliases:  []string{"rich"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!richest",
		MinArgs:  1,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			wallets, err := wallet.GetTop5RichestWallets()
			if err != nil {
				rcon.Tell(clientNum, "Couldnt get wallets")
				return
			}

			for i, w := range wallets {
				rcon.Tell(clientNum, fmt.Sprintf("[%d] %s %s%s", i+1, w.Name, cfg.Gambling.Currency, helpers.FormatMoney(w.Balance)))
			}
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "poorest",
		Aliases:  []string{"poor"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!poorest",
		MinArgs:  1,
		Handler: func(clientNum uint8, id int, name, xuid string, level int, args []string) {
			wallets, err := wallet.GetTop5PoorestWallets()
			if err != nil {
				rcon.Tell(clientNum, "Couldnt get wallets")
				return
			}

			for i, w := range wallets {
				rcon.Tell(clientNum, fmt.Sprintf("[%d] %s %s%s", i+1, w.Name, cfg.Gambling.Currency, helpers.FormatMoney(w.Balance)))
			}
		},
	})
}
