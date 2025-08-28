document.addEventListener('DOMContentLoaded', function() {
    // Redireciona para login se não estiver autenticado
    if (!localStorage.getItem('jwt_token')) {
        window.location.href = 'auth.html';
        return;
    }
    //
    // PROTEÇÃO DA PÁGINA E CONFIGURAÇÃO INICIAL
    //
    const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
    
    // Monta o resumo do pedido na tela
    montarResumoPedido(carrinho);
    
    // Aplica máscaras aos inputs do formulário
    aplicarMascaras();

    
    //
    // EVENT LISTENERS DO FORMULÁRIO
    //
    const form = document.getElementById('checkout-form');
    if (form) {
        form.addEventListener('submit', handleFormSubmit);
    }
    
    // Lógica para seleção dos botões de pagamento
    const btnPix = document.getElementById('btn-pix');
    const btnBoleto = document.getElementById('btn-boleto');
    if (btnPix && btnBoleto) {
        btnPix.addEventListener('click', () => selectPagamento('pix'));
        btnBoleto.addEventListener('click', () => selectPagamento('boleto'));
    }
});


//
// FUNÇÕES PRINCIPAIS
//

function montarResumoPedido(itens) {
    const pedidoResumoDiv = document.getElementById('pedido-resumo');
    if (!pedidoResumoDiv) return;
    
    const cupom = (localStorage.getItem('cupom') || '').toUpperCase();
    let total = 0;
    
    // Cria a lista de itens
    let itensHtml = itens.map(item => {
        total += item.preco * (item.quantidade || 1);
        return `<li>${item.nome} (${item.quantidade || 1}x) - <strong>R$ ${item.preco.toFixed(2)}</strong></li>`;
    }).join('');

    let descontoHtml = '';
    let totalComDesconto = total;

    // Aplica o desconto se o cupom for válido
    if (cupom === "HAD10") {
        const desconto = total * 0.10;
        totalComDesconto -= desconto;
        descontoHtml = `<div class="pedido-desconto">Cupom HAD10: -R$ ${desconto.toFixed(2)}</div>`;
    }

    // Monta o HTML final do resumo
    pedidoResumoDiv.innerHTML = `
        <strong>Resumo do pedido:</strong>
        <ul>${itensHtml}</ul>
        ${descontoHtml}
        <div class="pedido-total"><strong>Total:</strong> R$ ${totalComDesconto.toFixed(2)}</div>
    `;
}

async function handleFormSubmit(event) {
    event.preventDefault();

    // 1. Valida o formulário
    if (!validarFormulario()) {
        alert('Por favor, corrija os campos inválidos.');
        return;
    }

    // 2. Monta o objeto do pedido
    const pedido = montarObjetoPedido();

    // 3. Envia o pedido para a API
    showLoader();
    try {
        const data = await fetchApi('/pedido', {
            method: 'POST',
            body: JSON.stringify(pedido)
        });

        // 4. Se o pedido for bem-sucedido
        localStorage.removeItem('carrinho'); // Limpa o carrinho
        localStorage.removeItem('cupom'); // Limpa o cupom
        
        exibirCodigoPagamento(pedido.pagamento);

    } catch (error) {
        // 5. Trata os erros
        if (error.message.includes('Cupom HAD10 já utilizado')) {
            exibirModalCompraNegada();
        } else {
            console.error('Erro ao finalizar pedido:', error);
            alert(`Erro ao finalizar pedido: ${error.message}`);
        }
    } finally {
        hideLoader();
    }
}

function montarObjetoPedido() {
    const carrinho = JSON.parse(localStorage.getItem('carrinho') || '[]');
    const cupom = localStorage.getItem('cupom') || "";

    const itensPedido = carrinho.map(item => ({
        id: item.id,
        nome: item.nome,
        price: item.preco, // O backend espera 'price', não 'preco'
        quantidade: item.quantidade
    }));

    return {
        nome: document.getElementById('nome').value,
        cpf: document.getElementById('cpf').value,
        email: document.getElementById('email').value,
        telefone: document.getElementById('telefone').value,
        endereco: document.getElementById('endereco').value,
        numero: document.getElementById('numero').value,
        bairro: document.getElementById('bairro').value,
        cep: document.getElementById('cep').value,
        cidade: document.getElementById('cidade').value,
        uf: document.getElementById('uf').value,
        pagamento: document.getElementById('pagamento').value,
        itens: itensPedido,
        cupom: cupom,
    };
}


//
// FUNÇÕES DE VALIDAÇÃO E UI
//

