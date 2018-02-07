# lcache
golang实现一个基于内存的key-value缓存，模仿redis，简单实现了几个命令,最先想用gob，将数据编码存储，但是在解码时，碰到了一个问题，一直没有解决，所以对存储的数据没有经过编码

安装:

go get github.com/jfeige/lcache

