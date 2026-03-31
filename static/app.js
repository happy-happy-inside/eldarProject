const API = ""

function toggleTheme() {
    document.body.classList.toggle("dark")

    const isDark = document.body.classList.contains("dark")
    localStorage.setItem("theme", isDark ? "dark" : "light")
}

function loadTheme() {
    const saved = localStorage.getItem("theme")
    if (saved === "dark") {
        document.body.classList.add("dark")
    }
}


// ===== ROUTING =====

window.onload = () => {
    loadTheme()
    const path = window.location.pathname
    
    if (path === "/") loadCourses()
        if (path.startsWith("/course")) loadCourse()
            if (path.startsWith("/lesson")) loadLesson()
            }
        
// ===== COURSES =====

let allCourses = []

async function loadCourses() {

    const res = await fetch("/courses")
    const data = await res.json()
            
    allCourses = data
    currentPage = 1

    renderPage()
}


function renderCourses(courses) {
    const div = document.getElementById("courses")
    div.innerHTML = ""

    courses.forEach(c => {
        div.innerHTML += `
            <div class="card">
                <h3>📘 ${c.title}</h3>
                <p>${c.description}</p>
                <a href="/course?id=${c.id}">Open →</a>
            </div>
        `
    })
}

function filterCourses() {
    const query = document.getElementById("searchInput").value.toLowerCase()

    const filtered = allCourses.filter(c =>
        c.title.toLowerCase().includes(query) ||
        c.description.toLowerCase().includes(query)
    )

    renderCourses(filtered)
}

// ===== PAGE =====

let currentPage = 1
const perPage = 5

function renderPage() {
    const start = (currentPage - 1) * perPage
    const end = start + perPage

    const pageCourses = allCourses.slice(start, end)

    renderCourses(pageCourses)
    renderPagination()
}

function renderPagination() {
    const totalPages = Math.ceil(allCourses.length / perPage)
    const div = document.getElementById("pagination")

    div.innerHTML = ""

    if (currentPage > 1) {
        div.innerHTML += `<button onclick="changePage(${currentPage - 1})">←</button>`
    }

    for (let i = 1; i <= totalPages; i++) {
        div.innerHTML += `
            <button onclick="changePage(${i})">
                ${i}
            </button>
        `
    }

    if (currentPage < totalPages) {
        div.innerHTML += `<button onclick="changePage(${currentPage + 1})">→</button>`
    }
}

function changePage(page) {
    currentPage = page
    renderPage()
}

// ===== LESSON + TEST =====

let currentTest = []

async function loadLesson() {
    const id = new URLSearchParams(window.location.search).get("id")

    const res = await fetch("/lessons/" + id)
    const lesson = await res.json()

    document.getElementById("lessonTitle").innerText = lesson.title
    document.getElementById("content").innerText = lesson.content

    // тест
    const testRes = await fetch("/tests/" + id)
    const test = await testRes.json()

    currentTest = test

    const div = document.getElementById("test")
    div.innerHTML = ""

    test.forEach(q => {
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

    const res = await fetch("/submit/" + testId, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(answers)
    })

    const result = await res.json()

    document.getElementById("result").innerText = "Score: " + result.score + "%"
}