
## Миграции

Создаются миграции в дикертории */migrations*


### Накатить миграцию

```
migrate -database "postgresql://www:pass@localhost:15434/todo?sslmode=disable" -path migrations up
```

### Откатиться к миграции N

```
migrate -database "postgresql://www:pass@localhost:15434/todo?sslmode=disable" -path migrations goto N
```

