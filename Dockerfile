# Etapa de compilación
FROM golang:1.24-alpine AS builder

# Instalar las herramientas de compilación C necesarias para CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copiar el go.mod y go.sum para descargar dependencias
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY src/ .

# Compilar la aplicación con CGO habilitado
RUN go build -o main .

# Etapa final
FROM alpine:latest

# Instalar los certificados y la librería de sqlite para el binario compilado
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copiar el binario compilado de la etapa de compilación
COPY --from=builder /app/main .

# Copiar la carpeta de archivos estáticos (para el QR)
COPY --from=builder /app/statics ./statics

# Exponer el puerto correcto para esta aplicación
EXPOSE 3000

# Comando para correr la aplicación
CMD ["./main", "rest"] 