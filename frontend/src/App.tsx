import { Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home'
import CadastrarCliente from './pages/CadastrarCliente'
import RegistrarCompra from './pages/RegistrarCompra'
import HistoricoCompras from './pages/HistoricoCompras'
import Dashboard from './pages/Dashboard'
import AlertasDashboard from './pages/AlertasDashboard'


function App() {
  return (
    <div className="max-w-5xl mx-auto py-10 px-4">
      <h1 className="text-3xl font-bold mb-6">Smart Retention</h1>

      <nav className="mb-6 space-x-4">
        <Link to="/" className="text-blue-500 hover:underline">🏠 Home</Link>
        <Link to="/cadastrar" className="text-blue-500 hover:underline">➕ Cadastrar Cliente</Link>
        <Link to="/compras" className="text-blue-500 hover:underline">🧾 Registrar Compra</Link>
        <Link to="/historico" className="text-blue-500 hover:underline">📜 Histórico de Compras</Link>
        <Link to="/dashboard" className="text-blue-500 hover:underline">📊 Dashboard</Link>
        <Link to="/alertas" className="text-blue-500 hover:underline">🔔 Alertas</Link>
      </nav>

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/cadastrar" element={<CadastrarCliente />} />
        <Route path="/compras" element={<RegistrarCompra />} />
        <Route path="/historico" element={<HistoricoCompras />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/alertas" element={<AlertasDashboard />} />
      </Routes>
    </div>
  )
}

export default App
