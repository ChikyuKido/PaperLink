#### List books (prints JSON to stdout):

```bash
./digi4s-downloader list <username> <password>
```
#### Example:

```bash
./digi4s-downloader list myuser mypass > books.json
```
```json
[{
	"name": "buch 1",
	"dataCode": "tedf34r",
	"dataId": "4123"
}]
```

#### Download books (`DataId=output_path` pairs):

```bash
./digi4s-downloader download "ID1=/path/book1.pdf,ID2=/path/book2.pdf" <username> <password>
```

#### Example:

```bash
./digi4s-downloader download "12345=/tmp/book1.pdf,67890=/tmp/book2.pdf" myuser mypass
```
