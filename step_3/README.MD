# Step 3 - Build Book Quantity Tracker

## Studi Kasus

Cindy merupakan seorang pustakawan di Perpustakaan Alpha Beta. Dalam menjalankan kegiatan sehari-harinya, Cindy biasa melakukan kegiatannya secara manual, misal untuk mencatat peminjam buku, melacak jumlah buku tersedia setiap saat, dan mengatur daftar koleksi buku yag tersedia di perpustakaan tersebut.

Di perpustakaan tersebut, masing-masing anggota diharuskan untuk mendaftar keanggotaan di perpustakaan tersebut apabila ingin meminjam koleksi buku yang ada. Namun karena semuanya masih manual, orang-orang yang ingin mendaftar masih harus antri untuk daftar ke Cindy, dan identitas anggota perpustakaan masih menggunakan sebuah kartu anggota. Terlebih Cindy harus mencatat peminjam dalam buku catatan peminjam buku yang mana kegiatan tersebut melelahkan dan bila perpustakaan sedang ramai, bisa jadi Cindy bisa melewatkan beberapa peminjam karena terlalu banyaknya orang yang antri meminjam.

Maka dari itu, anda sebagai seorang konsultan IT dan seorang developer datang kepada Cindy untuk membantu Cindy melakukan digitalisasi Perpustakaan Alpha Beta. Ketika dilakukan wawancara terhadap Cindy, terdapat beberapa cerita yang diberikan oleh Cindy terkait keperluan digitalisasinya:

1. Koleksi buku yang ada di perpustkaan biasanya memiliki data-data sebagai berikut:
    1. ID Buku
    2. Judul Buku
    3. Pengarang Buku
    4. Penerbit Buku
    5. Ringkasan Buku
    6. Stok Buku
    7. Stok Maksimal Buku

2. Koleksi buku di perpustakaan bisa dilacak oleh Cindy, selain itu Cindy dapat melakukan penambahan koleksi buku, penghapusan koleksi buku, atau melakukan pembaruan terhadap informasi koleksi buku.

3. Anggota yang ingin mendaftar keanggotaan di perpustakaan secara digital perlu memberikan informasi sebagai berikut:
   1. User Name.
   2. Password.
   3. Email.

    Dan apabila user ingin melakukan login keanggotaannya untuk meminjam, anggota tersebut perlu login dengan data-data sebagai berikut:
    1. User Name.
    2. Password.

4. Ketika seorang anggota perpustakaan meminjam buku, biasanya anggota tersebut membawa bukunya ke Cindy, lalu Cindy akan memutuskan apakah buku tersebut boleh dipinjam atau tidak. Jika tidak boleh dipinjam karena alasan tertentu, Cindy akan mengarahkan anggota tersebut untuk meminjam buku yang lain beserta memberi tahu alasannya. Namun bila buku tersebut boleh dipinjam, Cindy akan mencatat kegiatan peminjaman tersebut dengan mencatat beberapa data sebagai berikut:
   1. ID buku yang dipinjam
   2. ID anggota yang meminjam.
   3. Waktu dan tanggal peminjaman.
   4. Tanggal tenggat buku harus dikembalikan.
   5. Status buku sudah dikembalikan atau belum.
   6. Tanggal buku dikembalikan.

5. Ketika ada permintaan peminjaman, Cindy perlu melakukan persetujuan apakah buku tersebut boleh dipinjam atau tidak. Sehingga semua permintaan peminjaman buku harus dicatat pada suatu daftar catatan tersebut agar Cindy bisa melakukan crosscek apakah anggota tersebut benar-benar ingin meminjam atau tidak sebelum Cindy memberikan persetujuan. Cindy menginginkan catatan tersebut memiliki data sebagai berikut:
    1. ID buku yang ingin dipinjam.
    2. ID anggota yang ingin meminjam.
    3. Waktu dan tanggal permintaan peminjaman.
    4. Status permintaan diterima atau tidak.
    5. Alasan apabila permintaan ditolak.

    Dengan catatan tersebut, Cindy bisa mengetahui dan melacak kapan buku ini ingin dipinjam dan misalnya ditanya kenapa buku ini tidak bisa dipinjam pada waktu tertentu, anggota perpustakaan bisa meninjau alasannya di lain waktu.

6. Alur peminjaman buku di Perpustakaan Alpha Beta ketika dilakukan digitalisasi menurut Cindy adalah sebagai berikut:
   1. Anggota melakukan login ke website Perpustakaan Alpha Beta.
   2. Anggota melihat daftar koleksi buku yang tersedia di perpustakaan.
   3. Anggota melakukan permintaan peminjaman buku selama persediaan buku masih ada, yang menyebabkan permintaan peminjamannya tercatat di catatan permintaan peminjaman buku
   4. Anggota tersebut menuju Cindy dan memberi tahu user name keanggotaannya kepada Cindy dan judul buku yang ingin dipinjam. Kemudian Cindy melakukan crosscek ke catatan permintaan peminjaman buku. Apabila tidak ditemukan, Cindy tidak melanjutkan proses peminjaman buku, namun bila ditemukan, Cindy akan memutuskan apakah peminjaman tersebut disetujui atau tidak.
   5. Apabila disetujui, Cindy akan melakukan persetujuan, yang menyebabkan data "status permintaan diterima" di catatan permintaan peminjaman buku berubah menjadi DISETUJUI, lalu data peminjaman buku yang baru tersebut akan tercatat di catatan peminjaman buku dengan waktu tenggat yang ditentukan oleh Cindy, biasanya 3 hari. Kemudian data stok buku akan diperbarui dan Cindy memberikan bukunya untuk dipinjam oleh anggota yang meminjam.
   6. Apabila tidak disetujui, Cindy akan menolak permintaan tersebut, yang menyebabkan data "status permintaan diterima" di catatan permintaan peminjaman buku berubah menjadi DITOLAK, kemudian Cindy akan mengisi alasan kenapa permintaan tersebut ditolak.

7. Alur pengembalian buku menurut Cindy dapat dilakukan sebagai berikut:
   1. Anggota menemui Cindy untuk mengembalikan buku.
   2. Anggota memberikan buku yang dipinjam dan memberitahukan user name-nya kepada Cindy.
   3. Cindy melakukan pengecekan di catatan peminjaman buku, lalu melakukan perubahan data bahwa buku sudah dikembalikan.
   4. Ketika buku sudah dikembalikan, sistem akan mengeluarkan denda yang harus dibayar oleh peminjam apabila terdapat perbedaan waktu antara "waktu buku harus dikembalikan" dan "waktu pengembalian buku dilakukan" sebanyak (n x Rp. 5000) apabila memang pengembaliannya telat. Waktu telat dihitung apabila sudah ganti hari dari "waktu buku harus dikembalikan".
   5. Ketika buku sudah dikembalikan, sistem akan memperbarui stok buku.

## Tugas

Berdasarkan studi kasus tersebut, buatlah interface-interface yang diperlukan untuk menjalankan semua bisnis proses tersebut.

Apabila terdapat hal yang kurang jelas, dapat langsung ditanyakan ke mentornya ya
