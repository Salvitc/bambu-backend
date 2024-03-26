# bambu-backend
Proyecto Backend de la tienda online "Bambu Shop". Perteneciente al Trabajo de Fin de Grado de Salvador Toledo.

Proyecto compuesto por un servidor API Rest desarrollado en Golang.

## Ejecución local

### Prerrequisitos
    - Base de datos Mongo funcionando
    - Go versión +1.21.1

Para usar el proyecto en local se debe crear un fichero .env en el directorio raíz del proyecto donde alojar las variables de entorno necesarias para la correcta ejecución del programa.

```
MONGODB_HOST=dirección de la base de datos mongo
MONGODB_PORT=Puerto a usar
MONGODB_DB=Nombre de la base de datos
SECRET_KEY=Clave para el cifrado de los tokens JWT
```

Para arrancar el proyecto, usar el siguiente comando:

`go run cms/backbu/main.go`

## Arquitectura
El siguiente diagrama de clases muestra la interacción de todos los elementos implicados en el proyecto.

![diagrama de clases](https://github.com/Salvitc/bambu-backend/blob/main/doc/backbu.png?raw=true)

## Despliegue en producción
    [TODO]