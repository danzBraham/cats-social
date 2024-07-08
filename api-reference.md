
# Cats Social API Reference

Welcome to the Cats Social API! This API allows cat owners to manage their cats and match them with other cats. Below are the available endpoints and their respective functionalities.

## Authentication and Authorization

#### Register user

`POST /v1/user/register`

Request:

```json
{
  "name": "frontname lastname", // minLength 5, maxLength 50
  "email": "your@email.com", // should be in email format
  "password": "secret" // minLength 5, maxLength 15
}
```

Response:

```json
{
  "message": "User registered successfully",
  "data": {
    "email": "your@email.com",
    "name": "frontname lastname",
    "accessToken": "random token"
  }
}
```

- `201` User successfully registered
- `409` conflict if email exists
- `400` request doesn’t pass validation
- `500` if server error

#### Login user

`POST /v1/user/login`

Request;

```json
{
  "email": "your@email.com", // should be in email format
  "password": "secret" // minLength 5, maxLength 15
}
```

Response:

```json
{
  "message": "User logged successfully",
  "data": {
    "email": "your@email.com",
    "name": "frontname lastname",
    "accessToken": "random token"
  }
}
```

- `200` User successfully logged
- `404` if user not found
- `400` if password is wrong
- `400` request doesn’t pass validation
- `500` if server error
## Managing Cats

> [!WARNING]
> All request here should use Bearer Token from accessToken auth route

#### Create cat

`POST /v1/cat`

Request:

```json
{
  "name": "", // minLength 1, maxLength 30
  "race": "" /** enum of:
			- "Persian"
			- "Maine Coon"
			- "Siamese"
			- "Ragdoll"
			- "Bengal"
			- "Sphynx"
			- "British Shorthair"
			- "Abyssinian"
			- "Scottish Fold"
			- "Birman" */,
  "sex": "", // enum of: "male" / "female"
  "ageInMonth": 1, // min: 1, max: 120082
  "description": "", // minLength 1, maxLength 200
  "imageUrls": [
    // minItems: 1, items: should be url
    "",
    "",
    ""
  ]
}
```

Response:

```json
{
  "message": "success",
  "data": {
    "id": "",
    "createdAt": ""
  }
}
```

- `201` successfully add cat
- `400` request doesn’t pass validation
- `401` request token is missing or expired

#### Get all cats

`GET /v1/cat`

| Parameter    | Type      | Description                                                                                                                                                                                      |
| :----------- | :-------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `id`         | `string`  | output based on the cat’s id                                                                                                                                                                     |
| `limit`      | `number`  | limit the output of data, default `limit=5`                                                                                                                                                      |
| `offset`     | `number`  | offset the output of data, default `offset=0`                                                                                                                                                    |
| `race`       | `enum`    | one of `Persian` or `Maine Coon` or `Siamese` or `Ragdoll` or `Bengal` or `Sphynx` or `British Sh or thair` or `Abyssinian` or `Scottish Fold` or `Birman`                                       |
| `sex`        | `enum`    | one of `male` or `female`                                                                                                                                                                        |
| `hasMatched` | `boolean` | cat has matched or not                                                                                                                                                                           |
| `ageInMonth` | `string`  | use it like this `ageInMonth=>4` searches data that have more than 4 months, `ageInMonth=<4` searches data that have less than 4 months, or `ageInMonth=4` searches data that have exact 4 month |
| `owned`      | `boolean` | cat that the user own                                                                                                                                                                            |
| `search`     | `string`  | contains the name of the cat                                                                                                                                                                     |

Response:

```json
{
  "message": "success",
  "data": [
    // ordered by newest first
    {
      "id": "",
      "name": "",
      "race": "",
      "sex": "",
      "ageInMonth": 1,
      "imageUrls": ["", "", ""],
      "description": "",
      "hasMatched": true,
      "createdAt": ""
    }
  ]
}
```

- `200` successfully get cats
- `401` request token is missing or expired

#### Update cat

`PUT /v1/cat/{id}`

Request Path Params

- `id` is the cat id that user want to edit

Request:

```json
{
  "name": "", // minLength 1, maxLength 30
  "race": "" /** enum of:
			- "Persian"
			- "Maine Coon"
			- "Siamese"
			- "Ragdoll"
			- "Bengal"
			- "Sphynx"
			- "British Shorthair"
			- "Abyssinian"
			- "Scottish Fold"
			- "Birman" */,
  "sex": "", // enum of: "male" / "female"
  "ageInMonth": 1, // min: 1, max: 120082
  "description": "", // minLength 1, maxLength 200
  "imageUrls": [
    // minItems: 1, items: should be url
    "",
    "",
    ""
  ]
}
```

Response:

- `200` successfully add cat
- `400` request doesn’t pass validation
- `401` request token is missing or expired
- `404` id is not found
- `400` sex is edited when cat is already requested to match

#### Delete cat

`DELETE /v1/cat/{id}`

Request Path Params

- `id` is the cat id that user want to delete

Response:

- `200` successfully delete cat
- `401` request token is missing or expired
- `404` id is not found

## Managing Cats

> [!WARNING]
> All request here should use Bearer Token from accessToken auth route

#### Create match request

`POST /v1/cat/match`

Request:

```json
{
  "matchCatId": "",
  "userCatId": "",
  "message": "" // minLength: 5, maxLength: 120
}
```

Response:

- `201` successfully send match request
- `404` if neither `matchCatId` / `userCatId` is not found
- `404` if `userCatId` is not belong to the user
- `400` if the cat’s gender is same
- `400` if both `matchCatId` &`userCatId` already matched
- `400` if `matchCatId` & `userCatId` is from the same owner
- `401` request token is missing or expired

#### Get match requests

`POST /v1/cat/match`

Response:

> [!NOTE]
> The information is shown in both the issuer and the receiver.

```json
{
  "message": "success",
  "data": [
    // ordered by newest first
    {
      "id": "",
      "issuedBy": {
        "name": "",
        "email": "",
        "createdAt": ""
      },
      "matchCatDetail": {
        "id": "",
        "name": "",
        "race": "",
        "sex": "",
        "description": "",
        "ageInMonth": 1,
        "imageUrls": ["", "", ""],
        "hasMatched": false,
        "createdAt": ""
      },
      "userCatDetail": {
        "id": "",
        "name": "",
        "race": "",
        "sex": "",
        "description": "",
        "ageInMonth": 1,
        "imageUrls": ["", "", ""],
        "hasMatched": false,
        "createdAt": ""
      },
      "message": "",
      "createdAt": ""
    }
  ]
}
```

- `200` successfully get match requests
- `401` request token is missing or expired

#### Approve match request

`POST /v1/cat/match/approve`

Request:

> [!NOTE]
> Once a match is approved, other match request that matches both the issuer and the receiver cat’s, get removed

```json
{
  "matchId": ""
}
```

Response:

- `200` successfully matches the cat match request
- `400` `matchId` is no longer valid
- `401` request token is missing or expired
- `404` `matchId` is not found

#### Reject match request

`POST /v1/cat/match/reject`

Request:

```json
{
  "matchId": ""
}
```

Response:

- `200` successfully reject the cat match request
- `400` `matchId` is no longer valid
- `401` request token is missing or expired
- `404` `matchId` is not found

#### Delete match request

`DELETE /v1/cat/match/{id}`

> [!WARNING]
> Match can only be deleted by issuer

Request Path Params

- `id` is the match id that user want to delete

Response:

- `200` successfully remove a match cat request
- `400` `matchId` is already approved / rejected
- `401` request token is missing or expired
- `404` `matchId` is not found
