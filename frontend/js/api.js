// A função que vamos usar em todos os lugares para fazer chamadas à API
async function fetchApi(path, options = {}) {
    // Pega o UUID de acesso (lógica que já existia)
    const apiUUID = localStorage.getItem("api_uuid");
    if (!apiUUID) {
        // Você pode redirecionar para uma página de erro ou pedir o UUID novamente
        throw new Error("Código de acesso (UUID) não encontrado.");
    }
    
    // Pega o token de login
    const token = localStorage.getItem('jwt_token');

    const defaultHeaders = {
        'Content-Type': 'application/json',
        'X-API-UUID': apiUUID,
    };

    // Se o token existir, adiciona ao cabeçalho de Autorização
    if (token) {
        defaultHeaders['Authorization'] = `Bearer ${token}`;
    }

    // Combina os cabeçalhos padrão com quaisquer outros que a chamada específica precise
    options.headers = { ...defaultHeaders, ...options.headers };

    const response = await fetch(`/api${path}`, options);

    if (response.status === 401) {
        // Se o token for inválido ou expirado, limpa e redireciona para o login
        localStorage.removeItem('jwt_token');
        window.location.href = 'auth.html';
        throw new Error('Não autorizado');
    }

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Erro na API');
    }
    
    // Se a resposta não tiver corpo (ex: status 204), retorna um objeto vazio
    const contentType = response.headers.get("content-type");
    if (contentType && contentType.indexOf("application/json") !== -1) {
        return response.json();
    }
    return {};
}