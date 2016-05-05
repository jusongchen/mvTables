// Copyright 2015 Jusong Chen
//
// Author: Jusong Chen (jusong.chen@gmail.com)

package main

import (
	"database/sql"
	"fmt"

	log "github.com/golang/glog"
)

//a row is represented as a slice of columns

type oraTable struct {
	TableOwner        string
	TableName         string
	PartitionName     string
	SubpartitionCount int
}

//Done() is executed as last step in parallel execution
func (o *oraTable) Done() {
	//log.Fatalf("%s %s.%s processed", o.TypeDesc, o.SchemaName, o.ObjName)
	//log.Fatalf("row with SchemaName %s processed", o.SchemaName)
	log.Infof("Processed:%+v", *o)
}

//Process()
//TODO:create cockroach table
func (o *oraTable) Process() {
	log.Infof("Processing %+v", *o)
}

//migTables get table definitions from source DB and for each table:
//		1.create destination table in cockroach
//		2. copy data from source to dest
//one table is processed in one goroutine while multiple tables can be processed concurrently

func migTables(db *sql.DB, DOP int) error {

	rs, err := db.Query(
		`select table_owner,table_name,partition_name,subpartition_count 
			from dba_tab_partitions
			where table_owner not in ('SYS')
			order by table_owner, table_name
			`)

	if err != nil {
		log.Fatal(err)
		return err
	}
	// close Query
	defer func() {
		rs.Close()
	}()

	o := oraTable{}
	//execute in sequencial order
	for rs.Next() {
		rs.Scan(&o.TableOwner, &o.TableName, &o.PartitionName, &o.SubpartitionCount)
		fmt.Printf("%#v\n", o)
	}

	return nil
}
