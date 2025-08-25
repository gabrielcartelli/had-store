document.addEventListener('DOMContentLoaded', function() {
    const API_UUID = "e3e6c6c2-9b7d-4c5e-8c1a-2f7b8f8e2a1d";

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

    function fetchComUUID(url, options = {}) {
        if (!options.headers) options.headers = {};
        options.headers["X-API-UUID"] = localStorage.getItem("api_uuid");
        return fetch(url, options);
    }

    const form = document.getElementById('checkout-form');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault(); // CORREÇÃO: impede o submit padrão do formulário
            let valid = true;

            // Limpa erros anteriores
            document.querySelectorAll('.input-error').forEach(span => span.textContent = '');
            document.querySelectorAll('.input-invalid').forEach(input => input.classList.remove('input-invalid'));

            // Dados pessoais
            const nome = document.getElementById('nome');
            if (!nome.value || nome.value.length < 3) {
                nome.classList.add('input-invalid');
                document.getElementById('erro-nome').textContent = 'Digite seu nome completo';
                valid = false;
            }

            const cpf = document.getElementById('cpf');
            if (!cpf.value.match(/^\d{3}\.\d{3}\.\d{3}-\d{2}$/)) {
                cpf.classList.add('input-invalid');
                document.getElementById('erro-cpf').textContent = 'CPF inválido. Formato: 000.000.000-00';
                valid = false;
            }

            const email = document.getElementById('email');
            if (!email.value.match(/^[^@\s]+@[^@\s]+\.[^@\s]+$/)) {
                email.classList.add('input-invalid');
                document.getElementById('erro-email').textContent = 'Digite um e-mail válido';
                valid = false;
            }

            const telefone = document.getElementById('telefone');
            if (!telefone.value.match(/^\(\d{2}\) \d{5}-\d{4}$/)) {
                telefone.classList.add('input-invalid');
                document.getElementById('erro-telefone').textContent = 'Telefone inválido. Formato: (99) 99999-9999';
                valid = false;
            }

            // Endereço
            const endereco = document.getElementById('endereco');
            if (!endereco.value) {
                endereco.classList.add('input-invalid');
                document.getElementById('erro-endereco').textContent = 'Preencha a rua';
                valid = false;
            }

            const numero = document.getElementById('numero');
            if (!numero.value.match(/^\d+$/)) {
                numero.classList.add('input-invalid');
                document.getElementById('erro-numero').textContent = 'Digite apenas números';
                valid = false;
            }

            const bairro = document.getElementById('bairro');
            if (!bairro.value) {
                bairro.classList.add('input-invalid');
                document.getElementById('erro-bairro').textContent = 'Preencha o bairro';
                valid = false;
            }

            const cep = document.getElementById('cep');
            if (!cep.value.match(/^\d{5}-\d{3}$/)) {
                cep.classList.add('input-invalid');
                document.getElementById('erro-cep').textContent = 'CEP inválido. Formato: 00000-000';
                valid = false;
            }

            const cidade = document.getElementById('cidade');
            if (!cidade.value) {
                cidade.classList.add('input-invalid');
                document.getElementById('erro-cidade').textContent = 'Preencha a cidade';
                valid = false;
            }

            const uf = document.getElementById('uf');
            if (!uf.value.match(/^[A-Za-z]{2}$/)) {
                uf.classList.add('input-invalid');
                document.getElementById('erro-uf').textContent = 'Digite a sigla do estado (ex: SP)';
                valid = false;
            }

            // Pagamento
            const pagamentoInput = document.getElementById('pagamento');
            if (!pagamentoInput.value) {
                document.getElementById('erro-pagamento').textContent = 'Selecione a forma de pagamento.';
                valid = false;
            }

            if (!valid) {
                // Não exibe alert, apenas mensagens nos inputs
                return;
            }

            const nomeValue = nome.value;
            const cpfValue = cpf.value;
            const emailValue = email.value;
            const pagamentoValue = pagamentoInput.value;

            const cupomInput = document.getElementById('cupom');
            const cupomValue = localStorage.getItem('cupom') || "";

            const pedidoMemoria = window.localStorage.getItem('carrinho');
            const carrinho = pedidoMemoria ? JSON.parse(pedidoMemoria) : [];

            // Garante que os itens do pedido têm os campos esperados pelo backend
            const itensPedido = carrinho.map(item => ({
                id: item.id,
                nome: item.nome,
                preco: item.preco,
                quantidade: item.quantidade
            }));

            const pedido = {
                nome: nomeValue,
                cpf: cpfValue,
                email: emailValue,
                pagamento: pagamentoValue,
                itens: itensPedido,
                total: itensPedido.reduce((acc, item) => acc + item.preco * item.quantidade, 0),
                cupom: cupomValue // Envia o cupom para o backend
            };

            fetchComUUID('https://hat-store-training.fly.dev/api/pedido', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(pedido)
            })
            .then(res => {
                if (res.status === 401) {
                    alert("Acesso não autorizado. Código de acesso inválido.");
                    localStorage.removeItem("api_uuid");
                    solicitarUUID();
                    return;
                }
                if (res.status === 403) {
                    // Compra negada por uso repetido do cupom HAD10
                    exibirModalCompraNegada();
                    return;
                }
                return res.json();
            })
            .then(data => {
                if (data) {
                    localStorage.removeItem('carrinho');
                    exibirCodigoPagamento(pagamentoValue);
                }
            })
            .catch(() => {
                alert("Erro ao finalizar pedido.");
            });
        });
    }

    // Função para modal de compra negada
    function exibirModalCompraNegada() {
        const modal = document.getElementById('modal');
        const modalTitle = document.getElementById('modal-title');
        const modalBody = document.getElementById('modal-body');
        if (modal && modalTitle && modalBody) {
            modalTitle.innerText = "Compra negada";
            modalBody.innerHTML = `
                <p style="color:#e11d48;font-weight:bold;">O cupom HAD10 só pode ser usado uma vez por CPF.</p>
                <p>Se você já utilizou esse cupom, finalize sua compra sem o cupom ou entre em contato com o suporte.</p>
                <button id="fechar-modal-negada" style="background:#232323;color:#c9a94b;border:none;border-radius:8px;padding:0.7rem 1.2rem;font-size:1rem;cursor:pointer;">Fechar</button>
            `;
            modal.style.display = 'flex';
            document.getElementById('fechar-modal-negada').onclick = function() {
                modal.style.display = 'none';
            };
        }
    }

    // Usa a modal já existente no HTML e não redireciona automaticamente
    function exibirCodigoPagamento(metodo) {
        let titulo = `Pagamento ${metodo === "pix" ? "PIX" : "Boleto"}`;
        let codigo = "";
        if (metodo === "pix") {
            codigo = "00020126580014BR.GOV.BCB.PIX0136b1e1f2e3d4c5b6a7f8e9d0c1b2a3f4g5h6i7j8k9l5204000053039865405120.005802BR5920Had Store6009SAO PAULO61080540900062070503***6304ABCD";
        } else if (metodo === "boleto") {
            codigo = "34191.79001 01043.510047 91020.150008 7 92180011000";
        } else {
            codigo = "Código não disponível";
        }

        // Preenche e exibe a modal já existente
        const modal = document.getElementById('modal');
        const modalTitle = document.getElementById('modal-title');
        const modalBody = document.getElementById('modal-body');
        if (modal && modalTitle && modalBody) {
            modalTitle.innerText = titulo;
            modalBody.innerHTML = `
                <p>Utilize o código abaixo para realizar o pagamento:</p>
                <pre id="codigo-pagamento" class="confirmacao-codigo">${codigo}</pre>
                <button id="copiar-codigo" style="background:#c9a94b;color:#fff;border:none;border-radius:8px;padding:0.7rem 1.2rem;font-size:1rem;cursor:pointer;margin-bottom:1rem;">Copiar código</button>
                <br>
                <button id="fechar-modal" style="background:#232323;color:#c9a94b;border:none;border-radius:8px;padding:0.7rem 1.2rem;font-size:1rem;cursor:pointer;">Fechar</button>
            `;
            modal.style.display = 'flex';

            document.getElementById('copiar-codigo').onclick = function() {
                const texto = document.getElementById('codigo-pagamento').innerText;
                navigator.clipboard.writeText(texto).then(() => {
                    this.innerText = "Copiado!";
                    setTimeout(() => { this.innerText = "Copiar código"; }, 1500);
                });
            };

            // CORREÇÃO: apenas fecha a modal, não redireciona
            document.getElementById('fechar-modal').onclick = function() {
                modal.style.display = 'none';
            };
        }
    }

    const pedidoResumo = document.getElementById('pedido-resumo');
    const carrinho = window.localStorage.getItem('carrinho');
    const cupomValue = (localStorage.getItem('cupom') || '').toUpperCase();

    if (pedidoResumo && carrinho) {
        const itens = JSON.parse(carrinho);
        let html = "<strong>Resumo do pedido:</strong><ul>";
        let total = 0;
        itens.forEach(item => {
            html += `<li>${item.nome} (${item.quantidade || 1}x) - <strong>R$ ${item.preco.toFixed(2)}</strong></li>`;
            total += item.preco * (item.quantidade || 1);
        });

        let desconto = 0;
        if (cupomValue === "HAD10") {
            desconto = total * 0.10;
            html += `<div class="pedido-desconto" style="color:#22c55e;font-weight:bold;margin-top:0.5em;">
                        Cupom HAD10 aplicado: -R$ ${desconto.toFixed(2)}
                    </div>`;
        }

        html += `<div class="pedido-total"><strong>Total com desconto:</strong> R$ ${(total - desconto).toFixed(2)}</div>`;
        pedidoResumo.innerHTML = html;
    }

    // Máscaras dos inputs
    if (window.jQuery && window.jQuery.fn.mask) {
        $('#cpf').mask('000.000.000-00');
        $('#telefone').mask('(00) 00000-0000');
        $('#cep').mask('00000-000');
        $('#numero').mask('00000');
    }

    // Seleção da forma de pagamento
    const btnPix = document.getElementById('btn-pix');
    const btnBoleto = document.getElementById('btn-boleto');
    const pagamentoInput = document.getElementById('pagamento');
    const buttons = [btnPix, btnBoleto];

    function selectPagamento(tipo) {
        buttons.forEach(btn => btn.classList.remove('selected'));
        if (tipo === 'pix') btnPix.classList.add('selected');
        if (tipo === 'boleto') btnBoleto.classList.add('selected');
        pagamentoInput.value = tipo;
        document.getElementById('erro-pagamento').textContent = '';
    }

    if (btnPix) btnPix.onclick = () => selectPagamento('pix');
    if (btnBoleto) btnBoleto.onclick = () => selectPagamento('boleto');
});