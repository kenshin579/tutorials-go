import { useState } from "react";

interface TodoInputProps {
  onAdd: (title: string) => void;
}

function TodoInput({ onAdd }: TodoInputProps) {
  const [title, setTitle] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (title.trim()) {
      onAdd(title.trim());
      setTitle("");
    }
  };

  return (
    <form className="todo-input" onSubmit={handleSubmit}>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="할 일을 입력하세요..."
        autoFocus
      />
      <button type="submit">추가</button>
    </form>
  );
}

export default TodoInput;
