package main

func main() {
	app := App{}
	app.Initialise(DbUser, DBPassword, DBName)
	app.Run("localhost:9999",)
}