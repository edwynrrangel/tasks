<div align="center">
  <h1>API para el registro de Tareas con autenticación JWT haciendo uso Go Fiber</h1>
</div>

# Introducción
Este repositorio contiene una API para el registro de tareas con autenticación JWT en Go Fiber. El objetivo de este repositorio es estructurar un proyecto base con Clean Architecture. Se establece la estructura de carpetas necesaria para el dominio, casos de usos, repositorios, controladores y elementos transversales.

# Tabla de contenido

- [Estructura de carpetas](#estructura-de-carpetas)


# Estructura de carpetas
Estructura base para el manejo de carpetas
```bash
.
├── /auth                       # Paquete de autenticación 
├── /comments                   # Paquete de comentararios para las tareas
├── /config                     # Paquete transversal con la configuración global
├── /errors                     # Paquete transversal para manejar los errores de la API
├── /jwt                        # Paquete transversal para manejar los token jwt
├── /logger                     # Paquete transversal para imprimir los logs generados en la aplicación
├── /middlewares                # Contiene los middleware personalizados para Fiber
|   ├── /role				    # Paquete para validar los roles de los usuarios
|	├── /session                # Paquete para validar la sesión de los usuarios
├── /password                   # Paquete transversal manejo de las contraseña
├── /roles                      # Paquete para el manejo de los roles para los usuarios
├── /states                     # Paquete para el manejo de los estados de las tareas
├── /tasks                      # Paquete para el manejo de las tareas
├── /token                      # Paquete para procesar los token jwt
├── /users                      # Paquete para el manejo de los usuarios en la API
├── /validator                  # Paquete transversal para validar las entradas de datos en los endpoint
├── /.env                       # Archivo con las variables de entorno a usar
├── /.env.template              # Plantilla con las variables de entorno que se pueden usar
├── go.mod						# Se define el path del módulo, además de las dependencias para el proceso de compilación
├── go.sum						# Este archivo lista el checksum de las dependencia directa e indirecta, además de incluir la versión
└── main.go					    # Punto de entrada para la función main
```

# Arquitectura limpia
La base de este repositorio es la arquitectura de puertos y adaptadores, también conocida como la arquitectura hexagonal.

![Clean Architecture](clean_architecture.png)
# Entidades del dominio
Todos los modelos de dominio se colocarán en cada paquete relacionado con las entedidades de la API. Contiene la definición *go struct* de cada entidad que forma parte del dominio y que puede ser utilizada en toda la aplicación.

# Puertos
Los puertos son la interfaz que pertenece al núcleo y definen cómo se debe abordar la comunicación entre los actores y el núcleo.

## Puertos primarios
Se definen los casos de uso que el núcleo implementara y se expone para ser consumido por los actores externos. Estos se encontraran adentro del archivo interface y son las interfaces **Service**.

## Puertos secundarios
Definen las acciones que la capa de datos debe implementar. Serían las interfaces **Repository**

# Adaptadores
Esta capa que sirve para transformar la comunicación entre actores externos y la lógica de la aplicación de forma que ambas quedan independientes.
## Adaptadores primarios
Es la implementación de los casos de uso definidos en el puerto primario. Estas implementaciones se encuentra adentro de los archivos **service.go**

## Adaptadores secundarios
Es la implementación de los puertos secundarios relacionados con la capa de datos. Estas implementaciones se encuentra adentro de los archivos **repository.go**

# Requests y Responses
Son un tipo de estructura que sirven únicamente para transportar datos, estas estructuras contienen las propiedades de la entidad. Las estructuras pueden tener su origen en una o más entidades.

Las estructuras del request se pueden usar para la validación de los datos de entrada.

# Infraestructura para la API
Esta sección permite exponer los puntos de entrada de la aplicación a través del framework Fiber. 
## Controladores
Son los elementos que contiene la lógica de los endpoint, en estas funciones se recibe la petición de los clientes y se llaman los casos de uso para el procesamiento de los datos, luego del que *core* procesa los datos y se encarga de retornar una respuesta a los clientes.

## Middleware
Este módulo contiene los middleware personalizados para Fiber.

# Elementos transversales
Esta compuesto por aquellos módulos que son transversales a la apliación.
## logger
Permite imprimir log estructurados y nivelados haciendo uso de la librería ZAP.
## errors
Son un conjunto de funciones que permite manejar los errores para la API.
## validator
Es un wrapper de la libraría github.com/go-playground/validator/v10 para validar la entrada de datos a los endpoint.

# Correr proyecto

## Clonar repo https
```
git clone https://github.com/edwynrrangel/tasks.git
```

## O clonar con ssh
```
git clone git@github.com:edwynrrangel/tasks.git
```
## Descargar dependencias
```
go install github.com/golang/mock/mockgen@v1.6.0
go generate ./...
go mod tidy
```

## Iniciar la aplicación
```
go run main.go
```

# Correr pruebas unitarias
Lo primero que se debe hacer es generar los mocks
```bash
go generate ./...
```

Luego de generar los mock se pueden correr las pruebas de la siguiente manera
```bash
go test -coverprofile=cover.out ./... 
go tool cover -func=cover.out
```
Para visualizar la cobertura en el navegador se debe correr el siguiente comando 
```bash
go tool cover -html=cover.out
```

# Correr API con docker-compose
Desde el directorio root, correr el siguiente comando
```
docker-compose up -d
```

# Casos de uso:
  - Login de usuario (debe considerar 3 perfiles). El login debe entregar un “token” o algún mecanismo que permita identificar el perfil el cual se usará para identificar sobre qué servicios puede consumir.
  - Los usuarios con perfil Administrador pueden realizar lo siguiente:
    - CRUD de usuarios
      - Los usuarios parten con una contraseña temporal la cual deben cambiar en el primer login
      - Los usuarios puedes tener perfil “Ejecutor” o “Auditor”, no puede crear otros usuarios tipo “Administrador”

    - CRUD de “tareas”
      - Una tarea tiene al menos los siguientes datos: título, descripción, fecha de vencimiento.
      - Cuando crea una tarea debe poder asignársela a un usuario con perfil “Ejecutor”.
      - No puede eliminar o actualizar una tarea en estado distinto a “Asignado”

  - Los usuarios con perfil Ejecutor pueden realizar lo siguiente:
    - Listar sus tareas asignadas y ver el detalle
    - Actualizar el estado de una tarea. Si la tarea ya está vencida no debe permitir esta acción.
    - Agregar un comentario sobre una tarea vencida.
    
  - Los usuarios con perfil Auditor pueden realizar lo siguiente:
    - Visualizar el listado de tareas asignadas a cualquier usuario y ver su estado.

  - Actualización de contraseña (para cualquier perfil)
  - Logout de usuario