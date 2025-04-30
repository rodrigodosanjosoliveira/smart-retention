import { useEffect, useState } from 'react'
import axios from 'axios'
import {
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer,
  BarChart, Bar, PieChart, Pie, Cell, Legend
} from 'recharts'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

const COLORS = ['#8884d8', '#82ca9d', '#ffc658', '#ff7f50', '#ffbb28']

interface Dashboard {
  total_clientes: number
  total_compras: number
  compras_por_mes: { mes: string, quantidade: number }[]
  itens_mais_comprados: { nome: string, quantidade: number }[]
  clientes_mais_ativos: { nome: string, quantidade: number }[]
}

export default function Dashboard() {
  const [data, setData] = useState<Dashboard | null>(null)

  useEffect(() => {
    api.get('/api/dashboard').then((res) => setData(res.data))
  }, [])

  if (!data) return <p className="text-center mt-10">Carregando dashboard...</p>

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      <h2 className="text-2xl font-semibold mb-4">📊 Dashboard</h2>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="bg-white p-4 rounded shadow text-center">
          <p className="text-gray-500 text-sm">Total de Clientes</p>
          <p className="text-xl font-bold">{data.total_clientes}</p>
        </div>
        <div className="bg-white p-4 rounded shadow text-center">
          <p className="text-gray-500 text-sm">Total de Compras</p>
          <p className="text-xl font-bold">{data.total_compras}</p>
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-8">
        <div>
          <h3 className="text-lg font-semibold mb-2">📆 Compras por Mês</h3>
          <ResponsiveContainer width="100%" height={250}>
            <LineChart data={data.compras_por_mes}>
              <XAxis dataKey="mes" />
              <YAxis />
              <Tooltip />
              <Line type="monotone" dataKey="quantidade" stroke="#8884d8" />
            </LineChart>
          </ResponsiveContainer>
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-2">🔥 Clientes Mais Ativos</h3>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={data.clientes_mais_ativos}>
              <XAxis dataKey="nome" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="quantidade" fill="#82ca9d" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div>
        <h3 className="text-lg font-semibold mb-2">🥇 Itens Mais Comprados</h3>
        <ResponsiveContainer width="100%" height={300}>
          <PieChart>
            <Pie
              data={data.itens_mais_comprados}
              dataKey="quantidade"
              nameKey="nome"
              cx="50%"
              cy="50%"
              outerRadius={100}
              label
            >
              {Array.isArray(data.itens_mais_comprados) && data.itens_mais_comprados.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Legend />
            <Tooltip />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}
