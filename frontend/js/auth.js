document.addEventListener('DOMContentLoaded', () => {
    // Seleciona os elementos do DOM
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterButton = document.getElementById('show-register-form');
    const showLoginButton = document.getElementById('show-login-form');
    
    const loginMessageEl = document.getElementById('login-message');
    const registerMessageEl = document.getElementById('register-message');

    // Função para alternar a visibilidade dos formulários
    const toggleForms = () => {
        loginForm.hidden = !loginForm.hidden;
        registerForm.hidden = !registerForm.hidden;
    };

    // Event listeners para os botões de alternância
    showRegisterButton.addEventListener('click', toggleForms);
    showLoginButton.addEventListener('click', toggleForms);

    // Função para exibir mensagens de feedback de forma acessível
    const showMessage = (element, message, isError = true) => {
        element.textContent = message;
        element.className = 'auth-message'; // Reset class
        element.classList.add(isError ? 'error' : 'success');
        element.style.display = 'block';
    };

    // Event listener para o formulário de REGISTRO
    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault(); // Impede o envio padrão do formulário
        
        const email = document.getElementById('register-email').value;
        const password = document.getElementById('register-password').value;

        try {
            const response = await fetch('/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                showMessage(registerMessageEl, 'Registro bem-sucedido! Você já pode fazer o login.', false);
                // Opcional: alternar para o formulário de login após sucesso
                setTimeout(() => {
                    toggleForms();
                    document.getElementById('login-email').value = email; // Preenche o email para facilitar
                    document.getElementById('login-password').focus();
                }, 2000);
            } else {
                const errorData = await response.text();
                showMessage(registerMessageEl, `Erro no registro: ${errorData}`);
            }
        } catch (error) {
            showMessage(registerMessageEl, 'Erro de conexão. Tente novamente.');
        }
    });

    // Event listener para o formulário de LOGIN
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        const rememberMe = document.getElementById('remember-me').checked;

        try {
            const response = await fetch('/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password, rememberMe }),
            });

            const data = await response.json();

            if (response.ok) {
                // 'Salva o token e o email
                localStorage.setItem('jwt_token', data.token);
                localStorage.setItem('user_email', email);

                showMessage(loginMessageEl, 'Login bem-sucedido! Redirecionando...', false);

                // Redireciona para a página principal após um breve momento
                setTimeout(() => {
                    window.location.href = 'index.html';
                }, 1500);

            } else {
                // Erro de login (senha errada, bloqueado, etc.)
                showMessage(loginMessageEl, data.error || 'Email ou senha inválidos.');
            }
        } catch (error) {
            showMessage(loginMessageEl, 'Erro de conexão ou dados inválidos. Tente novamente.');
        }
    });
});