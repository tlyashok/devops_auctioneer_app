Установка:
1. Ставим golang, make, запускаем docker engine
2. Утилита для миграций: 
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest  
2. Через терминал:  
   make env  
   make up  
   make migrate_up  
   make run  
