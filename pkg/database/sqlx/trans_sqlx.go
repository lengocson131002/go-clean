package database

import "github.com/lengocson131002/go-clean/pkg/database"

//==============================SQLX DB Transaction==================
// DB doesn't rollback, do nothing here
func (cdt *SqlxDBTx) Rollback() error {
	return nil
}

//DB doesnt commit, do nothing here
func (cdt *SqlxDBTx) Commit() error {
	return nil
}

// TransactionBegin starts a transaction
func (sdt *SqlxDBTx) TxBegin() (database.SqlGdbc, error) {
	tx, err := sdt.DB.Beginx()
	sct := SqlxConnTx{tx}
	return &sct, err
}

// DB doesnt rollback, do nothing here
func (cdt *SqlxDBTx) TxEnd(txFunc func() error) error {
	return nil
}

//===================================================================

//========================SQLX CONNECTION TRANSACTION================
//*sql.Tx can't begin a transaction, transaction always begins with a *sql.DB
func (sct *SqlxConnTx) TxBegin() (database.SqlGdbc, error) {
	return nil, nil
}

func (sct *SqlxConnTx) Rollback() error {
	return sct.DB.Rollback()
}

func (sct *SqlxConnTx) Commit() error {
	return sct.DB.Commit()
}

func (sct *SqlxConnTx) TxEnd(txFunc func() error) error {
	var err error
	tx := sct.DB

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()
	err = txFunc()
	return err
}

//================================================================
