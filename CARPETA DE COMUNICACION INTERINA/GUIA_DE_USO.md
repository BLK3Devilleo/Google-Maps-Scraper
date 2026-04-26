# 📍 Guía Maestra: Google Maps Scraper (Español)

Bienvenido a la guía de uso del **Google Maps Scraper**. Esta herramienta es extremadamente potente para extraer datos de negocios directamente desde Google Maps. A continuación, te explico paso a paso cómo ponerla a funcionar.

---

## 🛠 Requisitos Previos

Antes de empezar, asegúrate de tener instalado:
1.  **Go (Golang):** Si planeas ejecutarlo desde el código fuente o compilarlo. [Descargar Go](https://go.dev/dl/).
2.  **Docker (Recomendado):** La forma más fácil y limpia de ejecutarlo sin configurar dependencias locales. [Descargar Docker Desktop](https://www.docker.com/products/docker-desktop/).

---

## 🚀 Formas de Hacerlo Funcionar

Tienes tres formas principales de usar esta herramienta:

### 1. Interfaz Web (La más fácil) 🌐
Ideal si prefieres usar el ratón y ver los resultados en tu navegador.

**Usando Docker:**
Ejecuta este comando en tu terminal (PowerShell o CMD):
```powershell
mkdir data
docker run -v ${PWD}/data:/gmapsdata -p 8080:8080 gosom/google-maps-scraper -data-folder /gmapsdata
```
Luego, abre tu navegador en: `http://localhost:8080`

### 2. Línea de Comandos (CLI) 💻
Ideal para automatización y rapidez.

**Ejemplo básico:**
Para buscar "Restaurantes en Madrid" y guardar los resultados en un CSV:
1. Crea un archivo llamado `consultas.txt` y escribe dentro: `Restaurantes en Madrid`
2. Ejecuta:
```powershell
docker run -v ${PWD}/consultas.txt:/input.txt -v ${PWD}/resultados.csv:/results.csv gosom/google-maps-scraper -input /input.txt -results /results.csv -exit-on-inactivity 3m
```

### 3. Ejecución desde el Código (Para desarrolladores) 🏗
Si tienes Go instalado:
```bash
go build
./google-maps-scraper -input consultas.txt -results resultados.csv
```

---

## ⚙️ Configuraciones Importantes

Cuando lances el scraper, puedes añadir estas opciones para mejorar los resultados:

| Opción | Descripción | Ejemplo |
| :--- | :--- | :--- |
| `-email` | Extrae correos electrónicos de las webs de los negocios. | `-email` |
| `-depth` | Cuánto scroll hace en los resultados (por defecto 10). | `-depth 20` |
| `-lang` | Idioma de los resultados (por defecto 'en'). | `-lang es` |
| `-json` | Guarda los resultados en formato JSON en vez de CSV. | `-json` |
| `-proxies` | Usa proxies para evitar que Google te bloquee. | `-proxies "http://user:pass@host:port"` |

---

## 📂 ¿Dónde están mis datos?

- **En la Web UI:** Puedes descargarlos directamente desde el apartado "Jobs" en la interfaz del navegador.
- **En la CLI:** Estarán en el archivo que especifiques con el parámetro `-results` (ejemplo: `resultados.csv`).

---

## 💡 Consejos de Oro

1.  **Inactividad:** El parámetro `-exit-on-inactivity 3m` es muy útil para que el programa se cierre solo cuando termine de trabajar.
2.  **Bloqueos:** Si vas a extraer miles de datos, **usa proxies**. Google detectará muchas peticiones desde tu IP y podría bloquearte temporalmente.
3.  **Emails:** Activar `-email` hace que el proceso sea más lento, ya que el scraper tiene que entrar en cada página web del negocio para buscar el contacto.

---

> [!TIP]
> Si estás en Windows y usas Docker, asegúrate de que Docker Desktop esté abierto y funcionando antes de lanzar los comandos.

---
*Guía generada para la Carpeta de Comunicación Interina.*
