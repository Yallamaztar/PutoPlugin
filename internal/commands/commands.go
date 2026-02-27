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
	"plugin/internal/service/bank"
	"plugin/internal/service/link"
	"plugin/internal/service/player"
	"plugin/internal/service/stats"
	"plugin/internal/service/wallet"
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

func RegisterClientCommands(
	cfg config.Config,

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
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			id, err := player.GetDiscordIDByID(playerID)
			if err != nil {
				rcon.Tell(clientNum, "^1Error ^7checking your account, try again later")
				return
			}

			if id != "" {
				rcon.Tell(clientNum, "You have ^6already linked ^7your account")
				return
			}

			code, err := generateCode()
			if err != nil {
				rcon.Tell(clientNum, "^1Failed ^7to generate link code, ^6try again ^7later")
				return
			}

			if err = link.CreateLink(playerID, code); err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			rcon.Tell(clientNum, fmt.Sprintf("Your discord link code is: ^6%s", code))
			rcon.Tell(clientNum, "use ^6/link <code> ^7in discord to link your account")
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "gamble",
		Aliases:  []string{"g", "cf", "coinflip"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!gamble ^7<amount>",
		MinArgs:  1,
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			balance, err := wallet.GetBalance(playerID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			amount, err := helpers.ParseAmountArg(args[0], int64(balance))
			if err != nil {
				rcon.Tell(clientNum, fmt.Sprintf("%s (%q)", err, args[0]))
				return
			}

			res, err := gamble.Gamble(playerID, playerName, amount, cfg, player, wallet, bank, playerStats, gambleStats, walletStats, webhook)
			if err != nil {
				rcon.Tell(clientNum, err.Error())
				return
			}

			rcon.Tell(clientNum, res.Message)
			if res.Won {
				rcon.Say(fmt.Sprintf("%s just ^6won ^7%s%d!", playerName, cfg.Gambling.Currency, res.Amount))
			} else {
				rcon.Say(fmt.Sprintf("%s just ^6lost ^7%s%d!", playerName, cfg.Gambling.Currency, res.Amount))
			}

		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "pay",
		Aliases:  []string{"pp", "payplayer", "transfer"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!pay <player> <amount>",
		MinArgs:  2,
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			balance, err := wallet.GetBalance(playerID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			amount, err := helpers.ParseAmountArg(args[0], int64(balance))
			if err != nil {
				rcon.Tell(clientNum, fmt.Sprintf("%s (%q)", err, args[0]))
				return
			}

			t := reg.FindPlayer(args[0])
			if t == nil {
				rcon.Tell(clientNum, fmt.Sprintf("player %s couldnt be found", args[0]))
				return
			}

			target, err := player.GetPlayerByGUID(t.GUID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			res, err := pay.Pay(playerID, target.ID, amount, cfg, player, wallet, walletStats, webhook)
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
		MinArgs:  1,
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			bal, err := wallet.GetBalance(playerID)
			if err != nil {
				rcon.Tell(clientNum, unknownErr)
				return
			}

			rcon.Tell(clientNum, fmt.Sprintf("Your balance is %s%d", cfg.Gambling.Currency, bal))
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "bankbalance",
		Aliases:  []string{"bb", "bank", "bankbal"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!bankbalance",
		MinArgs:  1,
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
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
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {
			rcon.Tell(clientNum, cfg.Discord.InviteLink)
		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "richest",
		Aliases:  []string{"rich"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!richest",
		MinArgs:  1,
		Handler: func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {

		},
	})

	reg.RegisterCommand(register.Command{
		Name:     "poorest",
		Aliases:  []string{"poor"},
		MinLevel: LevelUser,
		Help:     "Usage: ^6!poorest",
		MinArgs:  1,
		Handler:  func(clientNum uint8, playerID int, playerName, xuid string, level int, args []string) {},
	})

}
