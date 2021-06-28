# iitk-coin
The project has the work carried out as a summer project offered by Programming Club, IIT Kanpur.

## Aim of the project
The project aims to build a pseudo-currency for use in the IITK Campus.

## How to run the code locally?
- After cloning the repository, run the following code in the ./iitk-coin directory
```
go run main.go
```
## Testing different endpoints
### Register:
- To register for a new user, create a new "data.json" file in any commond running directory, having the following format:
```
{
  "rollno":rollno_integer,
  "name":"name_string",
  "password":"password_string",
  "access":"access_string" // give "gensec" for authority access and "student" for student access
  }
```
- Next run the following code from a command line shell-
```
curl --location --request POST "http://localhost:8080/signup"  --header "Content-Type: application/json" -d @data.json 
```
### Login:
- To login and get a JWT Token create a new "data.json" file in any commond running directory, having the following format:
```
{
  "rollno":rollno_integer,
  "password":"password_string"
  }
```
- Next run the following code from a command line shell-
```
curl --location --request POST "http://localhost:8080/login"  --header "Content-Type: application/json" -d @data.json 
```
- The command line would run a JWT Token which can be used to access coin-related endpoints
### Award Coins:
- Can only be accessed by an authority, in my example a "gensec" access
- To award coins to a user create a new "data.json" file in any commond running directory, having the following format:
```
{
  "rollno":rollno_integer,
  "awarded":no_of_coins_integer
  }
```
- Next run the following code from a command line shell-
```
curl --location --request POST "http://localhost:8080/awardcoins"  --header  "Authorization: valid_jwt_token"  --header "Content-Type: application/json" -d @data.json 
```
### Get Coins:
- Returns the total number of coins a user currently has
- Run the following code from a command line shell-
```
curl --location --request GET "http://localhost:8080/getcoins"  --header  "Authorization: valid_jwt_token"  --header "Content-Type: application/json"
```
### Transfer Coins:
- Transfer coins from one user to another 
- Create a new "data.json" file in any commond running directory, having the following format:
```
{
  "reciever":reciever_rollno_integer,
  "coins":transfer_amount_integer
 }
```
- Next run the following code from a command line shell-
```
curl --location --request POST "http://localhost:8080/transfercoins"  --header  "Authorization: valid_jwt_token"  --header "Content-Type: application/json" -d @data.json 
```
The sender's jwt token should be passed in valid_jwt_token

### Records:
- Get the history of transactions of a user.
- Run the following code from a command line shell-
```
curl --location --request POST "http://localhost:8080/records"  --header  "Authorization: valid_jwt_token"  --header "Content-Type: application/json"
```
