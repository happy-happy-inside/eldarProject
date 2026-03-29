const API = ""

// ===== COURSE =====

window.createCourse = async function () {
    const title = document.getElementById("courseTitle").value
    const description = document.getElementById("courseDesc").value

    await fetch("/courses", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ title, description })
    })

    alert("Course created")
}

// ===== LESSON =====

window.createLesson = async function () {
    const course_id = parseInt(document.getElementById("lessonCourseId").value)
    const title = document.getElementById("lessonTitle").value
    const content = document.getElementById("lessonContent").value

    await fetch("/lessons", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ course_id, title, content, position: 1 })
    })

    alert("Lesson created")
}

// ===== TEST BUILDER =====

let questionCount = 0

window.addQuestion = function () {
    questionCount++

    const div = document.getElementById("questions")

    div.innerHTML += `
        <div id="q${questionCount}">
            <input placeholder="Question text" id="qtext${questionCount}">

            <div>
                <input placeholder="Answer 1" id="a${questionCount}_1">
                <input type="radio" name="correct${questionCount}" value="1"> correct
            </div>

            <div>
                <input placeholder="Answer 2" id="a${questionCount}_2">
                <input type="radio" name="correct${questionCount}" value="2"> correct
            </div>

            <div>
                <input placeholder="Answer 3" id="a${questionCount}_3">
                <input type="radio" name="correct${questionCount}" value="3"> correct
            </div>
        </div>
        <hr>
    `
}

window.createTest = async function () {
    const lessonId = document.getElementById("testLessonId").value

    const questions = []

    for (let i = 1; i <= questionCount; i++) {
        const text = document.getElementById(`qtext${i}`).value
        const correct = document.querySelector(`input[name="correct${i}"]:checked`)?.value

        const answers = []

        for (let j = 1; j <= 3; j++) {
            const val = document.getElementById(`a${i}_${j}`).value

            answers.push({
                answer: val,
                is_correct: correct == j
            })
        }

        questions.push({ text, answers })
    }

    await fetch(`/tests/${lessonId}`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ questions })
    })

    alert("Test created")
}