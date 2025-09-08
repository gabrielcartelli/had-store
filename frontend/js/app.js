    // Botão de carrinho no header abre a sidebar do carrinho
    const headerCartBtn = document.getElementById('header-cart-btn');
// Variáveis globais para os dados da aplicação
let todosChapeus = [];
let carrinho = [];

// Função principal que roda quando o HTML está pronto
document.addEventListener('DOMContentLoaded', () => {
    // Delegação de evento para garantir que o botão de fechar o carrinho funcione sempre
    document.addEventListener('click', function(e) {
        if (e.target && e.target.id === 'close-cart') {
            const cartSidebar = document.getElementById('cart-sidebar');
            if (cartSidebar) cartSidebar.classList.remove('open');
        }
    });
    // Botão de carrinho no header abre a sidebar do carrinho
    const headerCartBtn = document.getElementById('header-cart-btn');
    const cartSidebar = document.getElementById('cart-sidebar');
    if (headerCartBtn && cartSidebar) {
        headerCartBtn.addEventListener('click', () => {
            cartSidebar.classList.add('open');
            atualizarCarrinho();
        });
    }
    // Aplica máscara de CPF ao input da modal de consulta de pedidos
    $('#cpf-consulta-sidebar').mask('000.000.000-00');

    // Atualiza a UI para refletir o estado de login
    function setupAuthUI() {
        const authLinksContainer = document.getElementById('auth-links');
        if (!authLinksContainer) return;

        const token = localStorage.getItem('jwt_token');

        if (token) {
            // Usuário está logado: exibe o botão de Logout
            authLinksContainer.innerHTML = '<button id="logout-button" class="header-link" style="background:none; border:none; cursor:pointer;">Logout</button>';
            document.getElementById('logout-button').addEventListener('click', () => {
                localStorage.removeItem('jwt_token');
                window.location.reload(); // Recarrega a página para atualizar o estado
            });
        } else {
            // Usuário não está logado: exibe o link de Login
            authLinksContainer.innerHTML = '<a href="auth.html" class="header-link">Login</a>';
        }
    }
    setupAuthUI();


    //
    // INICIALIZAÇÃO DA PÁGINA
    //
    carregarChapeus();
    document.getElementById('search-hat').addEventListener('input', filtrarChapeus);

    const carrinhoSalvo = localStorage.getItem('carrinho');
    if (carrinhoSalvo) {
        carrinho = JSON.parse(carrinhoSalvo);
        atualizarCarrinho();
    }


    //
    // EVENT LISTENERS PARA SIDEBARS E AÇÕES
    //
    const cartMenuButton = document.getElementById('cart-menu-button');
    const closeCartButton = document.getElementById('close-cart');

    if (cartMenuButton) {
        cartMenuButton.addEventListener('click', () => {
            cartSidebar.classList.add('open');
            atualizarCarrinho();
        });

        closeCartButton.addEventListener('click', () => cartSidebar.classList.remove('open'));

        window.addEventListener('mousedown', (e) => {
            if (cartSidebar.classList.contains('open') && !cartSidebar.contains(e.target) && e.target !== cartMenuButton) {
                cartSidebar.classList.remove('open');
            }
        });
    }

    // Lógica para a sidebar de "Meus Pedidos"
    const pedidosMenuButton = document.getElementById('pedidos-menu-button');
    const pedidosSidebar = document.getElementById('pedidos-sidebar');
    const closePedidosButton = document.getElementById('close-pedidos');

    if (pedidosMenuButton) {
        pedidosMenuButton.addEventListener('click', () => {
            pedidosSidebar.classList.add('open');
            document.getElementById('pedidos-list').innerHTML = ''; // Limpa a lista
        });

        closePedidosButton.addEventListener('click', () => pedidosSidebar.classList.remove('open'));

        window.addEventListener('mousedown', (e) => {
            if (pedidosSidebar.classList.contains('open') && !pedidosSidebar.contains(e.target) && e.target !== pedidosMenuButton) {
                pedidosSidebar.classList.remove('open');
            }
        });
    }

    // Consulta de pedidos na sidebar
    const consultaForm = document.getElementById('consulta-form-sidebar'); // Corrigido para o ID da sidebar
    if (consultaForm) {
        consultaForm.addEventListener('submit', handleConsultaPedidos);
    }

    // Cupom no carrinho
    const aplicarCupomBtn = document.getElementById('aplicar-cupom');
    if (aplicarCupomBtn) {
        aplicarCupomBtn.addEventListener('click', aplicarCupom);
    }

    // Botão de Finalizar Compra
    const checkoutBtn = document.getElementById('checkout-button');
    if (checkoutBtn) {
        checkoutBtn.addEventListener('click', handleCheckout);
    }
});


