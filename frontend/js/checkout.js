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
                    alert(data.mensagem || "Pedido realizado com sucesso!");
                    localStorage.removeItem('carrinho');
                    window.location.href = "index.html";
                }
            })
            .catch(() => {
                alert("Erro ao finalizar pedido.");
            });
        });
    }
});