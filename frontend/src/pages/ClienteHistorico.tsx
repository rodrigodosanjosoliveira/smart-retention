import { useEffect, useState } from "react";
import {Link, useParams} from "react-router-dom";
import axios from "axios";

interface Cliente {
    nome: string;
    cnpj: string;
    telefone: string;
    endereco: string;
}

interface ItemCompra {
    nome: string;
    preco: number;
}

interface Compra {
    data: string;
    itens: ItemCompra[];
}

export default function ClienteHistorico() {
    const { id } = useParams<{ id: string }>()
    const [cliente, setCliente] = useState<Cliente | null>(null)
    const [compras, setCompras] = useState<Compra[]>([])
    const [carregando, setCarregando] = useState(true)

    useEffect(() => {
        axios.get(`${import.meta.env.VITE_API_URL}/clientes/${id}/historico`)
            .then(res => {
                setCliente(res.data.cliente)
                setCompras(res.data.historico)
            })
            .finally(() => setCarregando(false))
    }, [id])

    const formatarData = (iso: string) => {
        const data = new Date(iso)
        return data.toLocaleDateString("pt-BR")
    }

    if (carregando) return <p className="p-4">Carregando...</p>

    return (
        <div className="max-w-3xl mx-auto p-6">
            <Link to="/" className="text-blue-600 hover:underline text-sm mb-4 inline-block">
                ← Voltar para o dashboard
            </Link>
            <h1 className="text-2xl font-bold mb-4">Histórico de {cliente?.nome}</h1>
            <div className="mb-4 text-sm text-gray-700">
                <p><strong>CNPJ:</strong> {cliente?.cnpj}</p>
                <p><strong>Telefone:</strong> {cliente?.telefone}</p>
                <p><strong>Endereço:</strong> {cliente?.endereco}</p>
            </div>

            <h2 className="text-xl font-semibold mb-2">Compras</h2>
            {compras.length === 0 ? (
                <p className="text-gray-500">Nenhuma compra registrada.</p>
            ) : (
                <ul className="space-y-6">
                    {compras.map((compra, i) => (
                        <li key={i} className="border p-4 rounded bg-white shadow">
                            <p className="text-sm text-gray-600">{formatarData(compra.data)}</p>
                            <ul className="list-disc ml-5 mt-2 text-sm">
                                {compra.itens.map((item, j) => (
                                    <li key={j}>
                                        {item.nome} – R$ {item.preco.toLocaleString("pt-BR", {
                                        minimumFractionDigits: 2,
                                        maximumFractionDigits: 2,
                                    })}
                                    </li>
                                ))}
                            </ul>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    )
}
