package main


func main() {
	app := NewApp()
	defer app.Close()

	app.Run()
}
