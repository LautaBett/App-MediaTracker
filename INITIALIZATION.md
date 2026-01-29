
```markdown
# 🛠️ Guía de Inicialización y Ejecución

Esta guía detalla los pasos necesarios para configurar el entorno de desarrollo en macOS y ejecutar **MediaTracker**.

## 📋 1. Requisitos Previos (Solo la primera vez)

Si es la primera vez que configuras este proyecto en una Mac nueva, necesitas instalar las herramientas base.

### A. Instalar Go (Golang)
```bash
brew install go

```

### B. Instalar MongoDB

MongoDB no viene instalado por defecto en Mac. Usamos Homebrew:

```bash
brew tap mongodb/brew
brew install mongodb-community@7.0

```

### C. Preparar el Proyecto

Dentro de la carpeta del proyecto, descarga las librerías necesarias (Gin, Mongo Driver, etc.) que figuran en el `go.mod`:

```bash
go mod tidy

```

---

## 🚀 2. Ejecución Diaria (Development)

Cada vez que quieras trabajar en el proyecto, sigue estos dos pasos:

### Paso 1: Encender la Base de Datos

Asegúrate de que el servicio de MongoDB esté corriendo en segundo plano.

```bash
brew services start mongodb-community@7.0

```

> **Nota:** Si ya estaba corriendo, te dirá "already started". Si te da error de conexión, revisa que este servicio esté activo.

### Paso 2: Arrancar el Servidor

Ejecuta el archivo principal de Go. Esto iniciará el servidor web y la conexión a la base de datos.

```bash
go run main.go

```

Si todo es correcto, verás en la terminal:

> `¡Conectado a MongoDB exitosamente!`
> `Listening and serving HTTP on localhost:8080`

### Paso 3: Acceder a la App

Abre tu navegador web y visita:
👉 **http://localhost:8080**

---

## 🆘 Solución de Errores Comunes

* **Error: `connection refused` / `server selection error**`
* **Causa:** MongoDB está apagado.
* **Solución:** Ejecuta el comando del Paso 1 (`brew services start...`).


* **Error: `address already in use**`
* **Causa:** Ya tienes otra terminal corriendo el servidor.
* **Solución:** Busca la terminal abierta y ciérrala, o usa `killall main` si se quedó pegado.



```

```
