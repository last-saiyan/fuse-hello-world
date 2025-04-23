# fuse-hello-world
simple in memory fs using FUSE



```
mkdir temp
go run main.go temp/
```

in a different terminal, we can use the fs 

```
mkdir temp/asdf
touch temp/asdf/file

```
this fs exists only till the go program is running, if process is terminated, its gone
