import { useEffect, useState } from 'react'
import axios from 'axios'
import {Link} from "react-router-dom";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

interface Cliente {
  id: string
  nome: string
  cnpj: string
  telefone: string
  email: string
  endereco: string
  dias_compra: { dia_semana: number }[]
  itens: { nome: string }[]
}

interface Alerta {
  cliente_id: string
  nome_cliente: string
  motivo: string
  itens_faltantes?: string[]
}

export default function Home() {
  const [clientes, setClientes] = useState<Cliente[]>([])
  const [alertas, setAlertas] = useState<Alerta[]>([])

  useEffect(() => {
    api.get('/clientes').then((res) => {
      const normalized = res.data.map((c: any) => ({
        ...c,
        itens: Array.isArray(c.itens) ? c.itens : [],
        dias_compra: Array.isArray(c.dias_compra) ? c.dias_compra : [],
      }))
      setClientes(normalized)
    })
    api.get('/alertas/hoje').then((res) => setAlertas(Array.isArray(res.data) ? res.data : []))
  }, [])

  return (
    <>
      <section className="mb-10">
        <h2 className="text-xl font-semibold mb-2">⚠️ Alertas de Hoje</h2>
        {!alertas || alertas.length === 0 ? (
          <p className="text-gray-500">Nenhum alerta encontrado.</p>
        ) : (
          <ul className="space-y-3">
            {alertas.map((alerta, i) => (
              <li key={i} className="bg-red-100 border border-red-400 rounded p-4">
                <strong>{alerta.nome_cliente}</strong>: {alerta.motivo}
                {alerta.itens_faltantes && (
                  <ul className="ml-4 list-disc text-sm text-red-700">
                    {alerta.itens_faltantes.map((item, idx) => (
                      <li key={idx}>{item}</li>
                    ))}
                  </ul>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-2">📋 Clientes</h2>
        <div className="grid md:grid-cols-2 gap-4">
          {clientes.map((c) => (
            <div key={c.id} className="bg-white p-4 rounded shadow">
              <h3 className="font-bold text-lg">{c.nome}</h3>
              <Link
                  to={`/clientes/${c.id}/historico`}
                  className="text-sm text-blue-600 hover:underline"
              >
                Ver histórico
              </Link>
              <Link
                  to={`/clientes/${c.id}`}
                  className="text-sm text-yellow-600 hover:underline ml-4"
              >
                Editar
              </Link>

              <p className="text-sm">CNPJ: {c.cnpj}</p>
              <p className="text-sm">Telefone: {c.telefone}</p>
              <p className="text-sm">Email: {c.email}</p>
              <p className="text-sm">Endereço: {c.endereco}</p>
              <p className="text-sm mt-2">🛒 Itens: {Array.isArray(c.itens) ? c.itens.map(i => i.nome).join(', ') : ''}</p>
              <p className="text-sm">
                📆 Dias de Compra:{' '}
                {Array.isArray(c.dias_compra)
                  ? c.dias_compra.map(d =>
                      ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab'][d.dia_semana]
                    ).join(', ')
                  : ''}
              </p>
            </div>
          ))}
        </div>
      </section>
    </>
  )
}
