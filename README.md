Markdown
# 📚 Education Platform (Go + PostgreSQL)

Веб-приложение для создания и прохождения онлайн-курсов с уроками и тестами.

Проект включает:
- Backend на Go
- PostgreSQL базу данных
- Frontend (HTML/CSS/JS)
- Админ-панель для управления контентом
- Docker + Docker Compose для запуска

---

# 🚀 Функционал

## 👨‍🎓 Пользователь
- Просмотр списка курсов
- Просмотр уроков внутри курса
- Чтение контента урока
- Прохождение тестов
- Получение результата

## 🛠 Админка
- Создание курсов
- Добавление уроков
- Создание тестов с вопросами и ответами

---

# 🧱 Архитектура
. ├── main.go              # Backend (Go) ├── Dockerfile ├── docker-compose.yml ├── /templates │   ├── index.html │   ├── course.html │   ├── lesson.html │   └── admin.html ├── /static │   ├── app.js │   ├── admin.js │   └── style.css

---

# 🐳 Запуск проекта

## 1. Требования

- Docker
- Docker Compose

---

## 2. Запуск

`bash
docker-compose up --build
3. Доступ
Главная: http://localhost:8080⁠�
Админка: http://localhost:8080/admin⁠�
🗄 База данных
Используется PostgreSQL.
Таблицы:
courses
поле
тип
id
SERIAL
title
TEXT
description
TEXT
lessons
поле
тип
id
SERIAL
course_id
INT
title
TEXT
content
TEXT
position
INT
tests
поле
тип
id
SERIAL
lesson_id
INT
questions
поле
тип
id
SERIAL
test_id
INT
text
TEXT
answers
поле
тип
id
SERIAL
question_id
INT
answer
TEXT
is_correct
BOOLEAN
🔌 API
📘 Курсы
Получить все курсы

GET /courses
Получить курс

GET /courses/{id}
Создать курс

POST /courses
JSON
{
  "title": "Go курс",
  "description": "Изучение Go"
}
📗 Уроки
Получить все уроки

GET /lessons
Получить урок

GET /lessons/{id}
Создать урок

POST /lessons
JSON
{
  "course_id": 1,
  "title": "Урок 1",
  "content": "Текст",
  "position": 1
}
🧠 Тесты
Получить тест по уроку

GET /tests/{lesson_id}
Создать тест

POST /tests/{lesson_id}
JSON
{
  "questions": [
    {
      "text": "Вопрос?",
      "answers": [
        { "answer": "Да", "is_correct": true },
        { "answer": "Нет", "is_correct": false }
      ]
    }
  ]
}
✅ Отправка теста

POST /submit/{test_id}
JSON
{
  "1": 1,
  "2": 5
}
Ответ:
JSON
{
  "score": 80
}
🎨 Frontend
Страницы:
/ — список курсов
/course?id=1 — уроки курса
/lesson?id=1 — урок + тест
/admin — админка
⚠️ Возможные проблемы
❌ Уроки не отображаются
Проверь:
course_id совпадает
фильтр:
JavaScript
l.course_id === Number(id)
❌ JS ошибки (not defined)
Причина:
кэш браузера
Решение:

Ctrl + Shift + R
❌ admin.js не работает
Проверь:
HTML
<script defer src="/static/admin.js"></script>
🔧 Улучшения (TODO)
Авторизация пользователей
Сохранение прогресса
Редактирование курсов
Удаление данных
UI/UX улучшения
Пагинация
Загрузка изображений
🧠 Технологии
Go (net/http)
PostgreSQL
HTML / CSS / JavaScript
Docker
📌 Итог
Проект представляет собой полноценную образовательную платформу с:
backend API
базой данных
пользовательским интерфейсом
админ-панелью
Подходит как учебный проект или база для дальнейшего развития.
