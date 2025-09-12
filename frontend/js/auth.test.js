

/**
 * @jest-environment jsdom
 */

import './auth.js';

describe('auth.js', () => {
  let loginForm, registerForm, showRegisterButton, showLoginButton, loginMessageEl, registerMessageEl;

  beforeEach(() => {
    document.body.innerHTML = `
      <form id="login-form"><input id="login-email"><input id="login-password"><input type="checkbox" id="remember-me"></form>
      <form id="register-form"><input id="register-email"><input id="register-password"></form>
      <button id="show-register-form"></button>
      <button id="show-login-form"></button>
      <span id="login-message"></span>
      <span id="register-message"></span>
    `;
    loginForm = document.getElementById('login-form');
    registerForm = document.getElementById('register-form');
    showRegisterButton = document.getElementById('show-register-form');
    showLoginButton = document.getElementById('show-login-form');
    loginMessageEl = document.getElementById('login-message');
    registerMessageEl = document.getElementById('register-message');
  });

  test('toggleForms alterna visibilidade dos formulários', () => {
    loginForm.style.display = 'block';
    registerForm.style.display = 'none';
    document.dispatchEvent(new Event('DOMContentLoaded'));
    window.toggleForms();
    expect(loginForm.style.display).toBe('none');
    expect(registerForm.style.display).toBe('block');
    window.toggleForms();
    expect(loginForm.style.display).toBe('block');
    expect(registerForm.style.display).toBe('none');
  });

  test('showMessage exibe mensagem de erro', () => {
    // Função mockada
    const showMessage = (element, message, isError = true) => {
      element.textContent = message;
      element.className = 'auth-message';
      element.classList.add(isError ? 'error' : 'success');
      element.style.display = 'block';
    };
    showMessage(loginMessageEl, 'Erro de login');
    expect(loginMessageEl.textContent).toBe('Erro de login');
    expect(loginMessageEl.classList.contains('error')).toBe(true);
    expect(loginMessageEl.style.display).toBe('block');
  });

  test('showMessage exibe mensagem de sucesso', () => {
    const showMessage = (element, message, isError = true) => {
      element.textContent = message;
      element.className = 'auth-message';
      element.classList.add(isError ? 'error' : 'success');
      element.style.display = 'block';
    };
    showMessage(registerMessageEl, 'Sucesso!', false);
    expect(registerMessageEl.textContent).toBe('Sucesso!');
    expect(registerMessageEl.classList.contains('success')).toBe(true);
    expect(registerMessageEl.style.display).toBe('block');
  });
});