//
// FUNÇÕES DE DADOS E API
//
async function carregarChapeus() {
    showLoader();
    try {
        const hats = await fetchApi('/hats'); // Usa a nova função centralizada
        todosChapeus = hats;
        exibirChapeus(hats);
    } catch (error) {
        console.error("Erro ao carregar chapéus:", error);
        openModal("Erro", "<p>Não foi possível carregar os chapéus. Verifique sua conexão ou o código de acesso.</p>");
    } finally {
        hideLoader();
    }
}

async function handleConsultaPedidos(event) {
    event.preventDefault();
    const cpfInput = document.getElementById('cpf-consulta-sidebar'); // Corrigido
    const pedidosListDiv = document.getElementById('pedidos-list-sidebar'); // Corrigido

    const cpf = cpfInput.value.trim();
    if (!cpf) return;

    pedidosListDiv.innerHTML = "<p>Buscando pedidos...</p>";

    try {
        const pedidos = await fetchApi(`/pedidos?cpf=${encodeURIComponent(cpf)}`);
        pedidosListDiv.innerHTML = '';
        if (!pedidos || pedidos.length === 0) {
            pedidosListDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
            return;
        }
        pedidos.forEach(pedido => {
            const card = document.createElement('div');
            card.className = 'pedido-card';
            card.innerHTML = `
                <strong>Pedido:</strong> ${pedido.nome}<br>
                <strong>CPF:</strong> ${pedido.cpf}<br>
                <strong>E-mail:</strong> ${pedido.email}<br>
                <strong>Pagamento:</strong> ${pedido.pagamento}<br>
                <strong>Total:</strong> R$ ${pedido.total.toFixed(2)}
                <div class="pedido-itens">
                    <strong>Itens:</strong>
                    <ul>
                        ${pedido.itens.map(item => `<li>${item.nome} (${item.quantidade || 1}x) - R$ ${item.price.toFixed(2)}</li>`).join('')}
                    </ul>
                </div>
            `;
            pedidosListDiv.appendChild(card);
        });
    } catch (error) {
        console.error("Erro ao consultar pedidos:", error);
        pedidosListDiv.innerHTML = "<p>Erro ao consultar pedidos.</p>";
    }
}


