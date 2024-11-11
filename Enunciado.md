# Enunciado

Este proyecto de ejemplo de microservicios se centra en monitorizar las Olimpiadas de la Universidad de San Carlos de Guatemala mediante una arquitectura basada en microservicios y contenedores, desplegada en Google Cloud Platform (GCP) con Kubernetes.

### Objetivos:
- **Despliegue en GCP**: Utilización de Google Kubernetes Engine (GKE) para gestionar el tráfico en tiempo real generado por la competencia.
- **Microservicios en Golang y Rust**: Servicios para facultades y disciplinas implementados en Golang y Rust, con comunicación mediante gRPC.
- **Kafka y Redis**: Uso de Apache Kafka para la transmisión de eventos de estudiantes ganadores y perdedores, y Redis para almacenamiento temporal de resultados.
- **Monitoreo en Grafana y Prometheus**: Visualización de datos en Grafana y monitoreo de recursos con Prometheus.

### Componentes Clave:
1. **Generación de tráfico con Locust**: Simulación de tráfico que dirige datos al Ingress de Kubernetes.
2. **Microservicios de facultades**: Servicios que reciben solicitudes de estudiantes y los dirigen a los servidores de disciplinas correspondientes.
3. **Microservicios de disciplinas**: Estos deciden si un estudiante es ganador o perdedor y envían los resultados a Kafka.
4. **Consumo y visualización**: Datos procesados por consumidores de Kafka se almacenan en Redis, y Grafana los visualiza en tiempo real.

Este ejemplo ilustra la gestión y escalado automático de servicios distribuidos en la nube, adaptándose a altas demandas de tráfico mediante una infraestructura de contenedores y arquitectura de microservicios.

### Requisitos:
- GCP con cuenta de Google Cloud Platform (GCP) válida.
- Conocimientos básicos de Kubernetes y Docker.
- Conocimientos básicos de Kafka y Redis.
- Conocimientos básicos de Golang y Rust.
- Conocimientos básicos de Prometheus y Grafana.

### Arquitectura:
![Arquitectura](./imgs/arquitectura.png)

