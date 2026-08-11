package main

import "fmt"

func main() {
	client, err := NewErrorsClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = runLessonSync(client.Capture, "lesson-sync", syncLessons)
	if err != nil {
		fmt.Println("scheduled job failure surfaced")
		return
	}
	fmt.Println("scheduled job completed")
}
