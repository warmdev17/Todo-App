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

type User = {
  userId?: number;
  id?: number;
  username: string;
  token: string;
};

type ApiResponse<T> = {
  success: boolean;
  data: T;
  message?: string;
};

type LoginInput =
  | { username: string; password: string }
  | { email: string; password: string };

type AuthMode = "login" | "register";

function App() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editingTitle, setEditingTitle] = useState("");

  const [currentUser, setCurrentUser] = useState<string | null>(() =>
    localStorage.getItem("currentUser"),
  );
  const [authMode, setAuthMode] = useState<AuthMode | null>(null);
  const [authUsername, setAuthUsername] = useState("");
  const [authEmail, setAuthEmail] = useState("");
  const [authPassword, setAuthPassword] = useState("");
  const [authConfirmPassword, setAuthConfirmPassword] = useState("");
  const [authError, setAuthError] = useState("");

  const completedCount = todos.filter((todo) => todo.completed).length;
  const activeCount = todos.length - completedCount;
  const progressPercent = todos.length
    ? Math.round((completedCount / todos.length) * 100)
    : 0;

  function resetAuthForm() {
    setAuthUsername("");
    setAuthEmail("");
    setAuthPassword("");
    setAuthConfirmPassword("");
    setAuthError("");
  }

  function openAuthModal(mode: AuthMode) {
    resetAuthForm();
    setAuthMode(mode);
  }

  function closeAuthModal() {
    resetAuthForm();
    setAuthMode(null);
  }

  function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("currentUser");
    setCurrentUser(null);
  }

  function editTodo(todo: Todo) {
    setEditingId(todo.id);
    setEditingTitle(todo.title);
  }

  function cancelEditTodo() {
    setEditingId(null);
    setEditingTitle("");
  }

  async function getTodos() {
    const response = await fetch(`${apiUrl}/tasks`);
    const result: ApiResponse<Todo[]> = await response.json();

    if (result.success) {
      setTodos(result.data);
    }
  }

  async function saveEditTodo(id: number) {
    const trimmedTitle = editingTitle.trim();
    if (!trimmedTitle) return;

    const response = await fetch(`${apiUrl}/tasks/${id}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title: trimmedTitle }),
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

  async function handleLogin(input: LoginInput) {
    const response = await fetch(`${apiUrl}/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(input),
    });

    const result: ApiResponse<User> = await response.json();

    if (!response.ok || !result.success) {
      setAuthError(result.message ?? "Login failed");
      return;
    }

    localStorage.setItem("token", result.data.token);
    localStorage.setItem("currentUser", result.data.username);
    setCurrentUser(result.data.username);
    closeAuthModal();
  }

  async function handleRegister() {
    const username = authUsername.trim();
    const email = authEmail.trim();
    const password = authPassword.trim();
    const confirmPassword = authConfirmPassword.trim();

    if (!username || !email || !password || !confirmPassword) {
      setAuthError("Please fill all fields");
      return;
    }

    if (password !== confirmPassword) {
      setAuthError("Passwords do not match");
      return;
    }

    const response = await fetch(`${apiUrl}/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, email, password, confirmPassword }),
    });

    const result: ApiResponse<User> = await response.json();

    if (!response.ok || !result.success) {
      setAuthError(result.message ?? "Register failed");
      return;
    }

    setAuthError("");
    setAuthMode("login");
    setAuthPassword("");
    setAuthConfirmPassword("");
  }

  async function handleSubmitAuth(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();

    const emailOrUsername = authEmail.trim();
    const password = authPassword.trim();

    if (authMode === "login") {
      if (!emailOrUsername || !password) {
        setAuthError("Please enter your email/username and password");
        return;
      }

      if (emailOrUsername.includes("@")) {
        await handleLogin({ email: emailOrUsername, password });
      } else {
        await handleLogin({ username: emailOrUsername, password });
      }

      return;
    }

    await handleRegister();
  }

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
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
            {currentUser ? (
              <>
                <span className="current-user">@{currentUser}</span>
                <button type="button" onClick={handleLogout}>
                  Logout
                </button>
              </>
            ) : (
              <>
                <button type="button" onClick={() => openAuthModal("login")}>
                  Login
                </button>
                <button type="button" onClick={() => openAuthModal("register")}>
                  Register
                </button>
              </>
            )}
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

      {authMode && (
        <div className="modal-backdrop" role="presentation">
          <section className="auth-modal" aria-label="Authentication modal">
            <div className="modal-header">
              <div>
                <p className="eyebrow">account</p>
                <h2>{authMode === "login" ? "Login" : "Register"}</h2>
              </div>
              <button
                className="modal-close-btn"
                type="button"
                onClick={closeAuthModal}
                aria-label="Close modal"
              >
                <FontAwesomeIcon icon={faXmark} />
              </button>
            </div>

            <form className="auth-form" onSubmit={handleSubmitAuth}>
              {authMode === "register" && (
                <label>
                  <span>Username</span>
                  <input
                    value={authUsername}
                    onChange={(event) => setAuthUsername(event.target.value)}
                    placeholder="warmdev"
                  />
                </label>
              )}

              <label>
                <span>
                  {authMode === "login" ? "Email or username" : "Email"}
                </span>
                <input
                  value={authEmail}
                  onChange={(event) => setAuthEmail(event.target.value)}
                  placeholder={
                    authMode === "login"
                      ? "warmdev or warmdev@mail.com"
                      : "warmdev@mail.com"
                  }
                />
              </label>

              <label>
                <span>Password</span>
                <input
                  value={authPassword}
                  onChange={(event) => setAuthPassword(event.target.value)}
                  placeholder="••••••••"
                  type="password"
                />
              </label>

              {authMode === "register" && (
                <label>
                  <span>Confirm password</span>
                  <input
                    value={authConfirmPassword}
                    onChange={(event) =>
                      setAuthConfirmPassword(event.target.value)
                    }
                    placeholder="••••••••"
                    type="password"
                  />
                </label>
              )}

              {authError && <p className="auth-error">{authError}</p>}

              <button className="modal-submit" type="submit">
                {authMode === "login" ? "Login" : "Create account"}
              </button>
            </form>

            <button
              className="auth-switch-btn"
              type="button"
              onClick={() => {
                resetAuthForm();
                setAuthMode(authMode === "login" ? "register" : "login");
              }}
            >
              {authMode === "login"
                ? "Need an account? Register"
                : "Already have an account? Login"}
            </button>
          </section>
        </div>
      )}
    </main>
  );
}

export default App;
