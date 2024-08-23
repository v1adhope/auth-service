# Draft

REST

`access JWT SHA512`

`refresh base64` store as `bcrypt hash` can be used only once

Refresh valid only for tokens from the same session

`payload: ip`

If ip was changed send allert to user email (mock)

## Get token pair with GUID as req param

POST /tokens/{guid}

```json
{
    "access": "<SOME_TOKEN>",
    "refresh": "<SOME_TOKEN>"
}
```
resp

```errs
internal
not valid guid
```

## Refresh tokens (req `access` and `refresh` from equal session)

POST /tokens/refresh

```json
{
    "access": "<SOME_TOKEN>",
    "refresh": "<SOME_TOKEN>"
}
```
req

```json
{
    "access": "<SOME_TOKEN>",
    "refresh": "<SOME_TOKEN>"
}
```
resp

```errs
internal
not valid token pair
```
