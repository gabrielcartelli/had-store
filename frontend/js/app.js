let todosChapeus = [];
let carrinho = [];

// Carrega carrinho do localStorage ao iniciar
document.addEventListener('DOMContentLoaded', () => {
    carregarChapeus();
    document.getElementById('search-hat').addEventListener('input', filtrarChapeus);
    const carrinhoSalvo = localStorage.getItem('carrinho');
    if (carrinhoSalvo) {
        carrinho = JSON.parse(carrinhoSalvo);
        atualizarCarrinho();
    }

    const cartMenuButton = document.getElementById('cart-menu-button');
    const cartSidebar = document.getElementById('cart-sidebar');
    const closeCartButton = document.getElementById('close-cart');

    cartMenuButton.addEventListener('click', () => {
        cartSidebar.classList.add('open');
        atualizarCarrinho();
    });

    closeCartButton.addEventListener('click', () => {
        cartSidebar.classList.remove('open');
    });

    window.addEventListener('mousedown', function(e) {
        if (
            cartSidebar.classList.contains('open') &&
            !cartSidebar.contains(e.target) &&
            e.target !== cartMenuButton
        ) {
            cartSidebar.classList.remove('open');
        }
    });
});

function carregarChapeus() {
    fetch('/api/hats')
        .then(response => response.json())
        .then(hats => {
            todosChapeus = hats;
            exibirChapeus(hats);
        })
        .catch(() => {
            document.getElementById('hats').innerText = 'Não foi possível carregar os chapéus.';
        });
}

function exibirChapeus(hats) {
    const hatsDiv = document.getElementById('hats');
    hatsDiv.innerHTML = '';
    if (hats.length === 0) {
        hatsDiv.innerHTML = '<p style="text-align:center;color:#2563eb;font-weight:bold;">Nenhum chapéu encontrado.</p>';
        return;
    }
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
}

function filtrarChapeus() {
    const termo = document.getElementById('search-hat').value.toLowerCase();
    const filtrados = todosChapeus.filter(hat => hat.nome.toLowerCase().includes(termo));
    exibirChapeus(filtrados);
}

function adicionarAoCarrinho(id, nome, preco) {
    const itemExistente = carrinho.find(item => item.id === id);
    if (itemExistente) {
        itemExistente.quantidade += 1;
    } else {
        carrinho.push({ id, nome, preco, quantidade: 1 });
    }
    salvarCarrinho();
    atualizarCarrinho();
    // Abre o carrinho automaticamente ao adicionar item
    document.getElementById('cart-sidebar').classList.add('open');
}

function alterarQuantidade(id, delta) {
    const item = carrinho.find(item => item.id === id);
    if (!item) return;
    item.quantidade += delta;
    if (item.quantidade <= 0) {
        carrinho = carrinho.filter(i => i.id !== id);
    }
    salvarCarrinho();
    atualizarCarrinho();
}

function removerDoCarrinho(id) {
    carrinho = carrinho.filter(item => item.id !== id);
    salvarCarrinho();
    atualizarCarrinho();
}

function salvarCarrinho() {
    localStorage.setItem('carrinho', JSON.stringify(carrinho));
}

function atualizarCarrinho() {
    const cartDiv = document.getElementById('cart-items');
    const totalDiv = document.getElementById('total-amount');
    const checkoutBtn = document.getElementById('checkout-button');
    let total = 0;

    if (!cartDiv || !totalDiv || !checkoutBtn) return;

    if (carrinho.length === 0) {
        cartDiv.innerHTML = '<div class="cart-empty-message">Seu carrinho está vazio.</div>';
        totalDiv.innerHTML = '';
        checkoutBtn.style.display = 'none';
        return;
    }

    cartDiv.innerHTML = '';
    carrinho.forEach(item => {
        total += item.preco * item.quantidade;
        const itemDiv = document.createElement('div');
        itemDiv.className = 'cart-item';
        itemDiv.innerHTML = `
            <div class="cart-item-info">
                <img src="imagens/chapeu-${item.id}.jpg" alt="${item.nome}" onerror="this.src='imagens/default.jpg'">
                <div>
                    <strong>${item.nome}</strong><br>
                    <span>R$ ${item.preco.toFixed(2)}</span>
                </div>
            </div>
            <div class="cart-item-actions">
                <button onclick="alterarQuantidade(${item.id}, -1)" title="Diminuir">-</button>
                <span>${item.quantidade}</span>
                <button onclick="alterarQuantidade(${item.id}, 1)" title="Aumentar">+</button>
                <button onclick="removerDoCarrinho(${item.id})" title="Remover" class="remove-btn">&#10006;</button>
            </div>
        `;
        cartDiv.appendChild(itemDiv);
    });
    totalDiv.innerHTML = `<span style="font-size:1.1rem;font-weight:bold;">Total: R$ ${total.toFixed(2)}</span>`;
    checkoutBtn.style.display = 'block';
}

// Redireciona para página de checkout ao finalizar compra
document.getElementById('checkout-button').addEventListener('click', () => {
    if (carrinho.length === 0) {
        alert('Seu carrinho está vazio!');
        return;
    }
    window.location.href = "checkout.html";
});