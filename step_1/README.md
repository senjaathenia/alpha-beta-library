# Step 1 - Build User Authentication API

## Feature lists

1. Data user terdiri dari:
   1. Unique ID User
   2. User Name
   3. Email
   4. Password
   5. Created At
   6. Updated At
   7. Deleted At

2. CUD (Create-Update-Delete) data user.
   1. Definisi Create:
      1. *Request* yang diterima:

           ```json
           {
               "username": "string",
               "email": "string",
               "password_1": "string",
               "password_2": "string"
           }
           ```

      2. *Unique ID* di-*generate* secara otomatis, boleh di-*generate* oleh aplikasi API ataupun RDBMS yang digunakan.

      3. *Password* disimpan dalam BCrypt hash.

      4. Data user terbentuk apabila `password_1` dan `password_2` sama.

      5. Validasi password: minimal 8 huruf, alfanumerik + simbol, minimal memiliki 1 huruf besar, minimal memiliki 1 angka, dan minimal memiliki 1 simbol.

   2. Definisi Update:
       1. *Request* yang diterima:

           ```json
           {
               "username": "string",
               "email": "string",
               "password_1": "string",
               "password_2": "string"
           }
           ```

       2. `Username` wajib diisi, `email` dan `password_1` & `password_2` opsional.

       3. Data yang berubah sesuai dengan data yang diisi, bisa email saja, atau bisa password saja.

   3. Definisi Delete:
       1. *Request* yang diterima:

           ```json
           {
               "user_id": "string",
           }
           ```

       2. User ID disesuaikan dengan *Unique ID* yang telah ditentukan sebelumnya.

3. Validasi Username & Password.
    1. *Request* yang diterima:

        ```json
        {
            "username": "string",
            "password": "string",
        }
        ```

    2. Tampilkan pesan ke API client (*front end*) apabila username dan password valid atau tidak.

4. File SQL migration harus dibuat dan dimasukkan ke dalam folder `migrations`.
