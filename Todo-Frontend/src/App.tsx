import { useEffect, useState } from "react";
import { apiUrl } from "./config/env";
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
    <main>
      <h1>Todo App</h1>

      <ul>
        {todos.map((todo) => (
          <li key={todo.id}>
            {todo.id}: {todo.title} -{" "}
            {todo.completed ? "Hoàn thành" : "Chưa hoàn thành"}
          </li>
        ))}
      </ul>
    </main>
  );
}

export default App;
