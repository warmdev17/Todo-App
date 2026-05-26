let todoApp = document.querySelector(".todo-app");

console.log(todoApp);

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

async function main() {
  const todos = await getAllTodos();

  if (todos.success) {
    const todoData = todos.data;
    todoApp.innerHTML = todoData
      .map((todo) => {
        return `
        <li>${todo.id}: ${todo.title} (${todo.completed ? "Hoàn thành" : "Chưa hoàn thành"})</li>
      `;
      })
      .join("");
  }
}

main();
