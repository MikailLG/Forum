package main

import (
	"footconnect/database"
)

func main() {
	contexteBDD := database.InitialiserBDD()    	
	defer contexteBDD.Close()
}
