const pedidoMemoria = window.localStorage.getItem('carrinho');
const carrinho = pedidoMemoria ? JSON.parse(pedidoMemoria) : [];

function mostrarResumoCarrinho() {
    const resumoDiv = document.getElementById('pedido-resumo');
    if (!carrinho || carrinho.length === 0) {
        resumoDiv.innerHTML = '<p>Seu carrinho está vazio.</p>';
        document.getElementById('checkout-form').style.display = 'none';
        return;
    }
    let html = `<h3>Resumo do Pedido</h3><ul>`;
    carrinho.forEach(item => {
        html += `<li>${item.nome} - R$ ${item.preco.toFixed(2)} x ${item.quantidade}</li>`;
    });
    html += `</ul><strong>Total: R$ ${carrinho.reduce((acc, item) => acc + item.preco * item.quantidade, 0).toFixed(2)}</strong>`;
    resumoDiv.innerHTML = html;
}
mostrarResumoCarrinho();

document.getElementById('checkout-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const nome = document.getElementById('nome').value;
    const email = document.getElementById('email').value;
    const pagamento = document.getElementById('pagamento').value;

    const pedido = {
        nome,
        email,
        pagamento,
        itens: carrinho,
        total: carrinho.reduce((acc, item) => acc + item.preco * item.quantidade, 0)
    };

    fetch('/api/pedido', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(pedido)
    })
    .then(res => res.json())
    .then(data => {
        let msg = `<div class="pagamento-info"><h3>Pedido registrado!</h3>`;
        if (pagamento === 'boleto') {
            const codigoBoleto = "23793.38128 60002.123456 12345.678901 1 23456789012345";
            msg += `<p><strong>Boleto gerado:</strong></p>
                <pre id="boleto-codigo">${codigoBoleto}</pre>
                <button onclick="copiarTexto('boleto-codigo')">Copiar código do boleto</button>`;
        } else if (pagamento === 'pix') {
            const chavePix = "chave-pix-falsa-1234567890";
            msg += `<p><strong>Chave PIX:</strong></p>
                <pre id="pix-codigo">${chavePix}</pre>
                <button onclick="copiarTexto('pix-codigo')">Copiar chave PIX</button>`;
        }
        msg += `<p>Obrigado pela compra, ${nome}!</p>
            <button onclick="voltarParaListagem()" style="margin-top:1rem;">Voltar para a loja</button>
            </div>`;
        document.getElementById('resultado').innerHTML = msg;
        localStorage.removeItem('carrinho');
        document.getElementById('checkout-form').style.display = 'none';
        document.getElementById('pedido-resumo').style.display = 'none';
    });
});

// Função para copiar texto do elemento <pre>
function copiarTexto(elementId) {
    const texto = document.getElementById(elementId).innerText;
    navigator.clipboard.writeText(texto).then(() => {
        alert('Copiado para a área de transferência!');
    });
}

// Função para voltar para a listagem
function voltarParaListagem() {
    // Lógica para voltar para a página de listagem de produtos
    window.location.href = '/'; // Altere para a URL da sua página de listagem, se necessário
}