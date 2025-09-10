#!/bin/bash
# Script para rodar testes do backend e frontend e só fazer deploy se todos passarem

# Rodar testes do backend (Go) apenas na pasta handlers
pushd backend/handlers > /dev/null
BACKEND_TEST_RESULT=$(go test)
if echo "$BACKEND_TEST_RESULT" | grep -q FAIL; then
    echo "Backend tests failed. Deployment aborted."
    popd > /dev/null
    exit 1
fi
popd > /dev/null


# Rodar testes do frontend (Jest)
pushd frontend > /dev/null
FRONTEND_TEST_RESULT=$(npm test -- --json)
popd > /dev/null

# Verificar resultado dos testes do frontend
if ! echo "$FRONTEND_TEST_RESULT" | grep '"success":true' > /dev/null; then
    echo "Frontend tests failed. Deployment aborted."
    exit 1
fi

# Se todos passaram, faz deploy
fly deploy
