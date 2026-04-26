# 📘 Guía Avanzada: Funciones, Datos y Soluciones

Esta guía profundiza en las capacidades técnicas del **Google Maps Scraper**, los datos que extrae y cómo resolver los obstáculos más comunes.

---

## 🛠️ ¿Para qué sirve Go (Golang)?
El usuario preguntó: *¿No tengo Go para qué sirve?*
- **Respuesta:** **Go** es el lenguaje de programación en el que se escribió esta herramienta. Es como el "motor" bajo el capó.
- **Si NO tienes Go:** No pasa nada. Para eso usamos **Docker**. Docker empaqueta el motor (Go) y todas las piezas necesarias para que tú solo tengas que darle al "Play".

---

## 🌟 Funciones Principales

### 1. Interfaz Web (Web UI)
- **Qué es:** Una página web amigable que corre en tu propio ordenador.
- **Cómo usarla:** Una vez lanzada con Docker, vas a `localhost:8080`.
- **Ventaja:** Puedes crear "Jobs" (tareas), ver el progreso en tiempo real y descargar los resultados como CSV con un botón.

### 2. Extracción de Emails 📧
- **Cómo se activa:** Añadiendo el parámetro `-email` en la línea de comandos o marcando la opción en la web.
- **Qué hace:** Después de encontrar un negocio, el scraper entra en su página web oficial y busca correos electrónicos de contacto.

### 3. Modo Rápido (Fast Mode) ⚡
- **Cómo se activa:** Parámetro `-fast-mode`.
- **Qué hace:** Recoge solo los primeros 20 resultados por búsqueda de forma extremadamente rápida. Útil para "muestreos" veloces.

### 4. Soporte de Proxies 🛰️
- **Qué hace:** Permite cambiar tu dirección IP constantemente para que Google no te detecte como un robot y te bloquee.

---

## 📊 Datos que Recoge (Los 33+ Puntos)

La herramienta extrae prácticamente todo lo que ves en Google Maps. Algunos de los más importantes son:

| Dato | Descripción |
| :--- | :--- |
| **Título** | Nombre comercial del negocio. |
| **Categoría** | Tipo de negocio (Ej: Restaurante Italiano). |
| **Dirección** | Dirección exacta y código postal. |
| **Web** | URL del sitio oficial. |
| **Teléfono** | Número de contacto. |
| **Email** | Correos extraídos (solo si activas `-email`). |
| **Rating** | Puntuación (Ej: 4.5 estrellas). |
| **Reviews** | Número total de reseñas y texto de las mismas. |
| **Horario** | Horas de apertura y cierre. |
| **Coordenadas** | Latitud y Longitud GPS. |
| **Imágenes** | Enlaces a las fotos del local. |
| **Estado** | Indica si está "Abierto", "Cerrado temporalmente", etc. |

---

## ⚠️ Limitaciones y Restricciones

1.  **Bloqueos de Google:** Si haces miles de búsquedas seguidas sin proxies, Google te mostrará un CAPTCHA o bloqueará tu IP temporalmente.
2.  **Consumo de Recursos:** Al funcionar abriendo un navegador "invisible" (Chrome/Chromium), consume mucha **CPU y Memoria RAM**. No te asustes si los ventiladores de tu PC suenan más fuerte.
3.  **Velocidad de Emails:** Buscar correos triplica el tiempo de espera, porque requiere visitar webs externas.

---

## ❓ Problemas Comunes y Soluciones

| Problema | Causa Probable | Solución |
| :--- | :--- | :--- |
| **"Error connecting to Docker"** | Docker Desktop no está abierto. | Abre Docker Desktop y espera a que el icono esté en verde. |
| **No aparecen resultados** | La búsqueda es demasiado específica o errónea. | Prueba algo más simple como "Hoteles en Madrid". |
| **El archivo de salida está vacío** | El programa se cerró antes de guardar. | Asegúrate de usar `-exit-on-inactivity 3m`. |
| **Google me ha bloqueado** | Demasiadas peticiones desde tu IP. | Espera unas horas o usa el parámetro `-proxies`. |
| **Va muy lento** | Muchas búsquedas simultáneas. | Reduce el número de hilos con el parámetro `-c` (ejemplo: `-c 2`). |

---
*Guía detallada para la Carpeta de Comunicación Interina.*
