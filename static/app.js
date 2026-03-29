const API = "http://localhost:8080"

// ===== ROUTING =====

window.onload = () => {
    const path = window.location.pathname

    if (path === "/") loadCourses()
    if (path.startsWith("/course")) loadCourse()
    if (path.startsWith("/lesson")) loadLesson()
}

// ===== COURSES =====

async function loadCourses() {
    const res = await fetch(API + "/courses")
    const data = await res.json()

    const div = document.getElementById("courses")

    data.forEach(c => {
        div.innerHTML += `
            <div>
                <h3>${c.title}</h3>
                <p>${c.description}</p>
                <a href="/course?id=${c.id}">Open</a>
            </div>
        `
    })
}

// ===== COURSE =====

async function loadCourse() {
    const id = new URLSearchParams(window.location.search).get("id")

    const res = await fetch(API + "/lessons")
    const lessons = await res.json()

    const div = document.getElementById("lessons")

    lessons
        .filter(l => l.course_id == id)
        .forEach(l => {
       div.innerHTML += `
    <div class="card">
        <h3>${c.title}</h3>
        <p>${c.description}</p>
        <a href="/course?id=${c.id}">Open course →</a>
    </div>
`
        })
}

// ===== LESSON + TEST =====

let currentTest = []

async function loadLesson() {
    const id = new URLSearchParams(window.location.search).get("id")

    const res = await fetch(API + "/lessons/" + id)
    const lesson = await res.json()

    document.getElementById("lessonTitle").innerText = lesson.title
    document.getElementById("content").innerText = lesson.content

    // загрузка теста
    const testRes = await fetch(API + "/tests/" + id)
    const test = await testRes.json()

    currentTest = test

    const div = document.getElementById("test")

    test.forEach(q => {
        div.innerHTML += `<p><b>${q.text}</b></p>`

        q.answers.forEach(a => {
          div.innerHTML += `
    <div class="card">
        <p><b>${q.text}</b></p>

        ${q.answers.map(a => `
            <label>
                <input type="radio" name="q${q.id}" value="${a.id}">
                ${a.answer}
            </label><br>
        `).join("")}
    </div>
`
        })
    })
}

// ===== SUBMIT =====

async function submitTest() {
    const answers = {}

    currentTest.forEach(q => {
        const selected = document.querySelector(`input[name="q${q.id}"]:checked`)
        if (selected) {
            answers[selected.value] = parseInt(selected.value)
        }
    })

    const testId = currentTest[0]?.test_id || 1

    const res = await fetch(API + "/submit/" + testId, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(answers)
    })

    const result = await res.json()

    document.getElementById("result").innerText = "Score: " + result.score + "%"
}