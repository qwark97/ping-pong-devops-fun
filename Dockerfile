FROM python:3.8.5-buster

COPY . /app
WORKDIR /app

RUN pip install -r requirements.txt
EXPOSE 5050

CMD ["python", "-m", "pinged", "5050"]