import { useState, useEffect } from "react";
import "./App.css";
import TodoInput from "./components/TodoInput";
import TodoList from "./components/TodoList";
import { GetTodos, AddTodo, ToggleTodo, DeleteTodo } from "../wailsjs/go/backend/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

interface Todo {
  id: string;
  title: string;
  done: boolean;
  createdAt: string;
}

function App() {
  const [todos, setTodos] = useState<Todo[]>([]);

  const loadTodos = async () => {
    const result = await GetTodos();
    setTodos(result || []);
  };

  useEffect(() => {
    loadTodos();
    EventsOn("todos:reload", loadTodos);
  }, []);

  const handleAdd = async (title: string) => {
    await AddTodo(title);
    loadTodos();
  };

  const handleToggle = async (id: string) => {
    const updated = await ToggleTodo(id);
    setTodos(updated || []);
  };

  const handleDelete = async (id: string) => {
    const updated = await DeleteTodo(id);
    setTodos(updated || []);
  };

  const doneCount = todos.filter((t) => t.done).length;

  return (
    <div id="App">
      <h1>Wails Todo</h1>
      <p className="status">
        {todos.length}개 중 {doneCount}개 완료
      </p>
      <TodoInput onAdd={handleAdd} />
      <TodoList todos={todos} onToggle={handleToggle} onDelete={handleDelete} />
    </div>
  );
}

export default App;
