package connector

import (
	. "context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Connect struct {
	name  ConnectorName
	getDB dbGetter
}

func (p *Connect) Invoke(ctx Context, handle interface{}) error {

	var (
		isDB     bool
		isTx     bool
		isQ      bool
		handleDB func(*sqlx.DB) error
		handleTx func(*sqlx.Tx) error
		handleQ  func(Q) error
	)
	switch v := handle.(type) {
	case func(*sqlx.DB) error:
		isDB = true
		handleDB = v
	case func(*sqlx.Tx) error:
		isTx = true
		handleTx = v
	case func(Q) error:
		isQ = true
		handleQ = v
	}

	switch {
	case isDB:
		return handleDB(p.getDB())
	case isTx:
		return p.invokeTx(ctx, handleTx)
	case isQ:
		return handleQ(p.getDB())
	default:
		return fmt.Errorf("not support db handler %T", handle)
	}
}

func (p *Connect) invokeTx(ctx Context, handle handleTX) (err error) {
	// Begin a transcation with context.
	tx, err := p.getDB().BeginTxx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if re := recover(); re != nil {
			switch re := re.(type) {
			case error:
				err = re
			default:
				err = fmt.Errorf("%s", re)
			}
		}

		// Check Error & Rollback
		if err != nil {
			tx.Rollback()
			return
		}

		// Commit
		err = tx.Commit()
		return
	}()

	err = handle(tx)
	return
}
