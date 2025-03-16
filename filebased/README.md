# File as a Database

Building database using files

```go
func SaveData1(path string, data []byte) error {
    fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
    if err != nil {
        return err
    }
    defer fp.Close()

    _, err = fp.Write(data)
    if err != nil {
        return err
    }
    return fp.Sync() // fsync
}
```

This code creates the file if it does not exist, or truncates the existing one before writing the content. And most importantly, the data is not persistent unless you call fsync (fp.Sync() in Go).

It has some serious limitations:

- It updates the content as a whole; only usable for tiny data. This is why you don’t use Excel as a database.
- If you need to update the old file, you must read and modify it in memory, then overwrite the old file. What if the app crashes while overwriting the old file?
- If the app needs concurrent access to the data, how do you prevent readers from getting mixed data and writers from conflicting operations? That’s why most databases are client-server, you need a server to coordinate concurrent clients. (Concurrency is more complicated without a server, see SQLite).
