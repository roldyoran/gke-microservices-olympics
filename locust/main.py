from locust import HttpUser, TaskSet, task, between
import random
import json

class MyTasks(TaskSet):
    
    @task(1)
    def engineering(self):
        # List of random names
        names = ["Alice Doe", "Bob Smith", "Charlie Brown", "David Johnson", "Eve Williams", "Frank Jones", "Grace Lee", "Henry Thomas"]
    
        # Student data
        student_data = {
            "name": random.choice(names),  # Random name
            "age": random.randint(18, 28),  # Random age
            "faculty": "engineering",  # Faculty
            "discipline": random.choice([1, 2, 3])  # 1 = Swimming, 2 = Athletics, 3 = Boxing
        }
        
        # Send JSON as POST to the /engineering route
        headers = {'Content-Type': 'application/json'}
        self.client.post("/engineering", json=student_data, headers=headers)

    @task(1)
    def agronomy(self):
        # List of random names
        names = ["Juan Perez", "Maria Sanchez", "Pedro Rodriguez", "Ana Torres", "Luis Lopez", "Carlos Rodriguez", "Juan Garcia", "Maria Rodriguez"]
    
        # Student data
        student_data = {
            "name": random.choice(names),  # Random name
            "age": random.randint(20, 36),  # Random age
            "faculty": "agronomy",  # Faculty
            "discipline": random.choice([1, 2, 3])  # 1 = Swimming, 2 = Athletics, 3 = Boxing
        }
        
        # Send JSON as POST to the /agronomy route
        headers = {'Content-Type': 'application/json'}
        self.client.post("/agronomy", json=student_data, headers=headers)

class WebsiteUser(HttpUser):
    tasks = [MyTasks]
    wait_time = between(1, 5)  # Wait time between tasks (1 to 5 seconds)
