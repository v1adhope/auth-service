# Endpoints

## Generate tokens

POST /tokens/{guid}

```json
{
    "accessToken": "<SOME_TOKEN>",
    "refreshToken": "<SOME_TOKEN>"
}
```
resp


## Refresh tokens

POST /tokens/refresh

```json
{
    "accessToken": "<SOME_TOKEN>",
    "refreshToken": "<SOME_TOKEN>"
}
```
req body

```json
{
    "accessToken": "<SOME_TOKEN>",
    "refreshToken": "<SOME_TOKEN>"
}
```
resp
