import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import TerminalPage from './pages/TerminalPage';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/terminal/:robotId" element={<TerminalPage />} />
      </Routes>
    </BrowserRouter>
  );
}
