meteo
=================
home meteostation [description](http://elbroom.ru/post/meteo-1)

[Reader for serial port] (/blynk_serial.py)

[Scketch for Arduino] (/scketch.c)

Setup 
-----

1. Dowanload packages

1. Copy config

        cp config/app_tpl.yml config/app.yml
        cp config/sql_tpl.yml config/sql.yml
        
1. Set envariment variable

        export PATH_CONFIG=

1.  Make command
    
        go build -a -o app_ .
        ./app_

