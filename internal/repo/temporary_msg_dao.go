package repo

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TmpMsg struct {
	Id     int64
	ChatId int64
	MsgId  int64
}

type TmpMsgDao struct {
	db *pgxpool.Pool
}

func NewTmpMsgDao(db *pgxpool.Pool) TmpMsgDao {
	return TmpMsgDao{db: db}
}

func (r TmpMsgDao) GetAll(ctx context.Context) ([]TmpMsg, error) {
	query := `SELECT t.id, t.chatId, t.msgId FROM temporary_messages t`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]TmpMsg, 0)
	for rows.Next() {
		var tm TmpMsg
		if err := rows.Scan(&tm.Id, &tm.ChatId, &tm.MsgId); err != nil {
			slog.Warn(err.Error())
		}
		items = append(items, tm)
	}
	return items, nil
}

func (r TmpMsgDao) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM temporary_messages t WHERE t.id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r TmpMsgDao) Save(ctx context.Context, chatId int64, msgId int64) error {
	query := `INSERT INTO temporary_messages (chatId, msgId) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, chatId, msgId)
	return err
}
