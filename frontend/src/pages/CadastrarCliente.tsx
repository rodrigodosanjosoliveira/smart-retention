import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
})

const diasSemana = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab']

export default function CadastrarCliente() {
  const navigate = useNavigate()
  const [form, setForm] = useState({
    nome: '',
    cnpj: '',
    telefone: '',
    email: '',
    endereco: '',
    itens: [''],
    dias_compra: [] as number[],
  })

  const toggleDia = (dia: number) => {
    setForm((prev) => ({
      ...prev,
      dias_compra: prev.dias_compra.includes(dia)
        ? prev.dias_compra.filter((d) => d !== dia)
        : [...prev.dias_compra, dia],
    }))
  }

  const handleItemChange = (index: number, value: string) => {
    const novosItens = [...form.itens]
    novosItens[index] = value
    setForm((prev) => ({ ...prev, itens: novosItens }))
  }

  const adicionarItem = () => {
    setForm((prev) => ({ ...prev, itens: [...prev.itens, ''] }))
  }

  const removerItem = (index: number) => {
    const novosItens = form.itens.filter((_, i) => i !== index)
    setForm((prev) => ({ ...prev, itens: novosItens }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const payload = {
      nome: form.nome,
      cnpj: form.cnpj,
      telefone: form.telefone,
      email: form.email,
      endereco: form.endereco,
      itens: form.itens.filter(i => i.trim() !== '').map(i => ({ nome: i })),
      dias_compra: form.dias_compra.map(d => ({ dia_semana: d })),
    }

    try {
      await api.post('/clientes', payload)
      navigate('/')
    } catch (err) {
      alert('Erro ao cadastrar cliente')
      console.error(err)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-xl mx-auto">
      <h2 className="text-2xl font-semibold">Cadastrar Cliente</h2>

      <input
        className="w-full p-2 border rounded"
        placeholder="Nome"
        value={form.nome}
        onChange={(e) => setForm({ ...form, nome: e.target.value })}
      />
      <input
        className="w-full p-2 border rounded"
        placeholder="CNPJ"
        value={form.cnpj}
        onChange={(e) => setForm({ ...form, cnpj: e.target.value })}
      />
      <input
        className="w-full p-2 border rounded"
        placeholder="Telefone"
        value={form.telefone}
        onChange={(e) => setForm({ ...form, telefone: e.target.value })}
      />
    <input
        className="w-full p-2 border rounded"
        placeholder="Email"
        value={form.email}
        onChange={(e) => setForm({ ...form, email: e.target.value })}
    />
        
      <input
        className="w-full p-2 border rounded"
        placeholder="Endereço"
        value={form.endereco}
        onChange={(e) => setForm({ ...form, endereco: e.target.value })}
      />

      <div>
        <label className="block mb-1 font-semibold">Itens que costuma comprar:</label>
        {form.itens.map((item, index) => (
          <div key={index} className="flex gap-2 mb-2">
            <input
              className="flex-1 p-2 border rounded"
              placeholder={`Item ${index + 1}`}
              value={item}
              onChange={(e) => handleItemChange(index, e.target.value)}
            />
            <button type="button" onClick={() => removerItem(index)} className="text-red-500">Remover</button>
          </div>
        ))}
        <button type="button" onClick={adicionarItem} className="text-blue-500">+ Adicionar Item</button>
      </div>

      <div>
        <label className="block mb-1 font-semibold">Dias da Semana que costuma comprar:</label>
        <div className="flex flex-wrap gap-3">
          {diasSemana.map((dia, index) => (
            <label key={index} className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={form.dias_compra.includes(index)}
                onChange={() => toggleDia(index)}
              />
              {dia}
            </label>
          ))}
        </div>
      </div>

      <button
        type="submit"
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
      >
        Cadastrar
      </button>
    </form>
  )
}
