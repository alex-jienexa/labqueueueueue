# labQueueueueueue

Очередь для студентов от студентов к студентам. 

Записывайтесь в очереди, сражайтесь за места, предавайте и продавайте. 
Делайте всё ради того, чтобы первым сдать лабу у препода.

## Запуск через Docker Compose

1. Загрузите весь код и выполните `docker-compose build`. 
Очень рекомендуется при первой сборке включить VPN или какой-нибудь прокси, так как один из пакетов для Go вызывает ошибку 403 при установке (по крайней мере так было у меня).
2. После загрузки выполните `docker-compose up`. Должны появиться следующие элементы
```
 ✔ Network go_queueueue_default
 ✔ Volume "go_queueueue_postgres_data"
 ✔ Container go_queueueue-postgres-1
 ✔ Container go_queueueue-backend-1
 ✔ Container go_queueueue-frontend-1
 ```
 3. Перейдите по адресу http://localhost:3000 (для фронтенда) или http://localhost:8080/ping (для бекенда).
 4. ...
 5. Profit!

 Все методы для бекенда смотрите в [README бекенда](backend/README.md).