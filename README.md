# Web Auth using mail authentication
This is a template written in Go, which you can use for secure web authentication using mail authentication.
Some parts of codes follow good practices listed on "Go in Practice" written by Matt Butcher.


# DEMO
Menu Screen
![Image of Home Screen](https://github.com/froprintoai/modernWebAuth/blob/master/home.png?raw=true)

# Prerequisite
1.Modify conf.json for server path and gmail account. 
2.Install postgreSQL and execute the following. (This is an example where you create database with user name mark.)
```
    $ psql -f ~/go/src/Go_in_Practice/FB/data/setup.sql -U mark facebook
```

