# SUBD
Семестровый проект в рамках курса Технопарка по СУБД.

### Build:
```shell
sudo docker build -t subd-proj . ;
sudo docker run -p 5000:5000 --name subd-proj -t  subd-proj
```

### Test:
[Скомпилированная тестирующая программа](./linux64amd/tech-db-forum)
```shell
# функциональное тестирование:
./tech-db-forum func -uk http://localhost:5000/api -r report.html;

# заполнение:
./tech-db-forum fill --url=http://localhost:5000/api --timeout=900;

# нагрузочное тестирование:
./tech-db-forum perf --url=http://localhost:5000/api --duration=600 --step=60;
```

[Результаты функционального тестирования](./linux64amd/report.html)

Результаты нагрузочного тестирования: ~1200 rps

Тестовая документация содержится в файле [swagger.yml](./swagger.yml), для просмотра функционала апи можно
воспользоваться любым swagger-reader. Например: 
[Swagger Editor](https://editor.swagger.io/)

Подробное описание тестов можно найти в [репозитории](https://github.com/mailcourses/technopark-dbms-forum) курса технопарка по СУБД.