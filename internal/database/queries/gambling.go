package queries

const (
	InitGlobalStats         = `INSERT INTO global_stats (id, total_gambles, total_wagered, total_paid, last_update) VALUES (1, 0, 0, 0, CURRENT_TIMESTAMP);`
	GetGlobalStats          = `SELECT total_gambles, total_wagered, total_paid, last_update FROM global_stats WHERE id = 1;`
	UpdateGlobalAfterGamble = `UPDATE global_stats SET total_gambles = total_gambles + 1, total_wagered = total_wagered + ?, total_paid = total_paid + ?, last_update = CURRENT_TIMESTAMP WHERE id = 1;`
	ResetGlobalStats        = `UPDATE global_stats SET total_gambles = 0, total_wagered = 0, total_paid = 0, last_update = CURRENT_TIMESTAMP WHERE id = 1;`
)
