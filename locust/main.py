from locust import HttpUser, TaskSet, task, between
import random
import json

class MyTasks(TaskSet):
    
    @task(1)
    def ingenieria(self):
        # Lista de nombres aleatorios
        names = ["Alice Doe", "Bob Smith", "Charlie Brown", "David Johnson", "Eve Williams", "Frank Jones", "Grace Lee", "Henry Thomas"]
    
        
        # Datos del estudiante
        student_data = {
            "name": random.choice(names),  # Nombre al azar
            "age": random.randint(18, 28),  # Edad aleatoria
            "faculty": "Ingenieria",  # Facultad al azar
            "discipline": random.choice([1, 2, 3])  # 1 = Natación, 2 = Atletismo, 3 = Boxeo
        }
        
        # Enviar el JSON como POST a la ruta /ingenieria
        headers = {'Content-Type': 'application/json'}
        self.client.post("/ingenieria", json=student_data, headers=headers)


    @task(1)
    def agronomia(self):
        # Lista de nombres aleatorios
        names = ["Juan Perez", "Maria Sanchez", "Pedro Rodriguez", "Ana Torres", "Luis Lopez", "Carlos Rodriguez", "Juan Garcia", "Maria Rodriguez"]
    
        
        # Datos del estudiante
        student_data = {
            "name": random.choice(names),  # Nombre al azar
            "age": random.randint(20, 36),  # Edad aleatoria
            "faculty": "Agronomia",  # Facultad al azar
            "discipline": random.choice([1, 2, 3])  # 1 = Natación, 2 = Atletismo, 3 = Boxeo
        }
        
        # Enviar el JSON como POST a la ruta /agronomia
        headers = {'Content-Type': 'application/json'}
        self.client.post("/agronomia", json=student_data, headers=headers)

class WebsiteUser(HttpUser):
    tasks = [MyTasks]
    wait_time = between(1, 5)  # Tiempo de espera entre tareas (1 a 5 segundos)
