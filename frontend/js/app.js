const API_UUID = "e3e6c6c2-9b7d-4c5e-8c1a-2f7b8f8e2a1d";

// Solicita o UUID ao usuário se necessário
function solicitarUUID() {
    let uuid = localStorage.getItem("api_uuid");
    while (uuid !== API_UUID) {
        uuid = prompt("Informe o código de acesso:");
        if (uuid === null) {
            document.body.innerHTML = "<h2 style='color:#e11d48;text-align:center;margin-top:3rem;'>Acesso bloqueado.</h2>";
            throw new Error("Acesso bloqueado.");
        }
        localStorage.setItem("api_uuid", uuid);
    }
}
solicitarUUID();

// Função para fetch com UUID
function fetchComUUID(url, options = {}) {
    if (!options.headers) options.headers = {};
    options.headers["X-API-UUID"] = localStorage.getItem("api_uuid");
    return fetch(url, options);
}

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

    // Pedidos menu hamburguer
    const pedidosMenuButton = document.getElementById('pedidos-menu-button');
    const pedidosSidebar = document.getElementById('pedidos-sidebar');
    const closePedidosButton = document.getElementById('close-pedidos');
    const consultaForm = document.getElementById('consulta-form');
    const pedidosListDiv = document.getElementById('pedidos-list');

    pedidosMenuButton.addEventListener('click', () => {
        pedidosSidebar.classList.add('open');
        pedidosListDiv.innerHTML = '';
    });

    closePedidosButton.addEventListener('click', () => {
        pedidosSidebar.classList.remove('open');
    });

    window.addEventListener('mousedown', function(e) {
        if (
            pedidosSidebar.classList.contains('open') &&
            !pedidosSidebar.contains(e.target) &&
            e.target !== pedidosMenuButton
        ) {
            pedidosSidebar.classList.remove('open');
        }
    });

    // Consulta de pedidos
    if (consultaForm) {
        consultaForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const cpf = document.getElementById('cpf-consulta').value.trim();
            if (!cpf) return;

            document.getElementById('pedidos-list').innerHTML = "<p>Buscando pedidos...</p>";

            fetchComUUID(`/api/pedidos?cpf=${encodeURIComponent(cpf)}`)
                .then(res => res.json())
                .then(pedidos => {
                    const pedidosDiv = document.getElementById('pedidos-list');
                    pedidosDiv.innerHTML = '';
                    if (!pedidos || pedidos.length === 0) {
                        pedidosDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
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
                                    ${pedido.itens.map(item => `<li>${item.nome} (${item.quantidade || 1}x) - R$ ${item.preco.toFixed(2)}</li>`).join('')}
                                </ul>
                            </div>
                        `;
                        pedidosDiv.appendChild(card);
                    });
                })
                .catch(() => {
                    document.getElementById('pedidos-list').innerHTML = "<p>Erro ao consultar pedidos.</p>";
                });
        });

        // Máscara para CPF da consulta
        if (window.jQuery && window.jQuery.fn.mask) {
            $('#cpf-consulta').mask('000.000.000-00');
        } else {
            // Carrega jQuery e jQuery Mask se não estiverem presentes
            const jq = document.createElement('script');
            jq.src = "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.7.1/jquery.min.js";
            jq.onload = function() {
                const mask = document.createElement('script');
                mask.src = "https://cdnjs.cloudflare.com/ajax/libs/jquery.mask/1.14.16/jquery.mask.min.js";
                mask.onload = function() {
                    $('#cpf-consulta').mask('000.000.000-00');
                };
                document.body.appendChild(mask);
            };
            document.body.appendChild(jq);
        }
    }

    // Cupom no carrinho
    const cupomInput = document.getElementById('cupom');
    const aplicarCupomBtn = document.getElementById('aplicar-cupom');
    const erroCupom = document.getElementById('erro-cupom');
    if (cupomInput && aplicarCupomBtn && erroCupom) {
        aplicarCupomBtn.onclick = function() {
            const cupom = cupomInput.value.trim().toUpperCase();
            if (cupom) {
                localStorage.setItem('cupom', cupom);
                erroCupom.textContent = "Cupom aplicado!";
                erroCupom.classList.remove('input-error');
                erroCupom.classList.add('cupom-sucesso');
            } else {
                erroCupom.textContent = "Digite um cupom válido.";
                erroCupom.classList.remove('cupom-sucesso');
                erroCupom.classList.add('input-error');
            }
        };
    }
});

function carregarChapeus() {
    fetchComUUID('/api/hats')
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
                <span class="cart-item-nome">${item.nome}</span>
                <span class="cart-item-preco">R$ ${item.preco.toFixed(2)}</span>
                <div class="cart-item-actions">
                    <button onclick="alterarQuantidade(${item.id}, -1)" title="Diminuir">-</button>
                    <span>${item.quantidade}</span>
                    <button onclick="alterarQuantidade(${item.id}, 1)" title="Aumentar">+</button>
                    <button onclick="removerDoCarrinho(${item.id})" title="Remover" class="remove-btn">&#10006;</button>
                </div>
            </div>
        `;
        cartDiv.appendChild(itemDiv);
    });

    // Cupom de desconto
    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let desconto = 0;
    let descontoHtml = '';
    let totalComDesconto = total;

    if (cupom === 'HAD10') {
        desconto = total * 0.10;
        totalComDesconto = total - desconto;
        descontoHtml = `
            <div class="cart-cupom-info" style="color:#22c55e;font-weight:bold;margin-top:0.5em;">
                Cupom HAD10 aplicado: -R$ ${desconto.toFixed(2)}
            </div>
        `;
    }

    totalDiv.innerHTML = `
        <div style="font-size:1rem;">
            <span style="color:#888;">Subtotal:</span> <span style="font-weight:bold;">R$ ${total.toFixed(2)}</span>
        </div>
        ${descontoHtml}
        <div style="font-size:1.1rem;font-weight:bold;margin-top:0.5em;">
            Total com desconto: R$ ${totalComDesconto.toFixed(2)}
        </div>
    `;
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

// Loader functions
function showLoader() {
    document.getElementById('loader').style.display = 'flex';
}
function hideLoader() {
    document.getElementById('loader').style.display = 'none';
}

// Modal functions
function openModal(title, bodyHtml) {
    const modalTitle = document.getElementById('modal-title');
    const modalBody = document.getElementById('modal-body');
    const modal = document.getElementById('modal');
    if (modalTitle && modalBody && modal) {
        modalTitle.innerText = title;
        modalBody.innerHTML = bodyHtml;
        modal.style.display = 'flex';
    }
}
function closeModal() {
    const modal = document.getElementById('modal');
    if (modal) modal.style.display = 'none';
}

// Exemplo de uso do loader e modal em uma chamada de API
function buscarChapeus() {
    showLoader();
    fetchComUUID('/api/hats')
        .then(res => res.json())
        .then(hats => {
            hideLoader();
            // Renderize os chapéus normalmente
            // ...
        })
        .catch(() => {
            hideLoader();
            openModal("Erro", "<p>Não foi possível carregar os chapéus.</p>");
        });
}

// Exemplo de uso do modal para detalhes do produto
function mostrarDetalhesChapeu(chapeu) {
    openModal(
        chapeu.nome,
        `<img src="imagens/chapeu-${chapeu.id}.jpg" style="width:120px;border-radius:10px;margin-bottom:1rem;">
         <p>Preço: <strong>R$ ${chapeu.preco.toFixed(2)}</strong></p>
         <p>Descrição: Chapéu estiloso para todas as ocasiões.</p>`
    );
}

// Adicione animações extras conforme desejar nos eventos JS