This is a template for modern Web Application with authentification in Go.
Note that it doesn't include any CSS or JS file.
Some parts of codes follow good practices listed on "Go in Practice" written by Matt Butcher.

Here some required setup procedures before running the app.
1.Modify conf.json for server path and gmail account. 
2.Install postgreSQL and execute the following. (This is an example where you create database with user name mark.)
    $ psql -f ~/go/src/Go_in_Practice/FB/data/setup.sql -U mark facebook