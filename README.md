# Filensy

Filensy is a lightweight tool that allow users to save, retrieve and share files in a secure and simple way. It supports client-side encryption, custom file expiration, and public/private file access controls.


## Filensy CLI Guide

### Upload

```bash
filensy upload <file_path> [options]
```
##### Supported flags
--ttl <duration>: To set custom expiration for files. Ex: '15d', '1h'. The default is 30 days if ttl is not provided.

--password or -p <password>: To Encrypt the file based on the provided password.

--public: A bool flag that if enabled would make the file visibility public. File token isn't required to download public files. 

#### Example
```bash
filensy upload "path/to/file" -p "secret" --ttl 1d
```

Output: 

```bash
Upload Complete
File Code:  ABC123
File Access Token:  some-access-token
```
The file code is used to get and manage the file. File Acess Token is required for deleting files from the server and downloading private files.

---
### Download

```bash
filensy get <file_code> [options]
```
##### Supported flags
--dir or -d <directory>: Directory to store the file. Default directory can be configured using the "config" command. In case this isn't provided the current directory will be used.

--password or -p <password>: Password is needed for encrypted files.

--token or -t <token>: Token is needed for private files.
#### Example
```bash
filensy get ABC123 -p "secret" -t "access_token"
```
Note: In case the file was uploaded from the same machine, the password as well as file token will usually be saved on the user's computer. In that case there isn't a need to provide file token and password.

Simply Using
```bash
filensy get ABC123
```
will work.

---
### Find

```bash
filensy find <file_code> [options]
```
Find works exactly similar to "get" but it only shows the metadata about the file if it's available on the server. The "dir" flag isn't needed here since there is no downloading of the file.
#### Example
```bash
filensy find ABC123 -p "secret" -t "access_token"
```

Output Example: 
```bash
File MetaData
Name:  file.txt
Size:  100 KB
Is_Encrypted:  true
Visibility:  private
Uploaded At:  2026-06-11 01:55:37 PM IST
Expires At:  2026-07-11 01:55:37 PM IST
```
---
### Delete

```bash
filensy delete <file_code> [options]
```
Delete also takes the same flags, "password" and "token". It deletes the file permanatly from the server.
#### Example
```bash
filensy delete ABC123 -p "secret" -t "access_token"
```
---
### Config
The config command is used to configure global cli options. It takes a config key and the config value.
```bash
filensy config <key> <value>
```
Available configuration keys are:

"auth": This needs to be configured in case server auth is enabled. It attaches the "authorization" header for every request.
"dir": This is used to configure the default directory where files would be downloaded.
#### Example
```bash
filensy config dir "C:\Users\user\Downloads"
```