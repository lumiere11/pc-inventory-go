# 🖥️ PC Inventory Management System

Un sistema backend completo para gestión de inventario de componentes de PC desarrollado en **Go** con **Gin Framework**, **GORM**, **MySQL** y **Casbin** para autorización.

## 📋 ¿Qué hace este Backend?

Este sistema proporciona una **API REST completa** para gestionar un inventario de productos de PC con las siguientes funcionalidades:

### 🔐 **Sistema de Autenticación**
- **Registro de usuarios** con roles (admin/user)
- **Login con JWT** (JSON Web Tokens)
- **Autorización basada en roles** usando Casbin RBAC
- **Contraseñas encriptadas** con bcrypt

### 📦 **Gestión de Productos**
- **CRUD completo** de productos (Crear, Leer, Actualizar, Eliminar)
- **Categorización** de productos (Periféricos, Monitores, etc.)
- **Estados de inventario** (stock, sold out)
- **Actualización independiente de stock**
- **Campos detallados**: nombre, marca, modelo, descripción, precio, stock

### 🔍 **Búsqueda Inteligente**
- **Búsqueda pública** (sin necesidad de autenticación)
- **Búsqueda por múltiples campos** (nombre, marca, modelo, descripción)
- **Búsqueda fuzzy** con tolerancia a errores tipográficos
- **Filtrado por estado** (stock/sold out)
- **Búsqueda case-insensitive**

### 🏗️ **Arquitectura Modular**
- **Separación de responsabilidades** (handlers, models, routes, database)
- **Middleware de autenticación y autorización**
- **Configuración por variables de entorno**
- **Migraciones automáticas de base de datos**
- **Seeders con datos iniciales**

## 🚀 Cómo Ejecutar el Backend

### **Prerequisitos**
- [Go](https://go.dev/dl/) 1.20+ instalado
- **MySQL** ejecutándose (XAMPP, MySQL Server, etc.)
- Git instalado

### **Paso 1: Clonar el Repositorio**
```bash
git clone https://github.com/lumiere11/pc-inventory-go
cd pc-inventory-go
```

### **Paso 2: Configurar Base de Datos**

**Opción A: Usando XAMPP**
1. Inicia **XAMPP Control Panel**
2. Inicia el servicio **MySQL**
3. Abre **phpMyAdmin** (http://localhost/phpmyadmin)
4. Crea una base de datos llamada `pc_inventory`

**Opción B: MySQL Server**
```sql
CREATE DATABASE pc_inventory CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

### **Paso 3: Configurar Variables de Entorno**
Crea un archivo `.env` (opcional, usa valores por defecto si no existe):
```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=pc_inventory
GIN_MODE=debug
```

### **Paso 4: Instalar Dependencias**
```bash
go mod tidy
```

### **Paso 5: Ejecutar el Servidor**
```bash
go run cmd/server/main.go
```

El servidor se ejecutará en **http://localhost:8081**

### **Paso 6: Verificar que Funciona**
```bash
# Probar búsqueda pública
curl "http://localhost:8081/api/v1/products/search?q=mouse"

# Debería devolver productos con "mouse" en el nombre
```

## 👥 Usuarios por Defecto

El sistema crea automáticamente los siguientes datos iniciales:

### **🔑 Usuario Administrador**
- **Email**: `admin@admin.com`
- **Contraseña**: `password`
- **Rol**: `admin`
- **Permisos**: Crear, actualizar y eliminar productos

### **📊 Categorías Predefinidas**
- Periféricos
- Monitores  
- Gabinetes
- Procesadores
- Tarjetas Gráficas
- Memoria RAM
- Placas Madre
- Almacenamiento
- Fuentes de Poder
- Refrigeración
- Tarjetas de Sonido
- Tarjetas de Red
- Lectores Ópticos
- Cables y Conectores
- Ventiladores

### **📋 Estados de Inventario**
- **stock**: Productos disponibles
- **sold out**: Productos agotados

### **🖱️ Productos de Ejemplo**
- **Razer DeathAdder V3 Mouse** - Gaming mouse ergonómico
- **Razer BlackWidow V4** - Teclado mecánico gaming
- **Logitech G502 Mouse** - Mouse gaming de alto rendimiento  
- **Corsair M65 RGB Elite Mouse** - Mouse FPS con botón sniper

## 📡 Endpoints de la API

### **🌐 Endpoints Públicos (Sin Autenticación)**
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/register` | Registrar nuevo usuario |
| POST | `/api/v1/login` | Iniciar sesión |
| GET | `/api/v1/products/search` | Buscar productos |

### **🔒 Endpoints Protegidos (Requieren Autenticación)**
| Método | Endpoint | Descripción | Rol Requerido |
|--------|----------|-------------|---------------|
| POST | `/api/v1/products` | Crear producto | admin |
| PUT | `/api/v1/products/:id` | Actualizar producto | admin |
| PUT | `/api/v1/products/:id/stock` | Actualizar stock | admin |
| DELETE | `/api/v1/products/:id` | Eliminar producto | admin |

## 🧪 Ejemplos de Uso

### **1. Registrar Usuario**
```bash
curl -X POST http://localhost:8081/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "mipassword",
    "role": "user"
  }'
```

### **2. Iniciar Sesión**
```bash
curl -X POST http://localhost:8081/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@admin.com",
    "password": "password"
  }'
```

**Respuesta:**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@admin.com",
    "role": "admin"
  }
}
```

### **3. Buscar Productos**
```bash
# Búsqueda exacta
curl "http://localhost:8081/api/v1/products/search?q=mouse"

