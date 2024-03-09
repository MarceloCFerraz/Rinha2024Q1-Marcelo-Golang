# Use the official .NET SDK image as the base image
FROM python:3.11.7-slim-bullseye

# Set the working directory inside the container
WORKDIR /api

# Copy the contents of src directory
COPY src/ ./

RUN pip install -r requirements.txt

# Set the entry point for the container
ENTRYPOINT ["python", "init.py"]
