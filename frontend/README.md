# Hat Store - Projeto Completo

Este repositório contém o projeto **Hat Store**, uma loja fictícia de chapéus composta por frontend (HTML/CSS/JS) e backend (Go). O objetivo é servir como base para treinamento, avaliação e evolução de funcionalidades.

## Especificações Técnicas

- **Frontend:**  
  - Desenvolvido em HTML, CSS e JavaScript puro.
  - Estrutura modular com páginas para listagem, checkout, autenticação, pedidos, documentação e trabalho avaliativo.
  - Integração com backend via chamadas REST.

- **Backend:**  
  - Implementado em Go.
  - API RESTful para autenticação, listagem de produtos, registro de pedidos, consulta de pedidos e controle de estoque.
  - Autenticação via JWT.
  - Documentação da API disponível via Swagger.

## Como rodar localmente

1. **Clone o repositório:**
   ```sh
   git clone https://github.com/seu-usuario/hat-store.git
   ```

2. **Backend:**
   - Acesse a pasta do backend.
   - Instale as dependências do Go.
   - Execute o servidor:
     ```sh
     go run main.go
     ```
   - O backend será iniciado em `http://localhost:8080` (ou porta configurada).

3. **Frontend:**
   - Acesse a pasta `frontend`.
   - Para testes locais, abra o arquivo `index.html` diretamente no navegador.
   - Para integração completa, utilize um servidor local (ex: [Live Server](https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer) no VS Code) para evitar restrições de CORS.

4. **Configuração de ambiente:**
   - Certifique-se de que o frontend está configurado para consumir a API do backend local (`http://localhost:8080`).
   - Ajuste variáveis de ambiente conforme necessário.

## Deploy Remoto

- O deploy do projeto é realizado na plataforma [Fly.io](https://fly.io/).
- O site está disponível em produção em:  
  [https://hat-store-training.fly.dev/](https://hat-store-training.fly.dev/)

## Releases

- **Release 1.0.0:**  
  Estrutura base do site, com listagem de chapéus, carrinho, checkout e autenticação.

- **Release 1.2.0:**  
  Adição da página de documentação do trabalho avaliativo, detalhando objetivos e instruções para testes.

- **Release 1.2.1:**  
  Atualização do escopo do trabalho avaliativo, com ajustes nas instruções e requisitos.

- **Release 1.3.0:**  
  Evoluções do sistema, incluindo controle de estoque dos produtos e filtros avançados na listagem de chapéus.

## Links Úteis

- [Documentação de negócio](https://hat-store-training.fly.dev/documentacao.html)
- [Trabalho avaliativo](trabalho-avaliativo.html)
- [Site da CWI](https://cwi.com.br/)
- [GitLab QA](https://git.cwi.com.br/formacoes/cwi-crescer/edicao-27/qa)
- [Azure Board](https://dev.azure.com/cwi-formacao/CWI%20Crescer/_sprints/taskboard/CWI%20Crescer%20Team/CWI%20Crescer/Iteration%201)

---

**Observação:**  
O Hat Store é uma aplicação fictícia, sem cobranças ou compras reais.  
Para dúvidas técnicas, consulte a documentação ou entre em contato com o responsável (Gabriel Cartelli - gabriel.cartelli@cwi.com.br)