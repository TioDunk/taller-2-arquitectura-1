# DDD POC - Domain Driven Design en Go

Este proyecto es una **Proof of Concept (POC)** que demuestra los conceptos principales de **Domain Driven Design (DDD)** usando Go.

## 🎯 Objetivo

Crear un ejemplo práctico y sencillo de DDD que permita entender y explicar los conceptos fundamentales de esta arquitectura.

## 📋 Dominio: Sistema de Gestión de Pedidos

Hemos elegido un **sistema de gestión de pedidos de una tienda online** porque permite demostrar claramente todos los conceptos de DDD:

- **Entidades**: Cliente, Pedido, Producto
- **Value Objects**: Dirección, Moneda, Cantidad
- **Aggregates**: Pedido (con sus items)
- **Domain Services**: Cálculo de precios, validaciones
- **Repositories**: Para persistencia
- **Application Services**: Casos de uso

## 🏗️ Estructura del Proyecto (DDD)

```
├── domain/                 # Capa de Dominio (Core)
│   ├── entities/          # Entidades del dominio
│   ├── valueobjects/      # Value Objects
│   ├── aggregates/        # Aggregates
│   ├── services/          # Domain Services
│   └── repositories/      # Interfaces de repositorios
├── application/           # Capa de Aplicación
│   ├── services/          # Application Services
│   └── dto/              # Data Transfer Objects
├── infrastructure/        # Capa de Infraestructura
│   ├── repositories/      # Implementaciones de repositorios
│   └── database/         # Base de datos (en memoria)
└── presentation/          # Capa de Presentación
    └── http/             # Controladores HTTP
```

## 🚀 Cómo Ejecutar

```bash
# Instalar dependencias
go mod tidy

# Ejecutar la aplicación
go run main.go
```

## 📚 Conceptos DDD Implementados

### 1. **Entidades (Entities)**
- Tienen identidad única
- Ejemplo: `Customer`, `Order`, `Product`

### 2. **Value Objects**
- No tienen identidad, se comparan por valor
- Ejemplo: `Address`, `Money`, `Quantity`

### 3. **Aggregates**
- Agrupan entidades y value objects relacionados
- Ejemplo: `Order` (con OrderItems)

### 4. **Domain Services**
- Lógica de dominio que no pertenece a una entidad específica
- Ejemplo: `PricingService`, `OrderValidationService`

### 5. **Repositories**
- Abstraen el acceso a datos
- Interfaces en el dominio, implementaciones en infraestructura

### 6. **Application Services**
- Orquestan el flujo de la aplicación
- Coordinan entre dominio e infraestructura

## 🔗 Endpoints Disponibles

- `POST /customers` - Crear cliente
- `GET /customers/:id` - Obtener cliente
- `POST /products` - Crear producto
- `POST /orders` - Crear pedido
- `GET /orders/:id` - Obtener pedido
- `POST /orders/:id/add-item` - Agregar item al pedido

## 💡 Conceptos Clave a Explicar

1. **Separación de Capas**: Dominio, Aplicación, Infraestructura, Presentación
2. **Inversión de Dependencias**: El dominio no depende de infraestructura
3. **Ubiquitous Language**: El código refleja el lenguaje del dominio
4. **Encapsulación**: La lógica de negocio está en el dominio
5. **Testabilidad**: Fácil de testear por la separación de responsabilidades