function validarFormulario() {
    let isValid = true;
    document.querySelectorAll('.input-error').forEach(span => span.textContent = '');
    document.querySelectorAll('.input-invalid').forEach(input => input.classList.remove('input-invalid'));

    // Valida cada campo e atualiza 'isValid'
    isValid = validaCampo('nome', /.{3,}/, 'Digite seu nome completo') && isValid;
    isValid = validaCampo('cpf', /^\d{3}\.\d{3}\.\d{3}-\d{2}$/, 'CPF inválido. Formato: 000.000.000-00') && isValid;
    isValid = validaCampo('email', /^[^@\s]+@[^@\s]+\.[^@\s]+$/, 'Digite um e-mail válido') && isValid;
    isValid = validaCampo('telefone', /^\(\d{2}\) \d{5}-\d{4}$/, 'Telefone inválido. Formato: (99) 99999-9999') && isValid;
    isValid = validaCampo('endereco', /.+/, 'Preencha a rua') && isValid;
    isValid = validaCampo('numero', /^\d+$/, 'Digite apenas números') && isValid;
    isValid = validaCampo('bairro', /.+/, 'Preencha o bairro') && isValid;
    isValid = validaCampo('cep', /^\d{5}-\d{3}$/, 'CEP inválido. Formato: 00000-000') && isValid;
    isValid = validaCampo('cidade', /.+/, 'Preencha a cidade') && isValid;
    isValid = validaCampo('uf', /^[A-Za-z]{2}$/, 'Digite a sigla do estado (ex: SP)') && isValid;
    isValid = validaCampo('pagamento', /.+/, 'Selecione a forma de pagamento') && isValid;

    return isValid;
}

function validaCampo(id, regex, mensagemErro) {
    const campo = document.getElementById(id);
    const erroEl = document.getElementById(`erro-${id}`);

    if (!campo.value.match(regex)) {
        campo.classList.add('input-invalid');
        if (erroEl) erroEl.textContent = mensagemErro;
        return false;
    }
    return true;
}

function selectPagamento(tipo) {
    document.getElementById('btn-pix').classList.remove('selected');
    document.getElementById('btn-boleto').classList.remove('selected');
    
    if (tipo === 'pix') document.getElementById('btn-pix').classList.add('selected');
    if (tipo === 'boleto') document.getElementById('btn-boleto').classList.add('selected');

    document.getElementById('pagamento').value = tipo;
    document.getElementById('erro-pagamento').textContent = '';
}

function exibirCodigoPagamento(metodo) {
    const titulo = `Pagamento ${metodo === "pix" ? "PIX" : "Boleto"}`;
    const codigo = metodo === "pix"
        ? "00020126580014BR.GOV.BCB.PIX0136b1e1f2e3d4c5b6a7f8e9d0c1b2a3f4g5h6i7j8k9l5204000053039865405120.005802BR5920Had Store6009SAO PAULO61080540900062070503***6304ABCD"
        : "34191.79001 01043.510047 91020.150008 7 92180011000";

    const bodyHtml = `
        <p>Utilize o código abaixo para realizar o pagamento:</p>
        <pre id="codigo-pagamento" class="confirmacao-codigo">${codigo}</pre>
        <button id="copiar-codigo" class="auth-button" style="margin-bottom: 1rem;">Copiar código</button>
        <button id="fechar-modal-pagamento" class="auth-button" style="background: #444;">Fechar</button>
    `;
    openModal(titulo, bodyHtml);

    document.getElementById('copiar-codigo').addEventListener('click', function() {
        navigator.clipboard.writeText(codigo).then(() => {
            this.innerText = "Copiado!";
            setTimeout(() => { this.innerText = "Copiar código"; }, 2000);
        });
    });
    document.getElementById('fechar-modal-pagamento').addEventListener('click', () => {
        closeModal();
        window.location.href = 'index.html'; // Redireciona para home após fechar
    });
}

function exibirModalCompraNegada() {
    const bodyHtml = `
        <p style="color:#e11d48;font-weight:bold;">O cupom HAD10 só pode ser usado uma vez por CPF.</p>
        <p>Por favor, remova o cupom no carrinho e tente finalizar sua compra novamente.</p>
        <button id="fechar-modal-negada" class="auth-button">Entendi</button>
    `;
    openModal("Compra negada", bodyHtml);
    document.getElementById('fechar-modal-negada').addEventListener('click', closeModal);
}

function aplicarMascaras() {
    if (window.jQuery && window.jQuery.fn.mask) {
        $('#cpf').mask('000.000.000-00');
        $('#telefone').mask('(00) 00000-0000');
        $('#cep').mask('00000-000');
        $('#numero').mask('000000');
    }
}