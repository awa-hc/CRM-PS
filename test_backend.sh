#!/bin/bash

# Script para probar el backend del CRM
# Prueba el endpoint de actualización de proyectos

echo "=== Probando Backend CRM ==="
echo "Endpoint: http://localhost:8080/api/v1/projects/2"
echo ""

# 0. Hacer login para obtener el token
echo "0. Haciendo login para obtener token:"
LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "horacio@gmail.com",
    "password": "test123"
  }')

echo "Respuesta del login: $LOGIN_RESPONSE"

# Extraer el token de la respuesta
TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "Error: No se pudo obtener el token de autenticación"
  echo "Respuesta completa: $LOGIN_RESPONSE"
  exit 1
fi

echo "Token obtenido: $TOKEN"
echo ""
echo "==========================================="
echo ""

# 1. Primero obtener el proyecto actual
echo "1. Obteniendo proyecto actual (GET):"
curl -X GET "http://localhost:8080/api/v1/projects/2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n" \
  | jq .

echo ""
echo "==========================================="
echo ""

# 2. Actualizar el proyecto con datos de prueba
echo "2. Actualizando proyecto (PUT):"
curl -X PUT "http://localhost:8080/api/v1/projects/2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Proyecto Actualizado Test",
    "description": "Descripción actualizada desde curl",
    "client_id": 1,
    "status": "active",
    "priority": "high",
    "type": "construction",
    "address": "123 Test Street",
    "city": "Test City",
    "state": "Test State",
    "zip_code": "12345",
    "start_date": "2024-01-15T00:00:00Z",
    "end_date": "2024-06-15T00:00:00Z",
    "budget": 50000.00,
    "estimated_cost": 45000.00,
    "actual_cost": 0.00,
    "progress": 25,
    "notes": "Notas de prueba desde curl"
  }' \
  -w "\nStatus: %{http_code}\n" \
  | jq .

echo ""
echo "==========================================="
echo ""

# 3. Verificar que los cambios se guardaron
echo "3. Verificando cambios guardados (GET):"
curl -X GET "http://localhost:8080/api/v1/projects/2" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nStatus: %{http_code}\n" \
  | jq .

echo ""
echo "=== Prueba completada ==="