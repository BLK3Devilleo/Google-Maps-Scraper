# 🎨 Guía: Personalización y Búsquedas Geo-Específicas

En esta guía te explico cómo hacer que el scraper se vea como tú quieras y cómo encontrar exactamente esos negocios que no tienen web en tu zona.

---

## 🖌️ ¿Cómo modificar la interfaz web?

Sí, puedes modificarla totalmente. Los archivos que definen qué ves en tu navegador están aquí:
- **HTML (Estructura):** `web/static/templates/index.html`
- **CSS (Diseño/Colores):** `web/static/css/` (aquí encontrarás los estilos).

**Pasos para modificarla:**
1. Abre `web/static/templates/index.html` y cambia el texto, los colores o el logo.
2. Guarda los cambios.
3. **Importante:** Como estás usando Dokploy, después de guardar los cambios en tu PC local, debes **subirlos a GitHub** y darle a **Deploy** de nuevo para que Dokploy reconstruya la imagen con tu nueva interfaz.

---

## 📍 Búsqueda de Spas en un radio de 30km

Para buscar en un área específica de 30km a la redonda, necesitas usar parámetros geográficos.

### Opción A: Desde la Web UI
En la interfaz, busca los campos de "Latitude", "Longitude" y "Radius".
- **Radius:** Pon `30000` (esto son 30.000 metros = 30km).
- **Lat/Lon:** Busca las coordenadas de tu ciudad en Google (ej: Madrid es `40.4168, -3.7038`).

### Opción B: Desde la terminal (PowerShell o Docker)
```powershell
docker run gosom/google-maps-scraper -input consultas.txt -geo "TU_LATITUD,TU_LONGITUD" -radius 30000 -results resultados.csv
```

---

## 🚫 Cómo encontrar Spas SIN sitio web

Google Maps no permite filtrar directamente por "Negocios sin web" en la búsqueda inicial. El truco es el siguiente:

1.  **Haz el escaneo completo** de todos los Spas en tu radio de 30km.
2.  **Descarga el CSV** de resultados.
3.  **Filtra en Excel/Google Sheets**:
    - Abre el archivo CSV.
    - Busca la columna llamada `website`.
    - Filtra los resultados para mostrar solo las filas donde esa columna esté **vacía**.

**¿Por qué hacerlo así?**
Porque así te aseguras de no dejarte ninguno fuera. El scraper recogerá todos, y tú con un solo clic en Excel tendrás la lista de los que no tienen web (tus clientes potenciales).

---

## 💡 Pro-Tip para mejorar la puntería

Si pones un radio muy grande (30km), Google a veces se satura. Es mejor:
- Usar un **Zoom** de nivel `13` o `14` (cubre bien áreas urbanas).
- Si ves que faltan resultados, divide la ciudad en 2 o 3 puntos geográficos distintos con radios más pequeños (ej: 10km cada uno).

---
*Guía de personalización y búsqueda estratégica generada para la Carpeta de Comunicación Interina.*
