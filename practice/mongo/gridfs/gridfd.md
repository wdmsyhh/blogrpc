# gridFs

github.com/qiniu/qmgo 是一个用于 MongoDB 的 Go 语言驱动程序库，而 GridFS 是 MongoDB 的一种文件存储机制，用于存储大型二进制文件。

qmgo 库本身不提供对 GridFS 的直接支持，但您可以使用 MongoDB 原生驱动程序（如 go.mongodb.org/mongo-driver）来访问 GridFS 数据。您可以将 qmgo 与 MongoDB 原生驱动程序集成，以同时使用 qmgo 访问文档数据和 mongo-driver 访问 GridFS 数据。