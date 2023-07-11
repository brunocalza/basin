# SQLite to Filecoin Backuper

A CLI app that periodically makes backups of your SQLite database and upload it to Filecoin.

## Installing

```bash
go install github.com/brunocalza/basin/cmd/basin@latest
```

## Configuring

```bash
basin new 
```

## Usage

### Sync your database to a decentralized storage

```bash
basin sync
```

### List your snapshots

```bash
basin list
```

### Restore from a specific snapshot

```bash
basin restore [path]
```

### Status

```bash
basin status [CID]
```
