let todoApp = document.querySelector(".todo-app");
const createTodoForm = document.querySelector("#create-todo");

async function getAllTodos() {
  try {
    const response = await fetch("http://localhost:8080/tasks");

    if (!response.ok) {
      throw new Error(`HTTP error status: ${response.status}`);
    }

    const data = await response.json();

    return data;
  } catch (error) {
    console.error("Fetch error:", error);
    return null;
  }
}

async function createNewTodo(title) {
  const response = await fetch("http://localhost:8080/tasks", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ title }),
  });

  if (!response.ok) {
    throw new Error(`HTTP error status: ${response.status}`);
  }

  const result = await response.json();
  return result;
}

async function renderTodo(todos) {
  if (!todos) return;
  const todoData = todos.data;
  todoApp.innerHTML = todoData
    .map((todo) => {
      return `
        <li>${todo.id}: ${todo.title} (${todo.completed ? "Hoàn thành" : "Chưa hoàn thành"})</li>
      `;
    })
    .join("");
}

async function main() {
  const todos = await getAllTodos();

  if (todos.success) {
    await renderTodo(todos);
  }

  createTodoForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const title = createTodoForm.title.value.trim();

    if (!title) return;

    await createNewTodo(title);

    const todos = await getAllTodos();
    await renderTodo(todos);

    createTodoForm.title.value = "";
  });
}

main();
