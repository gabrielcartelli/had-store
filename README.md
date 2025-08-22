# Had Store

Este é um projeto de loja de chapéus desenvolvido em Go (backend) e JavaScript/HTML/CSS (frontend).

## Como rodar localmente

1. **Clone o repositório**
   ```sh
   git clone <url-do-repositorio>
   ```

2. **Backend**
   - Entre na pasta `backend`
   - Instale as dependências:
     ```sh
     go mod tidy
     ```
   - Rode o servidor:
     ```sh
     go run main.go
     ```
   - O backend roda na porta `8080` por padrão.

3. **Frontend**
   - Os arquivos estão na pasta `frontend`.
   - Abra o arquivo `index.html` em seu navegador para testar a interface.

## Deploy no Fly.io

1. Instale o Fly CLI:
   ```sh
   iwr https://fly.io/install.ps1 -useb | iex
   ```
2. Configure o app:
   ```sh
   fly launch
   ```
3. Faça o deploy:
   ```sh
   fly deploy
   ```

## Estrutura do projeto

```
backend/      # Código Go do servidor e API
frontend/     # Código do site (HTML, CSS, JS, imagens)
Dockerfile    # Configuração para deploy containerizado
fly.toml      # Configuração do Fly.io
```

## Funcionalidades

- Listagem de chapéus
- Pesquisa de chapéus
- Carrinho individual por usuário (localStorage)
- Menu lateral (hambúrguer) para o carrinho
- Finalização de compra com registro do pedido
- Suporte a PIX e Boleto (códigos fictícios)

## Documentação da API (Swagger)

A API possui documentação automática via Swagger.

- Gere a documentação com:
  ```sh
  swag init
  ```
- Acesse em produção:
  ```
  https://had-store.fly.dev/swagger/
  ```
- Os endpoints principais estão descritos e podem ser testados diretamente pela interface.

## Como contribuir

1. Faça um fork do projeto
2. Crie uma branch (`git checkout -b minha-feature`)
3. Faça suas alterações
4. Envie um pull request

---

Had