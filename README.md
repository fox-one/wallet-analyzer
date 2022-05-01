# wallet-analyzer

A Command-line tool that lists transaction snapshots of mixin application, allowing you to easily check and analyze accounts.

## Install

### go version <= 1.16

```
go get github.com/fox-one/wallet-analyzer
```

### else

```
go install github.com/fox-one/wallet-analyzer@latest
```

## Usage

```
$ wallet-analyzer -h
Usage of wallet-analyzer:
  -asset string
        Asset id
  -client string
        Mixin client id
  -end string
        End time, RFC3339 format
  -format string
        Snapshot format (default "id: {{ .SnapshotID }} -> (asset: {{ .AssetID }}, amount: {{ .Amount }})")
  -opponent string
        Opponent id
  -output string
        Output file path
  -secret string
        Mixin client secret
  -start string
        Start time, RFC3339 format
  -token string
        Access token

```

### Examples

#### Invoke by client and secret

If you don't have an existing token, you can set client and secret, it will open your browser and start the OAuth flow, and once the authentication is completed then you can copy the code from the address bar.

```
$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -client=$YOUR_CLIENT -secret=$YOUR_SECRET
OAuth Code: 8024fbb034283008e0d7cb60c23c083aa0c6ceae88cd7fd4d4cbf24f602ecce5

token: eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIzYjUwNmI4OS00MTVkLTQ1NTAtYTZjMS05ZTdmODk1NDk5NmMiLCJleHAiOjE2ODI5NTQyOTksImlhdCI6MTY1MTQxODI5OSwiaXNzIjoiYTQ4YjFlNfEtOWI1ZS00NzBkLWJjNzAtOTc0ZDE3ZWExNjxjIiwic3NwFjxiUFJPRklMRTpSRUFEIFNOQVBTSE9UUzpSRUFEIn0.HrLev8DpPkJfMSpsXvDhEDsl51UVnmYnNAlMSbblXvO8He8GiFzLawLEqL8O6Kh9oWuEIfcs-SUF9Hb_eLmcP_g4ULpLmiMRBkpM9i7d_7hbAgkKHxJHAGU7RW6hebE0BHXoJYm2nCwnZGu49Xn9-3ayAsuf0gL9Ben7jz1j72o

id: d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -> (asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b, amount: 9)

ids: ('d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c')

asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -> (count: 1, total: 9)

```

#### Invoke by token

As you can see, if you call through the client and secret, the output will print the generated token, save it, and then you can use the token call to avoid repeated authentication.

```
$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -token=$YOUR_TOKEN
id: d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -> (asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b, amount: 9)

ids: ('d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c')

asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -> (count: 1, total: 9)

```

#### Apply filters

You can set the start time, end time, asset id and opponent id to filter the snapshot results.

```
$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -token=$YOUR_TOKEN

$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -opponent=d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -token=$YOUR_TOKEN

$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -opponent=d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -token=$YOUR_TOKEN -start=2022-03-23T03:27:27Z -end=2022-05-01T03:27:27Z
```

#### Custom output snapshot format

By default, only the snapshoot id, asset id and amount are output, you can change the behavior by passing the format parameter, it's golang [template syntax](https://pkg.go.dev/text/template), the available data you can find in [Snapshoot Record](https://github.com/fox-one/mixin-sdk-go/blob/master/snapshot.go).

```
$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -token=$YOUR_TOKEN -format="id: {{ .SnapshotID }} -> (asset: {{ .AssetID }}, amount: {{ .Amount }}, type: {{ .Type }}, created: {{ .CreatedAt }})"
id: d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -> (asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b, amount: 9, type: deposit, created: 2022-01-18 02:59:58.487811 +0000 UTC)

ids: ('d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c')

asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -> (count: 1, total: 9)

```

#### Save output to file

Of course, you can use pipelines to save the output, but it will save unnecessary info, such as token, you can avoid this by using output param.

```
$ wallet-analyzer -asset=b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -client=$YOUR_CLIENT -secret=$YOUR_SECRET -output=./result.txt
OAuth Code: 1b7581510ccd465b3ccd246dfd03369cecf265d43ae12d5b79dcf62080057ac7

token: eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhaWQiOiIyNWY4OWIxYi0zMDhjLTQxZjYtOWJmMS0yNzViZDRmNzgwOTgiLCJleHAiOjE2ODI5NTgyNzUsImlhdCI6MTY1MTQyMjI3NSwiaXNzIjoiYTQ4YjFlNTEtOWI1ZS00NzBkLWJjNzAtOTc0ZDE3ZWExNjJjIiwxc2NwIjoiUFJPRklMRTpSRUFEIFNOQVBTSE9pzZLzaUFEzn1.BJS9WEe8iWQdoKTnDWXEa15P62Fjl0QZy6-NG_4OJjVZ96tVdMdNXZs9XQIHYpO2dzeIkXyVmVDUHXy7YsinRe6iEMa0puT58Htu57zk9Ybeazs71jrlhMlMm-5_3hkSVvLhkRb3hk-5Q79w2oEWuRESfxhcl-Cq06jdvgEKfLw

$ cat result.txt
id: d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c -> (asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b, amount: 9)

ids: ('d1f38bfc-f60e-4a6c-a588-0e2f8ad4297c')

asset: b91e18ff-a9ae-3dc7-8679-e935d9a4b34b -> (count: 1, total: 9)


```
