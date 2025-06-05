package main

func InsertTask(title string) error {
	_, err := db.Exec("INSERT INTO tasks (title) VALUES ($1)", title)
	return err
}
