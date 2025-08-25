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
            e.preventDefault();
            const nome = document.getElementById('nome').value;
            const cpf = document.getElementById('cpf').value;
            const email = document.getElementById('email').value;
            const pagamento = document.getElementById('pagamento').value;

            const pedidoMemoria = window.localStorage.getItem('carrinho');
            const carrinho = pedidoMemoria ? JSON.parse(pedidoMemoria) : [];

            if (!nome || !cpf || !email || !pagamento || carrinho.length === 0) {
                alert("Preencha todos os campos e adicione itens ao carrinho.");
                return;
            }

            // Garante que os itens do pedido têm os campos esperados pelo backend
            const itensPedido = carrinho.map(item => ({
                id: item.id,
                nome: item.nome,
                preco: item.preco,
                quantidade: item.quantidade
            }));

            const pedido = {
                nome,
                cpf,
                email,
                pagamento,
                itens: itensPedido,
                total: itensPedido.reduce((acc, item) => acc + item.preco * item.quantidade, 0)
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
                return res.json();
            })
            .then(data => {
                if (data) {
                    localStorage.removeItem('carrinho');
                    exibirCodigoPagamento(pagamento);
                }
            })
            .catch(() => {
                alert("Erro ao finalizar pedido.");
            });
        });
    }

    function exibirCodigoPagamento(metodo) {
        // Gera um código fictício
        let codigo = "";
        if (metodo === "pix") {
            codigo = "00020126580014BR.GOV.BCB.PIX0136b1e1f2e3d4c5b6a7f8e9d0c1b2a3f4g5h6i7j8k9l5204000053039865405120.005802BR5920Had Store6009SAO PAULO61080540900062070503***6304ABCD";
        } else if (metodo === "boleto") {
            codigo = "34191.79001 01043.510047 91020.150008 7 92180011000";
        } else {
            codigo = "Código não disponível";
        }

        // Cria o modal
        const modal = document.createElement('div');
        modal.style.position = 'fixed';
        modal.style.top = '0';
        modal.style.left = '0';
        modal.style.width = '100vw';
        modal.style.height = '100vh';
        modal.style.background = 'rgba(0,0,0,0.35)';
        modal.style.display = 'flex';
        modal.style.alignItems = 'center';
        modal.style.justifyContent = 'center';
        modal.style.zIndex = '99999';

        modal.innerHTML = `
            <div style="background:#fff;padding:2rem 2.5rem;border-radius:14px;box-shadow:0 2px 16px #2563eb33;max-width:420px;text-align:center;">
                <h2 style="color:#2563eb;margin-bottom:1.2rem;">Pagamento ${metodo === "pix" ? "PIX" : "Boleto"}</h2>
                <p style="margin-bottom:0.7rem;">Utilize o código abaixo para realizar o pagamento:</p>
                <pre id="codigo-pagamento" style="background:#f3f6fa;padding:1rem;border-radius:8px;font-size:1.1rem;word-break:break-all;margin-bottom:1.2rem;">${codigo}</pre>
                <button id="copiar-codigo" style="background:#2563eb;color:#fff;border:none;border-radius:8px;padding:0.7rem 1.2rem;font-size:1rem;cursor:pointer;margin-bottom:1rem;">Copiar código</button>
                <br>
                <button id="fechar-modal" style="background:#e0e7ff;color:#2563eb;border:none;border-radius:8px;padding:0.7rem 1.2rem;font-size:1rem;cursor:pointer;">Fechar</button>
            </div>
        `;

        document.body.appendChild(modal);

        document.getElementById('copiar-codigo').onclick = function() {
            const texto = document.getElementById('codigo-pagamento').innerText;
            navigator.clipboard.writeText(texto).then(() => {
                this.innerText = "Copiado!";
                setTimeout(() => { this.innerText = "Copiar código"; }, 1500);
            });
        };

        document.getElementById('fechar-modal').onclick = function() {
            document.body.removeChild(modal);
            window.location.href = "index.html";
        };
    }
});