//
// FUNÇÕES DE UI E MANIPULAÇÃO DO DOM
//
function exibirChapeus(hats) {
    const hatsDiv = document.getElementById('hats');
    if (!hatsDiv) return;
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
            <p>Preço: R$ ${hat.price.toFixed(2)}</p>
            <button class="add-to-cart-btn" data-id="${hat.id}" data-nome="${hat.nome}" data-price="${hat.price}">Adicionar ao carrinho</button>
        `;
        hatsDiv.appendChild(hatCard);
    });

    // Adiciona event listeners aos novos botões
    document.querySelectorAll('.add-to-cart-btn').forEach(button => {
        button.addEventListener('click', (event) => {
            const { id, nome, price } = event.target.dataset;
            adicionarAoCarrinho(parseInt(id), nome, parseFloat(price));
        });
    });
}

function filtrarChapeus() {
    const termo = 'Beanie';// document.getElementById('search-hat').value.toLowerCase();
    const filtrados = todosChapeus.filter(hat => hat.nome.toLowerCase().includes(termo));
    exibirChapeus(filtrados);
}

function aplicarCupom() {
    const cupomInput = document.getElementById('cupom');
    const erroCupom = document.getElementById('erro-cupom');
    if (!cupomInput || !erroCupom) return;

    const cupom = cupomInput.value.trim().toUpperCase();
    if (cupom === "HATOFF") {
        localStorage.setItem('cupom', cupom);
        erroCupom.textContent = "Cupom aplicado!";
        erroCupom.className = 'cupom-sucesso';
        atualizarCarrinho();
    } else {
        erroCupom.textContent = "Cupom inválido.";
        erroCupom.className = 'input-error';
        localStorage.removeItem('cupom');
    }
}


//
// FUNÇÕES DO CARRINHO
//
function adicionarAoCarrinho(id, nome, price) {
    const itemExistente = carrinho.find(item => item.id === id);
    if (itemExistente) {
        itemExistente.quantidade += 1;
    } else {
        carrinho.push({ id, nome, price, quantidade: 1 });
    }
    salvarCarrinho();
    atualizarCarrinho();
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
    if (!cartDiv || !totalDiv || !checkoutBtn) return;

    if (carrinho.length === 0) {
        cartDiv.innerHTML = '<div class="cart-empty-message">Seu carrinho está vazio.</div>';
        totalDiv.innerHTML = '';
        checkoutBtn.style.display = 'none';
        return;
    }

    let total = 0;
    cartDiv.innerHTML = '';
    carrinho.forEach(item => {
        total += item.price * item.quantidade;
        const itemDiv = document.createElement('div');
        itemDiv.className = 'cart-item';
        itemDiv.innerHTML = `
            <div class="cart-item-info">
                <span class="cart-item-nome">${item.nome}</span>
                <span class="cart-item-preco">R$ ${item.price.toFixed(2)}</span>
                <div class="cart-item-actions">
                    <button class="qty-btn" data-id="${item.id}" data-delta="-1" title="Diminuir">-</button>
                    <span>${item.quantidade}</span>
                    <button class="qty-btn" data-id="${item.id}" data-delta="1" title="Aumentar">+</button>
                    <button class="remove-btn" data-id="${item.id}" title="Remover" >&#10006;</button>
                </div>
            </div>
        `;
        cartDiv.appendChild(itemDiv);
    });

    // Adiciona event listeners aos botões do carrinho
    cartDiv.querySelectorAll('.qty-btn').forEach(btn =>
        btn.addEventListener('click', e => alterarQuantidade(parseInt(e.target.dataset.id), parseInt(e.target.dataset.delta)))
    );
    cartDiv.querySelectorAll('.remove-btn').forEach(btn =>
        btn.addEventListener('click', e => removerDoCarrinho(parseInt(e.target.dataset.id)))
    );

    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let descontoHtml = '';
    let totalComDesconto = total;

    if (cupom === 'HATOFF') {
        const desconto = total * 0.18;
        totalComDesconto -= desconto;
        descontoHtml = `<div class="cart-cupom-info">Cupom HATOFF aplicado: -R$ ${desconto.toFixed(2)}</div>`;
    }

    totalDiv.innerHTML = `
        <div class="cart-subtotal">Subtotal: <span>R$ ${total.toFixed(2)}</span></div>
        ${descontoHtml}
        <div class="cart-total">Total com desconto: <span>R$ ${totalComDesconto.toFixed(2)}</span></div>
    `;
    checkoutBtn.style.display = 'block';
}

function handleCheckout() {
    if (carrinho.length === 0) {
        alert('Seu carrinho está vazio!');
        return;
    }
    window.location.href = "checkout.html";
}


//
// FUNÇÕES UTILITÁRIAS (Loader, Modal)
//
function showLoader() {
    const loader = document.getElementById('loader');
    if (loader) loader.style.display = 'flex';
}
function hideLoader() {
    const loader = document.getElementById('loader');
    if (loader) loader.style.display = 'none';
}

function openModal(title, bodyHtml) {
    const modal = document.getElementById('modal');
    if (modal) {
        document.getElementById('modal-title').innerText = title;
        document.getElementById('modal-body').innerHTML = bodyHtml;
        modal.style.display = 'flex';
    }
}
function closeModal() {
    const modal = document.getElementById('modal');
    if (modal) modal.style.display = 'none';
}
// Expondo closeModal globalmente para o onclick no HTML
window.closeModal = closeModal;