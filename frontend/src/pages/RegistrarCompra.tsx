import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

interface Item {
  id: string
  nome: string
}

interface Cliente {
  id: string
  nome: string
  itens: Item[]
}

export default function RegistrarCompra() {
  const navigate = useNavigate()
  const [clientes, setClientes] = useState<Cliente[]>([])
  const [clienteSelecionado, setClienteSelecionado] = useState<Cliente | null>(null)
  const [itensSelecionados, setItensSelecionados] = useState<string[]>([])
  const [dataCompra, setDataCompra] = useState(() => new Date().toISOString().split('T')[0])

  useEffect(() => {
    api.get('/clientes').then((res) => {
      setClientes(res.data)
    })
  }, [])

  const handleClienteChange = (id: string) => {
    const cliente = clientes.find(c => c.id === id) || null
    setClienteSelecionado(cliente)
    setItensSelecionados([])
  }

  const toggleItem = (itemId: string) => {
    setItensSelecionados(prev =>
      prev.includes(itemId) ? prev.filter(id => id !== itemId) : [...prev, itemId]
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!clienteSelecionado || itensSelecionados.length === 0) {
      alert("Selecione um cliente e pelo menos um item.")
      return
    }

    const payload = {
      cliente_id: clienteSelecionado.id,
      itens: itensSelecionados.map(id => ({
        item_id: id,
        preco: 0,
      })),
      data: dataCompra,
    }

    try {
      console.log('Payload enviado:', payload)
      await api.post('/compras', payload)
      navigate('/')
    } catch (err) {
      console.error(err)
      alert("Erro ao registrar compra.")
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-xl mx-auto">
      <h2 className="text-2xl font-semibold">Registrar Compra</h2>

      <label className="block">
        Cliente:
        <select
          className="w-full p-2 border rounded mt-1"
          onChange={(e) => handleClienteChange(e.target.value)}
          value={clienteSelecionado?.id || ''}
        >
          <option value="">Selecione um cliente</option>
          {clientes.map(c => (
            <option key={c.id} value={c.id}>{c.nome}</option>
          ))}
        </select>
      </label>

      {clienteSelecionado && (
        <>
          <label className="block">
            Itens comprados:
            <div className="flex flex-wrap gap-3 mt-2">
              {clienteSelecionado.itens.map(item => (
                <label key={item.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={itensSelecionados.includes(item.id)}
                    onChange={() => toggleItem(item.id)}
                  />
                  {item.nome}
                </label>
              ))}
            </div>
          </label>

          <label className="block">
            Data da compra:
            <input
              type="date"
              className="w-full p-2 border rounded mt-1"
              value={dataCompra}
              onChange={(e) => setDataCompra(e.target.value)}
            />
          </label>
        </>
      )}

      <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">
        Registrar Compra
      </button>
    </form>
  )
}
