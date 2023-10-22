<div align="center">
  <h1>API para el registro de Tareas con autenticación JWT haciendo uso Go Fiber</h1>
</div>

# Introducción
Este repositorio contiene una API para el registro de tareas con autenticación JWT en Go Fiber. El objetivo de este repositorio es estructurar un proyecto base con Clean Architecture. Se establece la estructura de carpetas necesaria para el dominio, casos de usos, repositorios, controladores y elementos transversales.

# Tabla de contenido

- [Estructura de Carpetas](#estructura-de-carpetas)
- [Arquitectura Limpia](#arquitectura-limpia)
- [Entidades del Dominio](#entidades-del-dominio)
- [Puertos](#puertos)
  - [Puertos Primarios](#puertos-primarios)
  - [Puertos Secundarios](#puertos-secundarios)
- [Adaptadores](#adaptadores)
  - [Adaptadores Primarios](#adaptadores-primarios)
  - [Adaptadores Secundarios](#adaptadores-secundarios)
- [Requests y Responses](#requests-y-responses)
- [Infraestructura para la API](#infraestructura-para-la-api)
  - [Controladores](#controladores)
  - [Middleware](#middleware)
- [Elementos Transversales](#elementos-transversales)
  - [Logger](#logger)
  - [Errors](#errors)
  - [Validator](#validator)
- [Correr Proyecto](#correr-proyecto)
- [Correr Pruebas Unitarias](#correr-pruebas-unitarias)
- [Correr API con docker-compose](#correr-api-con-docker-compose)
- [Casos de Uso](#casos-de-uso)
  - [Caso de Uso 1: Inicio de Sesión de Usuario](#caso-de-uso-1-inicio-de-sesión-de-usuario)
  - [Caso de Uso 2: CRUD de Usuarios (Crear, Leer, Actualizar, Eliminar)](#caso-de-uso-2-crud-de-usuarios-crear-leer-actualizar-eliminar)
  - [Caso de Uso 3: CRUD de Tareas](#caso-de-uso-3-crud-de-tareas)
  - [Caso de Uso 4: Gestión de Tareas por Ejecutor](#caso-de-uso-4-gestión-de-tareas-por-ejecutor)
  - [Caso de Uso 5: Revisión de Tareas por Auditor](#caso-de-uso-5-revisión-de-tareas-por-auditor)
  - [Caso de Uso 6: Actualización de Contraseña](#caso-de-uso-6-actualización-de-contraseña)
  - [Caso de Uso 7: Cierre de Sesión de Usuario](#caso-de-uso-7-cierre-de-sesión-de-usuario)
- [Diagramas](#diagramas)
  - [Diagramas de Flujos](#diagramas-de-flujos)
  - [Diagramas de Componentes](#diagramas-de-componentes)
  - [Diagramas de Secuencias](#diagramas-de-secuencias)


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

# Casos de Uso

### Caso de Uso 1: Inicio de Sesión de Usuario
- **Actores**: Administrador, Ejecutor, Auditor.
- **Precondiciones**: El usuario debe tener una cuenta registrada en el sistema.
- **Flujo Principal**:
  1. El usuario ingresa su nombre de usuario y contraseña.
  2. El sistema valida las credenciales y genera un token que identifica el perfil del usuario.
  3. El sistema entrega el token al usuario y le permite acceder a las funcionalidades correspondientes a su perfil.
- **Postcondiciones**: El usuario tiene acceso al sistema con un perfil específico.

### Caso de Uso 2: CRUD de Usuarios (Crear, Leer, Actualizar, Eliminar)
- **Actores**: Administrador.
- **Precondiciones**: El administrador debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El administrador accede a la sección de gestión de usuarios.
  2. El administrador puede crear, ver, actualizar o eliminar usuarios.
  3. En caso de creación, el sistema genera una contraseña temporal para el nuevo usuario.
- **Postcondiciones**: Los datos de los usuarios en el sistema son actualizados según las acciones del administrador.

### Caso de Uso 3: CRUD de Tareas
- **Actores**: Administrador.
- **Precondiciones**: El administrador debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El administrador accede a la sección de gestión de tareas.
  2. El administrador puede crear, ver, actualizar o eliminar tareas.
  3. En caso de creación, el administrador asigna la tarea a un usuario con perfil Ejecutor.
- **Restricciones**: Las tareas solo pueden ser eliminadas o actualizadas si están en estado "Asignado".

### Caso de Uso 4: Gestión de Tareas por Ejecutor
- **Actores**: Ejecutor.
- **Precondiciones**: El ejecutor debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El ejecutor accede a la lista de tareas asignadas.
  2. El ejecutor puede ver el detalle de una tarea o actualizar el estado de una tarea.
  3. Si la tarea está vencida, el ejecutor puede agregar un comentario sobre la tarea.
- **Restricciones**: No se permite actualizar el estado de una tarea vencida.

### Caso de Uso 5: Revisión de Tareas por Auditor
- **Actores**: Auditor.
- **Precondiciones**: El auditor debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El auditor accede a la lista de tareas asignadas a cualquier usuario.
  2. El auditor puede ver el estado de las tareas.

### Caso de Uso 6: Actualización de Contraseña
- **Actores**: Administrador, Ejecutor, Auditor.
- **Precondiciones**: El usuario debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El usuario accede a la sección de actualización de contraseña.
  2. El usuario ingresa su contraseña actual y la nueva contraseña.
  3. El sistema valida y actualiza la contraseña del usuario.

### Caso de Uso 7: Cierre de Sesión de Usuario
- **Actores**: Administrador, Ejecutor, Auditor.
- **Precondiciones**: El usuario debe estar autenticado en el sistema.
- **Flujo Principal**:
  1. El usuario selecciona la opción de cerrar sesión.
  2. El sistema invalida el token del usuario y cierra la sesión.
- **Postcondiciones**: El usuario no tiene acceso al sistema hasta que vuelva a iniciar sesión.

# Diagramas

## Diagramas de flujos
### Caso de Uso 1: Inicio de Sesión de Usuario
![Caso de Uso 1](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%201.svg)

### Caso de Uso 2: CRUD de Usuarios (Crear, Leer, Actualizar, Eliminar)
![Caso de Uso 2](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%201.svg)

### Caso de Uso 3: CRUD de Tareas
![Caso de Uso 3](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%202.svg)

### Caso de Uso 4: Gestión de Tareas por Ejecutor
![Caso de Uso 4](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%203.svg)

### Caso de Uso 5: Revisión de Tareas por Auditor
![Caso de Uso 5](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%205.svg)

### Caso de Uso 6: Actualización de Contraseña
![Caso de Uso 6](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%206.svg)

### Caso de Uso 7: Cierre de Sesión de Usuario
![Caso de Uso 7](diagrams/diagrama%20de%20flujo%20caso%20de%20uso%207.svg)

## Diagramas de componentes
![componentes](diagrams/diagrama%20de%20componente.svg)

## Diagrama de secuencias
### Caso de Uso 1: Inicio de Sesión de Usuario
![Caso de Uso 1](diagrams/diagrama%20de%20secuencias%201.svg)

### Caso de Uso 2: CRUD de Usuarios (Crear, Leer, Actualizar, Eliminar)
![Caso de Uso 2](diagrams/diagrama%20de%20secuencias%202.svg)

### Caso de Uso 3: CRUD de Tareas
![Caso de Uso 3](diagrams/diagrama%20de%20secuencias%203.svg)

### Caso de Uso 4: Gestión de Tareas por Ejecutor
![Caso de Uso 4](diagrams/diagrama%20de%20secuencias%204.svg)

### Caso de Uso 5: Revisión de Tareas por Auditor
![Caso de Uso 5](diagrams/diagrama%20de%20secuencias%205.svg)

### Caso de Uso 6: Actualización de Contraseña
![Caso de Uso 6](diagrams/diagrama%20de%20secuencias%206.svg)

### Caso de Uso 7: Cierre de Sesión de Usuario
![Caso de Uso 7](diagrams/diagrama%20de%20secuencias%207.svg)

# Base de datos
[Modelo Relacional](https://dbdiagram.io/d/app-tasks-6531d265ffbf5169f01160f3)
