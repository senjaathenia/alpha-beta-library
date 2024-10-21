# Step 2 - Build User Authentication Middleware

## Feature lists

1. Buat satu endpoint bebas _request-response_-nya. Misal `GET http://localhost:8080/protected/hello`.

2. Buatlah endpoint `POST /login` yang berfungsi untuk memvalidasi username + password dengan request _username_ & _password_ sesuai yang telah dibuat di step 1.
   Namun response yang dikeluarkan merupakan _response_ token hasil pembangkitan JWT.
   Contoh:

   ```json
    {
        "message" : "Successful Login",
        "data": {
		    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0NjE5NTcxMzZ9.RB3arc4-OyzASAaUhC2W3ReWaXAt_z2Fd3BN4aWTgEY",
        },
        "errors": null,
        "code": 200
    }
   ```

   _Secret Key_ untuk JWT dapat dibuat bebas, namun disarankan merupakan string random.

3. Buat sebuah middleware yang berfungsi untuk melakukan autentikasi endpoint yang dibuat pada fitur nomor 1 dengan mengecek _Authorization Header_ dari request HTTP menuju endpoint yang telah dibuat. Apabila token yang dikirim tidak valid, API _client_ tidak boleh mengakses endpoint tersebut.

## Referensi

1. [JWT](https://jwt.io/)
2. [Echo Cookbook](https://echo.labstack.com/cookbook/jwt/)