
import { fetchApi } from './api.js';

// Prompt do UUID também em auth.html
const isDev = window.DEV_ENV === true;
const isTest = typeof process !== 'undefined' && process.env && process.env.JEST_WORKER_ID;
if (isDev && !isTest) {
    if (!localStorage.getItem('dev_uuid')) {
        const uuid = prompt('Digite o token de acesso ao ambiente de desenvolvimento:');
        if (uuid) {
            localStorage.setItem('dev_uuid', uuid);
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    // Seleciona os elementos do DOM
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterButton = document.getElementById('show-register-form');
    const showLoginButton = document.getElementById('show-login-form');
    
    const loginMessageEl = document.getElementById('login-message');
    const registerMessageEl = document.getElementById('register-message');

    // Função para alternar a visibilidade dos formulários usando display
    const toggleForms = () => {
        if (loginForm.style.display === 'none') {
            loginForm.style.display = 'block';
            registerForm.style.display = 'none';
        } else {
            loginForm.style.display = 'none';
            registerForm.style.display = 'block';
        }
    };
    window.toggleForms = toggleForms;

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
            await fetchApi('/auth/register', {
                method: 'POST',
                body: JSON.stringify({ email, password }),
            });

            showMessage(registerMessageEl, 'Registro bem-sucedido! Você já pode fazer o login.', false);
            // Opcional: alternar para o formulário de login após sucesso
            setTimeout(() => {
                toggleForms();
                document.getElementById('login-email').value = email; // Preenche o email para facilitar
                document.getElementById('login-password').focus();
            }, 2000);
        } catch (error) {
            showMessage(registerMessageEl, `Erro no registro: ${error.message}`);
        }
    });

    // Event listener para o formulário de LOGIN
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const email = document.getElementById('login-email').value;
        const password = document.getElementById('login-password').value;
        const rememberMe = document.getElementById('remember-me').checked;

        try {
            const data = await fetchApi('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password, rememberMe }),
            });

            // 'Salva o token e o email
            localStorage.setItem('jwt_token', data.token);
            localStorage.setItem('user_email', email);

            showMessage(loginMessageEl, 'Login bem-sucedido! Redirecionando...', false);

            // Redireciona para a página principal após um breve momento
            setTimeout(() => {
                window.location.href = 'index.html';
            }, 1500);

        } catch (error) {
            showMessage(loginMessageEl, `Erro no login: ${error.message || 'Credenciais inválidas.'}`);
        }
    });
});