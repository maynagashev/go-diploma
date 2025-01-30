FROM ubuntu:22.04

# Установка curl для healthcheck
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Копируем исполняемый файл accrual
COPY cmd/accrual/accrual_linux_amd64 /usr/local/bin/accrual

# Делаем файл исполняемым
RUN chmod +x /usr/local/bin/accrual

# Открываем порт
EXPOSE 8081

# Запускаем сервис
CMD ["/usr/local/bin/accrual", "-a", ":8081"] 