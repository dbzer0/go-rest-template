go-rest-template
===========

Проект представляет собой шаблон для создания приложений с набором различного функционала.
Go-rest-template реализует повседневные рутины (логирование, graceful shutdown, билд, и т.д.) и дает контроль над работой подпрограмм, оркестрируемых снаружи.

Просто хочется сэкономить немножко времени в создании новых приложений :)

В linux имя проекта можно задать так:

    sed -i s/PROJECTNAME/NEW_NAME/g Makefile Dockerfile app/Makefile app/main.go docker-compose.yml

Для MacOS:

    sed -i'' -e "s/PROJECTNAME/NEW_NAME/g" Makefile Dockerfile app/Makefile app/main.go docker-compose.yml
    

# Репозиторий

Далее в Dockerfile надо подставить путь к текущему репозиторию, изменив строку: `/go/src/github.com/dbzer0/go-rest-template`.
