document.addEventListener('DOMContentLoaded', () => {
    carregarChapeus();
});

function carregarChapeus() {
    fetch('http://localhost:8080/hats')
        .then(response => response.json())
        .then(hats => {
            const hatsDiv = document.getElementById('hats');
            hatsDiv.innerHTML = '';
            hats.forEach(hat => {
                const hatCard = document.createElement('div');
                hatCard.className = 'hat-card';
                hatCard.innerHTML = `
                    <img src="imagens/chapeu-${hat.id}.jpg" alt="${hat.nome}" onerror="this.src='imagens/default.jpg'">
                    <h3>${hat.nome}</h3>
                    <p>Preço: R$ ${hat.preco.toFixed(2)}</p>
                    <button onclick="adicionarAoCarrinho(${hat.id}, '${hat.nome}', ${hat.preco})">Adicionar ao carrinho</button>
                `;
                hatsDiv.appendChild(hatCard);
            });
        })
        .catch(() => {
            document.getElementById('hats').innerText = 'Não foi possível carregar os chapéus.';
        });
}

// Carrinho simples em memória
let carrinho = [];

function adicionarAoCarrinho(id, nome, preco) {
    const itemExistente = carrinho.find(item => item.id === id);
    if (itemExistente) {
        itemExistente.quantidade += 1;
    } else {
        carrinho.push({ id, nome, preco, quantidade: 1 });
    }
    atualizarCarrinho();
}

function atualizarCarrinho() {
    const cartDiv = document.getElementById('cart-items');
    cartDiv.innerHTML = '';
    let total = 0;
    carrinho.forEach(item => {
        total += item.preco * item.quantidade;
        const itemDiv = document.createElement('div');
        itemDiv.innerHTML = `
            <span>${item.nome} - R$ ${item.preco.toFixed(2)} x ${item.quantidade}</span>
            <button onclick="alterarQuantidade(${item.id}, 1)">+</button>
            <button onclick="alterarQuantidade(${item.id}, -1)">-</button>
        `;
        cartDiv.appendChild(itemDiv);
    });
    document.getElementById('total-amount').innerText = `Total: R$ ${total.toFixed(2)}`;
}

function alterarQuantidade(id, delta) {
    const item = carrinho.find(item => item.id === id);
    if (!item) return;
    item.quantidade += delta;
    if (item.quantidade <= 0) {
        carrinho = carrinho.filter(i => i.id !== id);
    }
    atualizarCarrinho();
}

// Finalizar compra (exemplo)
document.getElementById('checkout-button').addEventListener('click', () => {
    if (carrinho.length === 0) {
        alert('Seu carrinho está vazio!');
        return;
    }
    const formaPagamento = prompt('Selecione a forma de pagamento: PIX ou Boleto');
    if (formaPagamento && (formaPagamento.toLowerCase() === 'pix' || formaPagamento.toLowerCase() === 'boleto')) {
        alert(`Compra finalizada com pagamento via ${formaPagamento}!`);
        carrinho = [];
        atualizarCarrinho();
    } else {
        alert('Forma de pagamento inválida.');
    }
});