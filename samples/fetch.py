import requests

url = "http://localhost:8080/generate"
payload = {
    "language": "Python"
}

response = requests.post(url, json=payload)
print(response.text)

with open("python-course.txt", "w") as file:
    file.write(response.text)

print("Course saved to python-course.txt")
