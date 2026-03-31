package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// ===== MODELS =====

type Course struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Lesson struct {
	ID       int    `json:"id"`
	CourseID int    `json:"course_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Position int    `json:"position"`
	Image    string `json:"image"`
}

type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	IsCorrect  bool   `json:"is_correct"`
}

// ===== MAIN =====

func main() {
	connStr := "host=db user=postgres password=postgres dbname=education sslmode=disable"

	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		fmt.Println("Waiting for DB...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal(err)
	}

	createTables()

	// ===== API =====
	http.HandleFunc("/courses", coursesHandler)
	http.HandleFunc("/courses/", courseHandler)

	http.HandleFunc("/lessons", lessonsHandler)
	http.HandleFunc("/lessons/", lessonHandler)

	http.HandleFunc("/tests/", testHandler)
	http.HandleFunc("/submit/", submitTestHandler)

	// ===== STATIC =====
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/admin.html")
	})

	// ===== PAGES =====
	http.HandleFunc("/course", coursePage)
	http.HandleFunc("/lesson", lessonPage)
	http.HandleFunc("/", indexPage)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

// ===== DATABASE =====

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS courses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS lessons (
		id SERIAL PRIMARY KEY,
		course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
		title TEXT,
		content TEXT,
		position INTEGERl,
		image TEXT
	);

	CREATE TABLE IF NOT EXISTS tests (
		id SERIAL PRIMARY KEY,
		lesson_id INTEGER UNIQUE REFERENCES lessons(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		test_id INTEGER REFERENCES tests(id) ON DELETE CASCADE,
		question TEXT
	);

	CREATE TABLE IF NOT EXISTS answers (
		id SERIAL PRIMARY KEY,
		question_id INTEGER REFERENCES questions(id) ON DELETE CASCADE,
		answer TEXT,
		is_correct BOOLEAN
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

// ===== COURSES =====

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		rows, _ := db.Query("SELECT id,title,description FROM courses")
		defer rows.Close()

		list := []Course{}

		for rows.Next() {
			var c Course
			rows.Scan(&c.ID, &c.Title, &c.Description)
			list = append(list, c)
		}

		json.NewEncoder(w).Encode(list)

	case "POST":
		var c Course
		json.NewDecoder(r.Body).Decode(&c)

		db.QueryRow(
			"INSERT INTO courses (title,description) VALUES ($1,$2) RETURNING id",
			c.Title, c.Description,
		).Scan(&c.ID)

		json.NewEncoder(w).Encode(c)
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/courses/")

	switch r.Method {

	case "GET":
		var c Course
		err := db.QueryRow("SELECT id,title,description FROM courses WHERE id=$1", id).
			Scan(&c.ID, &c.Title, &c.Description)

		if err != nil {
			http.Error(w, "not found", 404)
			return
		}

		json.NewEncoder(w).Encode(c)

	case "PUT":
		var c Course
		json.NewDecoder(r.Body).Decode(&c)

		db.Exec("UPDATE courses SET title=$1,description=$2 WHERE id=$3",
			c.Title, c.Description, id)

	case "DELETE":
		db.Exec("DELETE FROM courses WHERE id=$1", id)
	}
}

// ===== LESSONS =====

func lessonsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		rows, _ := db.Query("SELECT id,course_id,title,content,position FROM lessons")
		defer rows.Close()

		list := []Lesson{}

		for rows.Next() {
			var l Lesson
			rows.Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.Position)
			list = append(list, l)
		}

		json.NewEncoder(w).Encode(list)

	case "POST":
		var l Lesson
		json.NewDecoder(r.Body).Decode(&l)

		db.QueryRow(
			"INSERT INTO lessons (course_id, title, content, image, position) VALUES ($1,$2,$3,$4,$5) RETURNING id",
			l.CourseID, l.Title, l.Content, l.Image, l.Position,
		).Scan(&l.ID)
		json.NewEncoder(w).Encode(l)
	}
}

func lessonHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/lessons/")

	switch r.Method {

	case "GET":
		var l Lesson
		err := db.QueryRow(
			"SELECT id,course_id,title,content,position FROM lessons WHERE id=$1", id,
		).Scan(&l.ID, &l.CourseID, &l.Title, &l.Content, &l.Position)

		if err != nil {
			http.Error(w, "not found", 404)
			return
		}

		json.NewEncoder(w).Encode(l)

	case "DELETE":
		db.Exec("DELETE FROM lessons WHERE id=$1", id)
	}
}

// ===== TESTS =====

func testHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tests/")

	switch r.Method {

	case "POST":
		lessonID, _ := strconv.Atoi(id)

		var payload struct {
			Questions []struct {
				Text    string   `json:"text"`
				Answers []Answer `json:"answers"`
			} `json:"questions"`
		}

		json.NewDecoder(r.Body).Decode(&payload)

		var testID int
		db.QueryRow(
			"INSERT INTO tests (lesson_id) VALUES ($1) RETURNING id",
			lessonID,
		).Scan(&testID)

		for _, q := range payload.Questions {
			var qID int

			db.QueryRow(
				"INSERT INTO questions (test_id,question) VALUES ($1,$2) RETURNING id",
				testID, q.Text,
			).Scan(&qID)

			for _, a := range q.Answers {
				db.Exec(
					"INSERT INTO answers (question_id,answer,is_correct) VALUES ($1,$2,$3)",
					qID, a.Answer, a.IsCorrect,
				)
			}
		}

		w.Write([]byte("test created"))

	case "GET":
		lessonID, _ := strconv.Atoi(id)

		var testID int
		db.QueryRow("SELECT id FROM tests WHERE lesson_id=$1", lessonID).Scan(&testID)

		rows, _ := db.Query("SELECT id,question FROM questions WHERE test_id=$1", testID)
		defer rows.Close()

		type FullQuestion struct {
			ID      int      `json:"id"`
			Text    string   `json:"text"`
			Answers []Answer `json:"answers"`
		}

		var result []FullQuestion

		for rows.Next() {
			var fq FullQuestion
			rows.Scan(&fq.ID, &fq.Text)

			ansRows, _ := db.Query("SELECT id,answer,is_correct FROM answers WHERE question_id=$1", fq.ID)

			for ansRows.Next() {
				var a Answer
				ansRows.Scan(&a.ID, &a.Answer, &a.IsCorrect)
				fq.Answers = append(fq.Answers, a)
			}
			ansRows.Close()

			result = append(result, fq)
		}

		json.NewEncoder(w).Encode(result)
	}
}

// ===== SUBMIT TEST =====

func submitTestHandler(w http.ResponseWriter, r *http.Request) {
	testIDStr := strings.TrimPrefix(r.URL.Path, "/submit/")
	testID, _ := strconv.Atoi(testIDStr)

	var answers map[string]int
	json.NewDecoder(r.Body).Decode(&answers)

	rows, _ := db.Query(`
		SELECT a.id, a.is_correct
		FROM answers a
		JOIN questions q ON a.question_id = q.id
		WHERE q.test_id=$1`, testID)

	defer rows.Close()

	correct := 0
	total := 0

	for rows.Next() {
		var id int
		var isCorrect bool
		rows.Scan(&id, &isCorrect)

		total++

		if selected, ok := answers[strconv.Itoa(id)]; ok {
			if selected == id && isCorrect {
				correct++
			}
		}
	}

	score := 0
	if total > 0 {
		score = (correct * 100) / total
	}

	json.NewEncoder(w).Encode(map[string]int{
		"score": score,
	})
}

// ===== PAGES =====

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func coursePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/course.html")
}

func lessonPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/lesson.html")
}
