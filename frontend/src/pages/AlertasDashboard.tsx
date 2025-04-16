import { useEffect, useState } from 'react'
import axios from 'axios'
import { Link } from 'react-router-dom'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

interface Alerta {
  cliente_id: string
  nome_cliente: string
  tipo: string
  motivo: string
  itens_faltantes?: string[]
}

export default function AlertasDashboard() {
  const [alertas, setAlertas] = useState<Alerta[]>([])

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws/alertas")

    socket.onmessage = (event) => {
      const novosAlertas = JSON.parse(event.data)
      setAlertas(novosAlertas)
    }

    socket.onerror = (err) => {
      console.error("Erro no WebSocket:", err)
    }

    const fetchFallback = () => {
      api.get('/alertas').then(res => setAlertas(res.data))
    }

    const interval = setInterval(fetchFallback, 30000) // fallback a cada 30s

    return () => {
      socket.close()
      clearInterval(interval)
    }
  }, [])

  const alertasPorTipo = (tipo: string) => alertas.filter(a => a.tipo === tipo)
  const totalAlertas = alertas.length

  const formatarData = (iso: string) => {
    const data = new Date(iso)
    return data.toLocaleDateString('pt-BR', { timeZone: 'UTC' })
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div className="flex justify-between items-center">
        <h2 className="text-2xl font-semibold mb-4">🔔 Alertas Inteligentes</h2>
        <Link to="/alertas" className="relative inline-block">
          <span className="text-blue-600 hover:underline">Ver todos</span>
          {totalAlertas > 0 && (
            <span className="absolute -top-2 -right-4 bg-red-600 text-white text-xs font-bold rounded-full px-2">
              {totalAlertas}
            </span>
          )}
        </Link>
      </div>

      <section>
        <h3 className="text-lg font-bold mb-2 text-red-600">🚨 Clientes Inativos</h3>
        {alertasPorTipo("inatividade").length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum cliente inativo.</p>
        ) : (
          <ul className="space-y-2">
            {alertasPorTipo("inatividade").map((a, i) => (
              <li key={i} className="bg-red-100 border-l-4 border-red-500 p-4 rounded">
                <strong>{a.nome_cliente}</strong>: {a.motivo}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h3 className="text-lg font-bold mb-2 text-orange-500">📅 Ausente no Dia Previsto</h3>
        {alertasPorTipo("dia_previsto").length === 0 ? (
          <p className="text-sm text-gray-500">Todos compraram no dia esperado.</p>
        ) : (
          <ul className="space-y-2">
            {alertasPorTipo("dia_previsto").map((a, i) => (
              <li key={i} className="bg-orange-100 border-l-4 border-orange-400 p-4 rounded">
                <strong>{a.nome_cliente}</strong>: {a.motivo}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h3 className="text-lg font-bold mb-2 text-yellow-600">📉 Itens Deixados de Comprar</h3>
        {alertasPorTipo("item_faltando").length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum item ausente identificado.</p>
        ) : (
          <ul className="space-y-4">
            {alertasPorTipo("item_faltando").map((a, i) => (
              <li key={i} className="bg-yellow-100 border-l-4 border-yellow-500 p-4 rounded">
                <p><strong>{a.nome_cliente}</strong>: {a.motivo}</p>
                <ul className="ml-4 list-disc text-sm mt-1">
                  {a.itens_faltantes?.map((item, j) => (
                    <li key={j}>{item}</li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}
