# pure-webserver

Pure Golang web service without any 3rd party libraries.


## Inner libs:

### database

simple json based database contains base required functionalities to use 
```go
WriteToCollection(&YOUR_MODULE)
GetFromCollection(&YOUR_MODULE) 
UpdateCollection(&YOUR_MODULE) 

// search on a query on single data (document).
Where(fieldName string, value interface{}) 
Update(&YOUR_MODULE) (*DBInnerModel, error)
All() *DBInnerModel
```


### HTTP engine:
simple HTTP engine developed on net/http package which supports in query params([example](https://github.com/amupxm/pure-webserver/blob/main/controller/httpEngine.go#L35)).

you can add your handler like this :
```go
func (e *engine) GetOne(c *httpEngine.ServerContext) {
	// check iid exists or not
	id, err := c.GetURLParam("iid")
	if err != nil {
		c.ErrorHandler(400, err)
		return
	}
	ee, err := e.ProductLogic.GetProductByID(id)
	if err != nil {
		c.ErrorHandler(400, errors.New(constants.NoData))
		return
	}
	c.JSON(200, ee)

}
```
## Usage
Use your system : be sure you had Golang compiler install on your device
```bash
touch database.json 

# to run the app
go run main.go

#to build the app
go build .
```
Use docker :
```bash
docker build --tag youruser/yourtag .
docker run  -p 8080:8080 youruser/yourtag 
```
**attention:** you can change bucket name and porn from `config/config.json`

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