# Búsqueda con errores tipográficos (fuzzy search)
curl "http://localhost:8081/api/v1/products/search?q=mou3se"

# Filtrar por estado
curl "http://localhost:8081/api/v1/products/search?q=mouse&status=stock"
```

### **4. Crear Producto (Requiere Token)**
```bash
curl -X POST http://localhost:8081/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TU_TOKEN_JWT_AQUI" \
  -d '{
    "name": "Nuevo Mouse Gaming",
    "brand": "Logitech",
    "model2": "G Pro X",
    "description": "Mouse gaming profesional",
    "stock": "25",
    "price": "99.99",
    "category_id": "1",
    "status_id": "1"
  }'
```

## 🔧 Configuración Avanzada

### **Variables de Entorno Disponibles**
```env
# Base de datos
DB_HOST=127.0.0.1        # Host de MySQL
DB_PORT=3306             # Puerto de MySQL  
DB_USER=root             # Usuario de MySQL
DB_PASSWORD=             # Contraseña de MySQL
DB_NAME=pc_inventory     # Nombre de la base de datos

# Aplicación
GIN_MODE=debug           # Modo Gin (debug/release)
PORT=8081               # Puerto del servidor
```

### **Estructura de la Base de Datos**
- **users**: Usuarios del sistema
- **categories**: Categorías de productos
- **statuses**: Estados de inventario
- **products**: Productos del inventario

## 🛠️ Desarrollo

### **Ejecutar Tests**
```bash
go test ./...
```

### **Compilar para Producción**
```bash
# Compilar binario
go build -o pc-inventory cmd/server/main.go

# Ejecutar en modo release
GIN_MODE=release ./pc-inventory
```

### **Estructura del Proyecto**
```
pc-inventory/
├── cmd/server/main.go      # Punto de entrada
├── database/
│   ├── database.go         # Configuración de DB
│   └── seeders.go          # Datos iniciales
├── handlers/               # Controladores HTTP
├── middlewares/            # Middleware de autenticación
├── models/                 # Modelos de datos
├── requests/               # Estructuras de validación
├── routes/                 # Configuración de rutas
├── model.conf             # Configuración Casbin
├── policy.csv             # Políticas RBAC
└── .env                   # Variables de entorno
```

## 🚨 Troubleshooting

### **Error: Connection Refused**
- Verifica que MySQL esté ejecutándose
- Revisa las variables de entorno de conexión
- Asegúrate de que la base de datos `pc_inventory` exista

### **Error: Login 401 Unauthorized**
- Verifica que el usuario `admin@admin.com` exista en la base de datos
- Confirma que los seeders se ejecutaron correctamente
- Revisa que la contraseña sea `password`

### **Error: Products Not Found**
- Verifica que los productos de ejemplo se crearon
- Revisa los logs del servidor para errores de seeding
- Confirma que las categorías y estados existen

---

**¿Necesitas ayuda?** Abre un issue en GitHub o revisa los logs del servidor para más detalles sobre errores específicos.