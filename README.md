# project-golang-crud

##### The diagram 
![clean_code_architecture](./clean-code-arch.jpg)

More explanation about clean code architecture can read from this
medium's post :
> https://medium.com/@imantumorang/golang-clean-archithecture-efd6d7c43047.

## Deskripsi Proyek: Sistem Manajemen Perpustakaan

#### Latar Belakang

Anda ditugaskan untuk mengembangkan aplikasi manajemen perpustakaan menggunakan bahasa pemrograman Go. Aplikasi ini harus mengikuti prinsip-prinsip arsitektur bersih (clean architecture) dan memanfaatkan framework Echo untuk pengelolaan HTTP request. Data akan disimpan di database PostgreSQL dan diakses melalui GORM, sebuah Object-Relational Mapping (ORM) untuk Go.

#### Tujuan
Membangun aplikasi RESTful API untuk manajemen buku dalam sebuah perpustakaan yang meliputi operasi CRUD (Create, Read, Update, Delete) dengan menggunakan arsitektur bersih dan framework Echo.

####  Fitur yang Diharapkan

Manajemen Buku:

* CRUD Buku: Menyediakan endpoint untuk membuat, membaca, memperbarui, dan menghapus buku.
* Detail Buku: Menyimpan informasi seperti ID, judul, penulis, penerbit, serta tanggal pembuatan dan pembaruan.

Konfigurasi dan Koneksi Database:

* Menggunakan file konfigurasi .env untuk menyimpan detail koneksi database.
* Menggunakan GORM untuk interaksi dengan database PostgreSQL.

Dokumentasi API:

* Menyediakan dokumentasi API dalam format Swagger/OpenAPI.

### Struktur Proyek

1. conf/
* Deskripsi: Folder ini menyimpan file konfigurasi untuk aplikasi, seperti config.env, yang berisi variabel lingkungan seperti URL database.
* File: config.env: Berisi variabel lingkungan seperti seperti detail koneksi untuk menghubungkan ke database PostgreSQL.

2. domains/
* Deskripsi: Folder ini berisi definisi model atau struktur data yang digunakan dalam aplikasi. Model ini biasanya merepresentasikan entitas dari domain aplikasi dan juga interface untuk repository dan usecase.
* File: models.go: Berisi struktur model User yang dikelola oleh GORM dan digunakan dalam operasi CRUD.

3. pkg/
* Deskripsi: Folder ini berisi modul-modul aplikasi yang mengimplementasikan arsitektur bersih. Terdapat tiga subfolder utama:

4. pkg/config/
* Deskripsi: Berisi konfigurasi dan inisialisasi aplikasi, termasuk pengaturan database.
* File: config.go: Menyediakan fungsi Init untuk menginisialisasi koneksi database.

5. pkg/delivery/
* Deskripsi: Folder ini menangani penerimaan dan pengolahan HTTP request. Di sini didefinisikan handler HTTP untuk berbagai endpoint.
* File: handler.go: Mengimplementasikan handler untuk CRUD operations, termasuk pembuatan, pengambilan, pembaruan, dan penghapusan data.

6. pkg/usecase/
* Deskripsi: Folder ini berisi logika bisnis dan aturan aplikasi. Usecase menghubungkan antara delivery layer dan repository layer.
* File: usecase.go: Mendefinisikan interface UserUsecase dan implementasinya, yang memproses logika aplikasi untuk operasi CRUD.

7. pkg/repository/
* Deskripsi: Folder ini mengelola akses data dan integrasi dengan database. Repository berfungsi sebagai lapisan akses data dan menyediakan fungsi untuk berinteraksi dengan database.
* File: repository.go: Implementasi fungsi CRUD untuk model User yang berkomunikasi dengan database menggunakan GORM.

8. swagger/
* Deskripsi: Folder ini berisi file dokumentasi API dalam format Swagger/OpenAPI. Dokumen ini mendeskripsikan API dan endpoint yang tersedia untuk aplikasi.
* File: swagger.yaml: Dokumentasi API menggunakan Swagger/OpenAPI.

#### Instruksi Pengerjaan
1. Inisialisasi Proyek:
    * Buat folder dan file sesuai dengan struktur proyek di atas.
    * Implementasikan model dan interface di domains/models.go.
    * Implementasikan konfigurasi dan koneksi database di pkg/config/config.go.

2. Implementasi Repository:

    * Implementasikan akses data CRUD di pkg/repository/repository.go.

3. Implementasi Usecase:

    * Implementasikan logika bisnis CRUD di pkg/usecase/usecase.go.

4. Implementasi Handler:

    * Implementasikan endpoint HTTP untuk CRUD operations di pkg/delivery/handler.go.

5. Dokumentasi API:

    * Buat dokumentasi API di swagger/swagger.yaml.

6. Pengujian:

    * Jalankan aplikasi dan uji setiap endpoint untuk memastikan semuanya berfungsi dengan baik.
    
7. Dokumentasi dan Pembersihan:
    
    * Pastikan semua file terstruktur dengan baik dan dokumentasi API up-to-date.
