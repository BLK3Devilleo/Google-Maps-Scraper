# 🚀 Despliegue en Dokploy: Solución paso a paso

Tu despliegue falló porque Dokploy intentó usar un archivo llamado `docker-compose.yml` que no existía. He creado ese archivo por ti en la raíz del proyecto.

---

## 🛠️ ¿Qué he cambiado?
He añadido un archivo **`docker-compose.yml`** en la carpeta principal. Este archivo le dice a Dokploy:
1.  **Construir** la aplicación usando el `Dockerfile` que ya tienes.
2.  **Abrir el puerto 8092**, que es donde vivirá la interfaz web según tu preferencia.
3.  **Crear un volumen persistente** llamado `gmapsdata` para que no pierdas tus datos si el servidor se reinicia.
4.  **Activar el modo Web** automáticamente.

---

## 📋 Pasos para completar el despliegue

1.  **Sube los cambios a GitHub**: Asegúrate de que el nuevo archivo `docker-compose.yml` esté en tu repositorio de GitHub.
2.  **Reintenta el Deploy en Dokploy**:
    - Ve a tu proyecto en Dokploy.
    - Asegúrate de que el **Compose Type** sea `docker-compose`.
    - Haz clic en **Deploy**.
3.  **Configura el Dominio**: Una vez que el despliegue sea exitoso (dirá ✅), configura un dominio o usa la IP de tu servidor con el puerto `:8092`.

---

## ⚠️ ¡IMPORTANTE: Recursos del Servidor!

El scraper utiliza un navegador interno para extraer los datos. Esto consume **mucha memoria RAM**.
- **Recomendación:** Tu servidor debe tener al menos **2GB de RAM** (preferiblemente 4GB) para que no se cuelgue durante el proceso de extracción.
- Si ves que el despliegue falla por "Out of Memory", tendrás que aumentar la potencia de tu VPS.

---

## 📂 Gestión de Datos en Dokploy

El volumen que he configurado guardará todo en una carpeta interna del contenedor. Si necesitas acceder a los archivos CSV directamente desde el servidor:
- En Dokploy, dentro de la configuración del servicio, verás el apartado de **Mounts/Volumes**.
- Los datos se guardarán en `/gmapsdata` dentro del contenedor.

---
*Guía de despliegue generada para la Carpeta de Comunicación Interina.*
