import { useEffect, useState } from 'react'
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

interface Compra {
  id: string
  data: string
  nome_cliente: string
  itens: { nome: string; preco?: number }[]
}

export default function HistoricoCompras() {
  const [compras, setCompras] = useState<Compra[]>([])

  useEffect(() => {
    api.get('/api/compras').then((res) => setCompras(res.data))
  }, [])

  const formatarData = (iso: string) => {
    const data = new Date(iso)
    return data.toLocaleDateString('pt-BR', { timeZone: 'UTC' })
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">🧾 Histórico de Compras</h2>
      {compras?.length === 0 ? (
        <p className="text-gray-500">Nenhuma compra registrada.</p>
      ) : (
        <ul className="space-y-4">
          {compras?.map((compra) => (
            <li key={compra.id} className="bg-white p-4 rounded shadow">
              <p className="text-sm text-gray-500">📅 {formatarData(compra.data)}</p>
              <p className="font-bold text-lg">🧑‍🍳 {compra.nome_cliente}</p>
              <p className="text-sm">
                🛒 Itens:
                {compra.itens.map((i, index) => (
                  <span key={index}>
                    {i.nome}{i.preco !== undefined ? ` (${i.preco.toFixed(2)})` : ''}{index < compra.itens.length - 1 ? ', ' : ''}
                  </span>
                ))}
              </p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
