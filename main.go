package main

import "log"

func main() {
	db, err := NewMysqlStorage()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully Connected to mysql database!")
	api := NewApiServer(":"+Envs.PORT, db)
	api.Run()
}
