# go-cognito-authy
Congito User Pool quick auth CLI for testing purposes.

# Examples
First create user(s) in Cognito Pool ( or import them ) and then use one of commands below

## Authenticate
Try to authenticate using `myAwsProfile` profile and `eu-central-1` as region
```
go-cognito-authy -profile myAwsProfile -region eu-central-1 auth --username <username-in-pool> --password '<some-magical-password>'  --clientID <app-client-id>
```

Received response  with tokens 
```
> [SHELL] RafPe $ go-cognito-authy --profile cloudy --region eu-central-1 auth --username rafpe --password 'Password.0ne!'  --clientID 2jxxxiuui123
{
AuthenticationResult: {
    AccessToken: "eyJraWQiOiJ0QXVBNmxtNngrYkxoSmZ",
    ExpiresIn: 3600,
    IdToken: "eyJraWQiOiJ0bHF2UElTV0pn",
    RefreshToken: "eyJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R-TpkR_uompG7fyajYeFvn-rJVC_tDO4pB3",
    TokenType: "Bearer"
},
ChallengeParameters: {}
}
```

Received response with challenge ( if it was NEW_PASSWORD_REQUIRED then I would admin reset my pass ;) ) 
```
> [INSERT] RafPe $ go-cognito-authy --profile cloudy --region eu-central-1 auth --username rafpe --password 'Password.0ne!'  --clientID 2jxxxiuui123
{
ChallengeName: "NEW_PASSWORD_REQUIRED",
ChallengeParameters: {
    requiredAttributes: "[]",
    userAttributes: "{\"email_verified\":\"true\",\"email\":\"mee@rafpe.ninja\"}",
    USER_ID_FOR_SRP: "rafpe"
},
Session: "bCqSkLeoJR_ys...."
}
```

## Admin 
Administrative actions
* Reset password
```
> [INSERT] RafPe $ go-cognito-authy --profile cloudy -region eu-central-1 admin reset-pass --username rafpe --pass-new 'Password.0ne2!' --clientID 2jxxxiuui123 --userPoolID  eu-central-1_CWNnTiR0j --session "bCqSkLeoJR_ys...."
```
