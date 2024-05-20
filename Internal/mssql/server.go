package mssql

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

func (c *ConnectSQLSetting) ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;",
		c.Server, c.User, c.Password, c.Database)

	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *ConnectSQLSetting) GetDatabasesOnServer() ([]DB, error) {
	conn, err := c.ConnectDB()
	if err != nil {
		return nil, err
	}

	var arrDB []DB

	txtQuery := `SELECT 
    		name, 
    		database_id, 
    		state_desc, 
    		recovery_model_desc  
		FROM 
		    Sys.Databases  
		WHERE name NOT IN ('master','model','msdb','tempdb') 
		ORDER BY name`

	rows, err := conn.Query(txtQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var db DB

		if err = rows.Scan(&db.Name, &db.ID, &db.State, &db.RecoveryModel); err != nil {
			continue
		}

		arrDB = append(arrDB, db)
	}

	return arrDB, nil
}

func (c *ConnectSQLSetting) ShrinkDatabases(db []DB, chanRes chan string) {
	conn, err := c.ConnectDB()
	if err != nil {
		return
	}

	chAll := 0
	for _, v := range db {
		if v.Mark == "" {
			continue
		}
		chAll++
	}

	ch := 0
	for _, v := range db {
		if v.Mark == "" {
			continue
		}

		ch++
		chanRes <- fmt.Sprintf("(%d/%d) | Base: %s", ch, chAll, v.Name)

		txtQuery := fmt.Sprintf(`DBCC SHRINKDATABASE ([%s], 0);`, v.Name)

		_, err = conn.Exec(txtQuery)
		if err != nil {
			//chanRes <- fmt.Sprintf("(%d/%d) | Base: %s. Err: %s", ch, chAll, v.Name, err.Error())
			continue
		}
		//chanRes <- fmt.Sprintf("(%d/%d) | Base: %s OK", ch, chAll, v.Name)
	}

	close(chanRes)

}

func (c *ConnectSQLSetting) GetConfig1C() ([]ConfigDB, error) {
	conn, err := c.ConnectDB()
	if err != nil {
		return nil, err
	}

	var arrCDB []ConfigDB

	txtQuery := `SELECT 
    		FileName
			, Creation
			, Modified
			, Attributes
			, DataSize
			, BinaryData
			, PartNo  
		FROM 
		    dbo.Config
		WHERE 
			FileName = 'versions'`

	rows, err := conn.Query(txtQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cdb ConfigDB

		if err = rows.Scan(&cdb.FileName, &cdb.Creation, &cdb.Modified, &cdb.Attributes,
			&cdb.DataSize, &cdb.BinaryData, &cdb.PartNo); err != nil {
			continue
		}

		arrCDB = append(arrCDB, cdb)
	}

	return arrCDB, nil
}
