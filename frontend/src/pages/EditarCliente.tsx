import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
})

const diasSemana = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sab']

export default function EditarCliente() {
    const { id } = useParams()
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

    const [errors, setErrors] = useState<Record<string, string>>({})
    const [erroServidor, setErroServidor] = useState<string | null>(null)

    useEffect(() => {
        api.get(`/clientes/${id}`)
            .then((res) => {
                const cliente = res.data
                setForm({
                    nome: cliente.nome,
                    cnpj: cliente.cnpj,
                    telefone: cliente.telefone,
                    email: cliente.email || '',
                    endereco: cliente.endereco,
                    itens: cliente.itens.map((i: any) => i.nome),
                    dias_compra: cliente.dias_compra.map((d: any) => d.dia_semana),
                })
            })
            .catch(() => setErroServidor('Erro ao carregar cliente'))
    }, [id])

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

    const validar = () => {
        const novosErros: Record<string, string> = {}
        if (!form.nome) novosErros.nome = "Nome é obrigatório"
        if (!form.cnpj) novosErros.cnpj = "CNPJ é obrigatório"
        if (!form.telefone) novosErros.telefone = "Telefone é obrigatório"
        if (!form.endereco) novosErros.endereco = "Endereço é obrigatório"
        return novosErros
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const validados = validar()
        if (Object.keys(validados).length > 0) {
            setErrors(validados)
            return
        }

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
            await api.put(`/clientes/${id}`, payload)
            navigate('/')
        } catch (err) {
            setErroServidor('Erro ao atualizar cliente')
        }
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-4 max-w-xl mx-auto">
            <h2 className="text-2xl font-semibold">Editar Cliente</h2>

            {erroServidor && <p className="text-red-600">{erroServidor}</p>}

            <div>
                <input
                    className={`w-full p-2 border rounded ${errors.nome ? "border-red-500" : ""}`}
                    placeholder="Nome"
                    value={form.nome}
                    onChange={(e) => setForm({ ...form, nome: e.target.value })}
                />
                {errors.nome && <p className="text-red-500 text-sm">{errors.nome}</p>}
            </div>

            <div>
                <input
                    className={`w-full p-2 border rounded ${errors.cnpj ? "border-red-500" : ""}`}
                    placeholder="CNPJ"
                    value={form.cnpj}
                    onChange={(e) => setForm({ ...form, cnpj: e.target.value })}
                />
                {errors.cnpj && <p className="text-red-500 text-sm">{errors.cnpj}</p>}
            </div>

            <div>
                <input
                    className={`w-full p-2 border rounded ${errors.telefone ? "border-red-500" : ""}`}
                    placeholder="Telefone"
                    value={form.telefone}
                    onChange={(e) => setForm({ ...form, telefone: e.target.value })}
                />
                {errors.telefone && <p className="text-red-500 text-sm">{errors.telefone}</p>}
            </div>

            <input
                className="w-full p-2 border rounded"
                placeholder="Email (opcional)"
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
            />

            <div>
                <input
                    className={`w-full p-2 border rounded ${errors.endereco ? "border-red-500" : ""}`}
                    placeholder="Endereço"
                    value={form.endereco}
                    onChange={(e) => setForm({ ...form, endereco: e.target.value })}
                />
                {errors.endereco && <p className="text-red-500 text-sm">{errors.endereco}</p>}
            </div>

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
                className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
            >
                Salvar alterações
            </button>
            <button
                type="button"
                onClick={async () => {
                    const confirmar = window.confirm("Tem certeza que deseja excluir este cliente?")
                    if (!confirmar) return

                    try {
                        await api.delete(`/clientes/${id}`)
                        navigate('/')
                    } catch (err) {
                        setErroServidor("Erro ao excluir cliente")
                    }
                }}
                className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 ml-4"
            >
                Excluir cliente
            </button>

        </form>
    )
}
