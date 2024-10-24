# Library Management Project

## *What is this?*

Di project kali ini, kita akan membangun sebuah project aplikasi yang akan digunakan untuk mengatur daftar buku perpustakaan menggunakan Go Programming Language.

## *Any details*

1. Aplikasi dibangun dengan menggunakan Go Programming Language.

2. Aplikasi dibuat secara *full stack* (*Front End* & *Back end* digabung menjadi satu).

3. *Source code* aplikasi di push di dalam folder ini, dengan cara membuat folder baru untuk masing-masing *mentee* dan *push source code*-nya di sana.

4. Pengetahuan tambahan yang perlu diasah:
   1. HTML
   2. CSS
   3. JavaScript (ES2015)
   4. RDBMS (disarankan PostgreSQL)
   5. *Bonus tambahan*: Docker

5. Library yang diperbolehkan untuk digunakan:
   1. Sebuah *web framework*. (Echo Web Framework)
   2. ORM. Gunakan **GORM V2**.
   3. JSONIter untuk marshal & unmarshal JSON.
   4. [Swag](https://github.com/swaggo/swag) (Swagger API Documentation).
   5. Any of Go built in library.
   6. [Testify](https://github.com/stretchr/testify).
   7. Library tambahan di luar list ini yang ingin digunakan bisa ditanyakan terlebih dahulu.

6. Library yang tidak diperbolehkan untuk digunakan:
   1. jQuery.

7. Struktur folder dibebaskan ke *mentee*, namun sangat disarankan untuk mengikuti kaidah [Go Standard Project Layout](https://github.com/golang-standards/project-layout).

8. Arsitektur kode dibebaskan ke *mentee*, namun sangat disarankan untuk mengikuti kaidah [Go Clean Architecture](https://github.com/bxcodec/go-clean-arch).

## Common Details

### API Details

1. API dibangun dengan kaidah REST API.

2. Response dari sisi API seminimalnya memiliki data *field* sebagai berikut:

    ```json
    {
        "message": "string",
        "data": "object, nullable",
        "errors": "object, nullable",
        "code": "int"
    }
    ```

3. *message* berisikan pesan hasil proses dari aplikasi bagian *back end*, apakah proses yang dilakukan berhasil, gagal, dan sebagainya.

4. *data* berisikan data hasil proses dari aplikasi bagian *back end*. Diisi `null` bila tidak ada data-nya.

5. *errors* berisikan satu atau lebih `error` yang terjadi pada saat aplikasi bagian *back end* melakukan suatu proses. Diisi `null` bila tidak ada `error`-nya.

6. *code* berisikan http status code hasil dari status code proses yang didapatkan dari aplikasi bagian *back end*.

7. Semua kalimat di dalam respon API harus ditulis dalam bahasa inggris.

### Front End Details

1. Bagian *Front end* harus memanfaatkan Go Template.

2. Kode JavaScript harus *full native*, tidak boleh menggunakan bantuan jQuery.

3. *Front end* dibangun dalam Bahasa Indonesia formal.

### Code Development Details

1. Setiap code yang di push harus lolos *linter*. Gunakan [golangci-lint](https://golangci-lint.run/).

2. Setiap fitur harus dikerjakan dalam *branch* yang berbeda. *Checkout* dari latest *main*.

3. Semua kode yang akan di push harus dibuat PR terlebih dahulu.

4. Semua commit harus mengikuti kaidah [Conventional Commits](https://www.conventionalcommits.org/id/v1.0.0/).