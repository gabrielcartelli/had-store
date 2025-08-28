// frontend/js/pedidos.js

document.addEventListener('DOMContentLoaded', () => {
    // Redireciona para login se não estiver autenticado
    if (!localStorage.getItem('jwt_token')) {
        window.location.href = 'auth.html';
        return;
    }
    const consultaForm = document.getElementById('consulta-form-pedidos');
    if (!consultaForm) return;

    // Aplica máscara de CPF no campo de input
    $('#cpf-consulta-pedidos').mask('000.000.000-00');

    // Adiciona o event listener para o envio do formulário
    consultaForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        
        const cpfInput = document.getElementById('cpf-consulta-pedidos');
        const pedidosListDiv = document.getElementById('pedidos-list');
        const cpf = cpfInput.value.trim();

        if (!cpf) {
            alert('Por favor, informe um CPF.');
            return;
        }

        pedidosListDiv.innerHTML = "<p>Buscando seus pedidos...</p>";
        showLoader();

        try {
            // Usa a função fetchApi centralizada para fazer a chamada
            const pedidos = await fetchApi(`/pedidos?cpf=${encodeURIComponent(cpf)}`);
            
            pedidosListDiv.innerHTML = ''; // Limpa a mensagem de "Buscando..."
            if (!pedidos || pedidos.length === 0) {
                pedidosListDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
                return;
            }

            // Cria e exibe os cards de pedido
            pedidos.forEach(pedido => {
                const card = document.createElement('div');
                card.className = 'pedido-card';
                card.innerHTML = `
                    <p><strong>Cliente:</strong> ${pedido.nome}</p>
                    <p><strong>CPF:</strong> ${pedido.cpf}</p>
                    <p><strong>Pagamento:</strong> ${pedido.pagamento}</p>
                    <p><strong>Total:</strong> R$ ${pedido.total.toFixed(2)}</p>
                    <div class="pedido-itens">
                        <strong>Itens do Pedido:</strong>
                        <ul>
                            ${pedido.itens.map(item => `<li>${item.nome} (${item.quantidade || 1}x) - R$ ${item.preco.toFixed(2)}</li>`).join('')}
                        </ul>
                    </div>
                `;
                pedidosListDiv.appendChild(card);
            });

        } catch (error) {
            console.error('Erro ao consultar pedidos:', error);
            pedidosListDiv.innerHTML = "<p>Ocorreu um erro ao buscar seus pedidos. Tente novamente mais tarde.</p>";
        } finally {
            hideLoader();
        }
    });
});