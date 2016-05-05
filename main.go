//Copyright 2015 Jusong Chen
//
//// Author:  Jusong Chen (jusong.chen@gmail.com)

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/golang/glog"
	_ "github.com/mattn/go-oci8"
)

var (
	dsn = flag.String("logon", "", `
	<logon> is: {<username>[/<password>][@<connect_identifier>] | / } [AS {SYSDBA | SYSOPER | SYSASM}] 

	Example: Connect to database using Net Service Name and the database net service name is ORCL.
	    myusername/mypassword@ORCL

	Example: Connect to database using Easy Connect and the Service name is ORCL.

	    myusername/mypassword@localhost/ORCL
		`)

	dop = flag.Int("DOP", 4*runtime.NumCPU(), " Degree of Concurrency/Parallelism")
	//interval = flag.Duration("i", 5*time.Second, "interval between each ping")

	//port = flag.String("port", "80", "web server port number")
)

func main() {
	flag.Parse() //   SetupDB()
	setNLS_lang()

	db, err := sql.Open("oci8", *dsn)

	if db == nil || err != nil {
		log.Fatal(err)
	}
	if log.V(2) {
		log.Infof("DOP %d", *dop)
		log.Flush()
	}
	migTables(db, *dop)
	log.Flush()
}

func setNLS_lang() {

	nlsLang := os.Getenv("NLS_LANG")
	if !strings.HasSuffix(nlsLang, "UTF8") {
		i := strings.LastIndex(nlsLang, ".")
		if i < 0 {
			os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")
		} else {
			nlsLang = nlsLang[:i+1] + "AL32UTF8"
			fmt.Fprintf(os.Stderr, "NLS_LANG error: should be %s, not %s!\n",
				nlsLang, os.Getenv("NLS_LANG"))
		}
	}

}
