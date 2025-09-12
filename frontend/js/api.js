// Detecta ambiente de desenvolvimento e de teste
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

// A função que vamos usar em todos os lugares para fazer chamadas à API
export async function fetchApi(path, options = {}) {
    // Pega o token de login
    const token = localStorage.getItem('jwt_token');
    const userEmail = localStorage.getItem('user_email');
    let devUuid;
    // Só solicita/enviará o UUID se estiver em ambiente de desenvolvimento
    const isDev = window.DEV_ENV === true;
    const isTest = typeof process !== 'undefined' && process.env && process.env.JEST_WORKER_ID;
    if (isDev && !isTest) {
        devUuid = localStorage.getItem('dev_uuid');
        if (!devUuid) {
            devUuid = prompt('Digite o UUID de acesso ao ambiente de desenvolvimento:');
            if (devUuid) {
                localStorage.setItem('dev_uuid', devUuid);
            }
        }
    }

    const defaultHeaders = {
        'Content-Type': 'application/json',
    };

    // Se o token existir, adiciona ao cabeçalho de Autorização
    if (token) {
        defaultHeaders['Authorization'] = `Bearer ${token}`;
    }
    if (userEmail) {
        defaultHeaders['X-User-Email'] = userEmail;
    }
    if (isDev && devUuid) {
        defaultHeaders['X-Dev-UUID'] = devUuid;
    }

    // Combina os cabeçalhos padrão com quaisquer outros que a chamada específica precise
    options.headers = { ...defaultHeaders, ...options.headers };

    let url = path;
    if (!url.startsWith('/auth')) {
        url = `/api${path}`;
    }
    const response = await fetch(url, options);

    if (response.status === 401) {
        // Se o token for inválido ou expirado, limpa e redireciona para o login
        localStorage.removeItem('jwt_token');
        if (!isTest) {
            window.location.href = 'auth.html';
        }
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