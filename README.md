# GeoBasedService
Copy these Three folders into ~/go/src/ folder
Go to Individual folder and run go get ./...  to get the dependencies then run go build
Example
  Go to GeoLoginServer
  run go get ./... on terminal
  run go build
  you will  get a build which you can run
  
API's

http://localhost:6000/v1/checkuser

Input:
{
    "phoneno":"1234567881"


}

Output:
{
    "error": "Username already exist",
    "success": true
}




http://localhost:6000/v1/register

Input:

{
    "username":"Neerajs",
    "email":"abc@d.coo",
    "phoneno":"1234567881",
    "password":"passwrd",
    "long":"40.20",
    "lat":"5.20"
    
}

Output:

{
    "result": "Registration Successful",
    "success": true
}



http://localhost:6000/v1/login

Input:

{
    "phoneno":"1234567881",
    "password":"passwrd",
    "long":"48.22",
    "Lat":"5.22",
    “devid”:”deviceid”

    
}

Output:


{
    "data": {
        "userid": "5df342d121c07acd420b43c0",
        "username": "Neerajs",
        "email": "abc@d.coo",
        "phoneno": "1234567881",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkZXZpZCI6IiIsImZpcnN0bmFtZSI6IiIsImxhc3RuYW1lIjoiIiwicGhvbmVubyI6IjEyMzQ1Njc4ODEifQ.c4T03lA7-r-cf5HVC96gdiJih5_asORfFx4c_eaR_IY",
        "location": {
            "type": "Point",
            "coordinates": [
                48.22,
                5.22
            ]
        }
    },
    "success": true
}



http://localhost:6001/v1/updatelocation

Input:

{
    "UID":"5df342d121c07acd420b43c0",
    "long":"49.22",
    "lat":"5.22"
    
}

Output:

{
    "success": true,
    "result": "Data Updated successfully!!"
}



http://localhost:6001/v1/getpeople

Input:

{
"page":1,
"Lat":"5",
"Long":"48",
"UID":"5df34094eaa0ada3690ff7a"
}

Output:

{
    "data": [
        {
            "userid": "5df342d121c07acd420b43c0",
            "username": "Neerajs",
            "email": "abc@d.coo",
            "phoneno": "1234567881",
            "password": "$2a$05$84eMnS0qzgqrosqf2Kt43ulr4jGPn36fgrPA4aX9ZSILWXWoJfZEW",
            "location": {
                "type": "Point",
                "coordinates": [
                    49.22,
                    5.22
                ]
            },
            "Dist": 137.4681832270385
        }
    ],
    "error": "",
    "page": 1,
    "success": true,
    "total_pages": 1
}









Forgot Password:
http://localhost:6000/v1/forgot
Method:  POST
Sample:   {
    "Email":                            //This can be email or phone number,
    "Password":"12345678"
}
Output success::{
    "error": "",
    "result": "Password changed successfully"
}

Change Password:
http://localhost:6000/v1/change
Method:  POST
Sample:   {
    "UID":                            //UserID,
    "Password":"New Password",
    “Oldpassword”:”Old Password”
}
Output success::{
    "error": "",
    "result": "Password changed successfully"
}





http://localhost:6001/v1/uploademergency

{"UID":"5df342d121c07acd420b43c0",
"Ename":"aaa1",
"Ephonenumber":"bcd1"

}

{
    "success": true,
    "result": "Uploaded successfully!!"
}


http://localhost:6001/v1/UpdateProfile

{"UID":"5df342d121c07acd420b43c0",
"Email":"aaa1",
"Username":"bcd1"

}

{
    "success": true,
    "result": "Uploaded successfully!!"
}


http://localhost:6001/v1/uploadimage

Type is the factor if you leave it black it will become user profile image otherwise it will be aadhaar 

{
    "success": true,
    "result": "Image Uploaded successfully!!"
}



http://localhost:6001/v1/getprofile
{“Userid”:””}

http://localhost:4000/v1/download{userimageurl}



