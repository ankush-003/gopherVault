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

## Atomic Renaming

Many problems are solved by not updating data in-place. You can write a new file and delete the old file. Not touching the old file data means:

If the update is interrupted, you can recover from the old, intact file.
Concurrent readers won’t get half written data.
How will readers find the new file? This is solved by the renaming pattern:

```go
func SaveData2(path string, data []byte) error {
    tmp := fmt.Sprintf("%s.tmp.%d", path, randomInt())
    fp, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
    if err != nil {
        return err
    }
    defer func() { // 4. discard the temporary file if it still exists
        fp.Close() // not expected to fail
        if err != nil {
            os.Remove(tmp)
        }
    }()

    if _, err = fp.Write(data); err != nil { // 1. save to the temporary file
        return err
    }
    if err = fp.Sync(); err != nil { // 2. fsync
        return err
    }
    err = os.Rename(tmp, path) // 3. replace the target
    return err
}
```

Renaming a file to an existing one replaces it atomically. But pay attention to the meaning of the jargon, whenever you see “X is atomic”, you should ask “X is atomic with respect to what?” In this case:

Rename is atomic w.r.t. concurrent readers; a reader opens either the old or the new file.
Rename is NOT atomic w.r.t. power loss; it’s not even durable. You need an extra fsync on the parent directory, which is discussed later.
Why does renaming work?
Filesystems keep a mapping from file names to file data, so replacing a file by renaming simply points the file name to the new data without touching the old data. That’s why atomic renaming is possible in filesystems. And the operation cost is constant regardless of the data size.

On Linux, the replaced old file may still exist if it’s still being opened by a reader; it’s just not accessible from a file name. Readers can safely work on whatever version of the data it got, while writer won’t be blocked by readers. However, there must be a way to prevent concurrent writers. The level of concurrency is multi-reader-single-writer, which is what we will implement.

## Append Only Logs

One way to do incremental updates is to just append the updates to a file. This is called a “log” because it’s append-only. It’s safer than in-place updates because no data is overwritten; you can always recover the old data after a crash.

The reader must consider all log entries when using the log. For example, here is a log-based KV with 4 entries:

| Operation | Key | Value |
|-----------|-----|-------|
| set       | a   | 1     |
| set       | b   | 2     |
| set       | a   | 3     |
| del       | b   | -     |

The final state is a=3.

Logs are an essential component of many databases. But logs alone are not enough to build a DB because:

It’s not an indexing data structure; readers must read all entries.
It has no way to reclaim space from deleted data.