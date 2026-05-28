import { useEffect, useState } from "react";
import { apiUrl } from "./config/env";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheck,
  faCirclePlus,
  faPenToSquare,
  faTrashCan,
  faXmark,
} from "@fortawesome/free-solid-svg-icons";

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
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editingTitle, setEditingTitle] = useState("");

  const completedCount = todos.filter((todo) => todo.completed).length;
  const activeCount = todos.length - completedCount;
  const progressPercent = todos.length
    ? Math.round((completedCount / todos.length) * 100)
    : 0;

  function editTodo(todo: Todo) {
    setEditingId(todo.id);
    setEditingTitle(todo.title);
  }

  function cancelEditTodo() {
    setEditingId(null);
    setEditingTitle("");
  }

  async function saveEditTodo(id: number) {
    const trimmedTitle = editingTitle.trim();
    if (!trimmedTitle) return;

    const response = await fetch(`${apiUrl}/tasks/${id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id,
        title: trimmedTitle,
      }),
    });

    const result: ApiResponse<Todo> = await response.json();

    if (result.success) {
      setTodos((currentTodos) =>
        currentTodos.map((todo) => (todo.id === id ? result.data : todo)),
      );
      cancelEditTodo();
    }
  }

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
      setTodos((currentTodos) => [result.data, ...currentTodos]);
      setTitle("");
    }
  }

  async function handleDeleteTodo(id: number) {
    const response = await fetch(`${apiUrl}/tasks/${id}`, {
      method: "DELETE",
    });

    const result: ApiResponse<Todo> = await response.json();

    if (result.success) {
      setTodos((currentTodos) => currentTodos.filter((todo) => todo.id !== id));
    }
  }

  async function handleToggleCompleted(id: number, isChecked: boolean) {
    const response = await fetch(`${apiUrl}/tasks/${id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ completed: isChecked }),
    });

    const result: ApiResponse<Todo> = await response.json();

    if (result.success) {
      setTodos((currentTodos) =>
        currentTodos.map((todo) => (todo.id === id ? result.data : todo)),
      );
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
    <main className="app-shell">
      <section className="app-card">
        <header className="topbar">
          <div>
            <p className="eyebrow">warmdev todo</p>
            <h1>Tasks</h1>
          </div>

          <div className="auth-slot">
            {/* Later: put Login / Logout / User menu here */}
            <span>Guest mode</span>
          </div>
        </header>

        <section className="summary-row" aria-label="Todo summary">
          <div className="summary-card">
            <span>Total</span>
            <strong>{todos.length}</strong>
          </div>
          <div className="summary-card">
            <span>Active</span>
            <strong>{activeCount}</strong>
          </div>
          <div className="summary-card">
            <span>Done</span>
            <strong>{completedCount}</strong>
          </div>
          <div className="progress-card">
            <div className="progress-info">
              <span>Progress</span>
              <strong>{progressPercent}%</strong>
            </div>
            <div className="progress-track">
              <div
                className="progress-fill"
                style={{ width: `${progressPercent}%` }}
              />
            </div>
          </div>
        </section>

        <form className="composer" onSubmit={handleCreateTodo}>
          <input
            value={title}
            onChange={(event) => setTitle(event.target.value)}
            placeholder="Add a new task..."
            className="todo-title-input"
          />
          <button type="submit" className="add-todo-btn">
            <FontAwesomeIcon icon={faCirclePlus} />
            Add
          </button>
        </form>

        <section className="task-panel">
          <div className="panel-header">
            <h2>Today</h2>
            <span>{activeCount} active</span>
          </div>

          {todos.length === 0 ? (
            <div className="empty-state">
              <strong>No task yet</strong>
              <p>Create your first task and ship something small.</p>
            </div>
          ) : (
            <ul className="todo-list">
              {todos.map((todo) => (
                <li
                  className={`todo-item ${todo.completed ? "is-completed" : ""}`}
                  key={todo.id}
                >
                  <label className="check-wrap">
                    <input
                      type="checkbox"
                      checked={todo.completed}
                      onChange={(event) =>
                        handleToggleCompleted(
                          todo.id,
                          event.currentTarget.checked,
                        )
                      }
                    />
                    <span className="custom-check">
                      <FontAwesomeIcon icon={faCheck} />
                    </span>
                  </label>

                  <div className="task-content">
                    {editingId === todo.id ? (
                      <div className="edit-row">
                        <input
                          type="text"
                          value={editingTitle}
                          onChange={(event) =>
                            setEditingTitle(event.target.value)
                          }
                          className="edit-title-input"
                          autoFocus
                        />
                        <button
                          className="save-btn"
                          type="button"
                          onClick={() => saveEditTodo(todo.id)}
                          aria-label="Save todo"
                        >
                          <FontAwesomeIcon icon={faCheck} />
                        </button>
                        <button
                          className="cancel-btn"
                          type="button"
                          onClick={cancelEditTodo}
                          aria-label="Cancel edit"
                        >
                          <FontAwesomeIcon icon={faXmark} />
                        </button>
                      </div>
                    ) : (
                      <span className="todo-title">{todo.title}</span>
                    )}
                  </div>

                  <div className="action-btn">
                    <button
                      className="icon-btn"
                      type="button"
                      aria-label="Edit todo"
                      onClick={() => editTodo(todo)}
                    >
                      <FontAwesomeIcon icon={faPenToSquare} />
                    </button>
                    <button
                      className="icon-btn danger"
                      type="button"
                      aria-label="Delete todo"
                      onClick={() => handleDeleteTodo(todo.id)}
                    >
                      <FontAwesomeIcon icon={faTrashCan} />
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>
      </section>
    </main>
  );
}

export default App;
