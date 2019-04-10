# contacts-api
A simple REST API for contacts written in go. Includes options to export/import contacts via csv file.

**Requirements**
 - Postgres database
 - go installed
 - dep installed
 - some free time

**Installation**<br/><br/>
```go get github.com/SquancherSquanch/contacts-api```

Once the package lands configure the development.yaml file to suit your data base needs.
<br/>
<br/>
  **config key/value description**
  - **host:** baseURL
  - **port:** port
  - **user:** user assigned to the database
  - **password:** password needed to connect to the database
  - **name:** name of the database
  - **db:** name of the table in the data base where the contact entries reside
  
 Ensure your postgres database is running and configured.
 
 Run commands: 
 <br/>
 ```dep ensure```
 <br/>
 ```go run main.go```
 <br/>**or**<br/>
 ```go build``` and start generated file
 
 **End Points**
 <br/><br/>
 **[GET]:**<br/>
 
 **Retrieve list of all contacts<br/>**
   baseurl/api/entry<br/><br/>
 
 **Retrieve a single contact**<br/>
   baseurl/api/entry?id=0<br/>
   *id is an integer that represents an id in the contacts table*<br/><br/>
 
 **Export contacts via csv file**<br/>
   baseurl/api/entry/export<br/><br/>
 
 **[POST]:**<br/>
 
 **Create a new contact**<br/>
   baseurl/api/entry<br/>
   *json data must be provided with this call*<br/>
      **example:**<br/>
      ```{
          "first_name": "tom",
          "last_name": "dob",
          "email": "tom.dobs@gmail.com",
          "phone": "5555555555"
          }
      ```<br/><br/>
 **Import contacts with a csv**<br/>
   baseurl/api/entry<br/>
   *csv file must be provided with headers of [Content-Disposition: form-data; file; filename.csv, Content-Type: text/csv]*<br/><br/>
 
 **[PUT]:**<br/>
 
 **Update contact**<br/>
   baseurl/api/entry<br/>
   *json data must be provided with this call*<br/>
      **example:**<br/>
      ```{
          "id": "4"
          "first_name": "tom",
          "last_name": "dob",
          "email": "tom.dobs@gmail.com",
          "phone": "5555555555"
          }
      ```<br/><br/>
 
 **[DELETE]:**<br/>
 
 **Delete contact**
   baseurl/api/entry?id=0<br/>
   *id is an integer that represents an id in the contacts table*<br/><br/>
 
  




