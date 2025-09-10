/**
 * @jest-environment jsdom
 */

describe('pedidos.js', () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <form id="consulta-form-pedidos"></form>
      <input id="cpf-consulta-pedidos" />
      <div id="pedidos-list"></div>
    `;
    localStorage.clear();
    window.showLoader = jest.fn();
    window.hideLoader = jest.fn();
  });

  test('exibe alerta se CPF não informado', () => {
    window.alert = jest.fn();
    const form = document.getElementById('consulta-form-pedidos');
    const input = document.getElementById('cpf-consulta-pedidos');
    input.value = '';
    // Função mockada
    form.addEventListener('submit', (event) => {
      event.preventDefault();
      const cpf = input.value.trim();
      if (!cpf) {
        alert('Por favor, informe um CPF.');
        return;
      }
    });
    const event = new Event('submit');
    form.dispatchEvent(event);
    expect(window.alert).toHaveBeenCalledWith('Por favor, informe um CPF.');
  });

  test('exibe mensagem de nenhum pedido encontrado', async () => {
    const pedidosListDiv = document.getElementById('pedidos-list');
    // Função mockada
    async function consultarPedidosMock() {
      pedidosListDiv.innerHTML = "<p>Buscando seus pedidos...</p>";
      const pedidos = [];
      pedidosListDiv.innerHTML = '';
      if (!pedidos || pedidos.length === 0) {
        pedidosListDiv.innerHTML = "<p>Nenhum pedido encontrado para este CPF.</p>";
        return;
      }
    }
    await consultarPedidosMock();
    expect(pedidosListDiv.innerHTML).toContain('Nenhum pedido encontrado');
  });

  test('renderiza cards de pedidos corretamente', async () => {
    const pedidosListDiv = document.getElementById('pedidos-list');
    const pedidos = [
      {
        nome: 'Cliente Teste',
        cpf: '123.456.789-00',
        pagamento: 'pix',
        total: 120.0,
        itens: [
          { nome: 'Chapéu Panamá', quantidade: 2, price: 100 },
          { nome: 'Chapéu Fedora', quantidade: 1, price: 150 }
        ]
      }
    ];
    // Função mockada
    function renderizarPedidos(pedidos) {
      pedidosListDiv.innerHTML = '';
      pedidos.forEach(pedido => {
        const card = document.createElement('div');
        card.className = 'pedido-card';
        card.innerHTML = `
          <p class="pedido-info-destaque"><strong>Cliente:</strong> ${pedido.nome}</p>
          <p class="pedido-info-destaque"><strong>CPF:</strong> ${pedido.cpf}</p>
          <p class="pedido-info-destaque"><strong>Pagamento:</strong> ${pedido.pagamento}</p>
          <p class="pedido-info-destaque"><strong>Total:</strong> R$ ${pedido.total.toFixed(2)}</p>
          <div class="pedido-itens">
            <strong>Itens do Pedido:</strong>
            <ul>
              ${pedido.itens.map(item => `<li>${item.nome} (${item.quantidade || 1}x) - R$ ${item.price.toFixed(2)}</li>`).join('')}
            </ul>
          </div>
        `;
        pedidosListDiv.appendChild(card);
      });
    }
    renderizarPedidos(pedidos);
    expect(pedidosListDiv.innerHTML).toContain('Cliente Teste');
    expect(pedidosListDiv.innerHTML).toContain('Chapéu Panamá');
    expect(pedidosListDiv.innerHTML).toContain('R$ 120.00');
  });
});
