import { useEffect, useState } from "react";
import { apiUrl } from "./config/env";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrashCan } from "@fortawesome/free-solid-svg-icons";
import { library } from "@fortawesome/fontawesome-svg-core";

/* import all the icons in Free Solid, Free Regular, and Brands styles */
import { fas } from "@fortawesome/free-solid-svg-icons";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fab } from "@fortawesome/free-brands-svg-icons";

import "./App.css";

type Todo = {
  id: number;
  title: string;
  completed: boolean;
};

type ApiResponse<T> = {
  success: boolean;
  data: T;
};

function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState("");

  library.add(fas, far, fab);

  async function handleCreateTodo(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedTitle = title.trim();
    if (!trimmedTitle) return;

    const response = await fetch(`${apiUrl}/tasks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title: trimmedTitle }),
    });

    const result: ApiResponse<Todo> = await response.json();

    if (result.success) {
      setTodos((currentTodos) => [...currentTodos, result.data]);
      setTitle("");
    }
  }

  useEffect(() => {
    async function getTodos() {
      const response = await fetch(`${apiUrl}/tasks`);
      const result: ApiResponse<Todo[]> = await response.json();

      if (result.success) {
        setTodos(result.data);
      }
    }

    getTodos();
  }, []);

  return (
    <main className="container">
      <h1>Todo App</h1>

      <ul className="todo-list">
        {todos.map((todo) => (
          <li className="todo-item" key={todo.id}>
            <div className="todo-item-content">
              <input type="checkbox" checked={todo.completed} />
              <span>{todo.title}</span>
            </div>
            <FontAwesomeIcon icon={faTrashCan} />
          </li>
        ))}
      </ul>

      <form id="create-todo-form" onSubmit={handleCreateTodo}>
        <input
          value={title}
          onChange={(event) => setTitle(event.target.value)}
          placeholder="Enter todo..."
          className="todo-title-input"
        />
        <button type="submit" className="add-todo-btn">
          Add
        </button>
      </form>
    </main>
  );
}

export default App;
