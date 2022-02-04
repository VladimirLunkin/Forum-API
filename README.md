# tech-db-forum

Проект "Форумы" второго семестра Технопарка по курсу "Базы данных".

Суть задания заключается в реализации API к базе данных проекта «Форумы» по документации к этому API.

Таким образом, на входе:

* документация к API в файле ./swagger.yaml;

На выходе:

* репозиторий, содержащий все необходимое для разворачивания сервиса в Docker-контейнере.

## Документация к API
Документация к API предоставлена в виде спецификации OpenAPI: [swagger.yml](https://github.com/VladimirLunkin/tech-db-forum/swagger.yml) 

## Требования к проекту
Проект должен включать в себя все необходимое для разворачивания сервиса в Docker-контейнере.

При этом:

* файл для сборки Docker-контейнера должен называться Dockerfile и располагаться в корне репозитория;
* реализуемое API должно быть доступно на 5000-ом порту по протоколу http;
* допускается использовать любой язык программирования;
* крайне не рекомендуется использовать ORM.

Контейнер будет собираться из запускаться командами вида:
```
docker build -t <username> https://github.com/mailcourses/technopark-dbms-forum-server.git
docker run -p 5000:5000 --name <username> -t <username>
```
