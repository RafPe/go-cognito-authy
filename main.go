package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/urfave/cli"
)

var (
	appName, appVer string
	ses             *session.Session
	sessionParams   session.Options
	cip             *cognitoidentityprovider.CognitoIdentityProvider
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "Used for quick testing auth on Cognito Auth Pool"
	app.Version = appVer
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Rafpe ( https://rafpe.ninja )",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "profile,p",
			Value:  "-",
			EnvVar: "AWS_PROFILE",
		},
		cli.StringFlag{
			Name:   "region",
			Value:  "",
			EnvVar: "AWS_DEFAULT_REGION",
		},
		cli.StringFlag{
			Name:   "access-key",
			EnvVar: "AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "secret-key",
			EnvVar: "AWS_SECRET_ACCESS_KEY",
		},
	}

	app.Before = func(c *cli.Context) error {
		var sesErr error

		if c.String("profile") != "-" {

			sessionParams = session.Options{
				Profile: c.String("profile"),
				Config:  aws.Config{Region: aws.String(c.String("region"))},
			}

			ses, sesErr = session.NewSessionWithOptions(sessionParams)
			if sesErr != nil {
				fmt.Println(sesErr)
				os.Exit(1)
			}
		}

		ses, sesErr = session.NewSession(&aws.Config{
			Region:      aws.String(c.String("region")),
			Credentials: credentials.NewStaticCredentials(c.String("access-key"), c.String("secret-key"), ""),
		})
		if sesErr != nil {
			fmt.Println(sesErr)
			os.Exit(1)
		}

		cip = cognitoidentityprovider.New(ses)

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "auth",
			Usage: "Authenticates user",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username",
					Usage: "username used in Cognito User Pool",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "Password for the username",
				},
				cli.StringFlag{
					Name:  "clientID",
					Value: "-",
					Usage: "App clientID from Cognito User Pool",
				},
			},
			Action: cmdAuthenticateUser,
		},
		{
			Name:  "admin",
			Usage: "Admin actions",
			Subcommands: []cli.Command{
				{
					Name:  "reset-pass",
					Usage: "Administratively resets password",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "username",
							Usage: "username used in Cognito User Pool",
						},
						cli.StringFlag{
							Name:  "pass-new",
							Usage: "New password for the username",
						},
						cli.StringFlag{
							Name:  "clientID",
							Value: "IP",
							Usage: "App clientID from Cognito User Pool",
						},
						cli.StringFlag{
							Name:  "userPoolID",
							Value: "IP",
							Usage: "Cognito User Pool id",
						},
						cli.StringFlag{
							Name:  "session",
							Value: "IP",
							Usage: "Session param from auth action",
						},
					},
					Action: cmdAdminResetPassword,
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Action = func(c *cli.Context) error {

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func cmdAdminResetPassword(c *cli.Context) error {
	username := c.String("username")
	passNew := c.String("pass-new")
	clientID := c.String("clientID")
	userPoolID := c.String("userPoolID")
	session := c.String("session")

	params := &cognitoidentityprovider.AdminRespondToAuthChallengeInput{
		ChallengeName: aws.String("NEW_PASSWORD_REQUIRED"),
		ChallengeResponses: map[string]*string{
			"NEW_PASSWORD": aws.String(passNew),
			"USERNAME":     aws.String(username),
		},
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
		Session:    aws.String(session),
	}

	adminChallengeResp, adminChallengeErr := cip.AdminRespondToAuthChallenge(params)
	if adminChallengeErr != nil {
		return adminChallengeErr
	}
	fmt.Println(adminChallengeResp)

	return nil
}

func cmdChangePassword(c *cli.Context) error {

	accessToken := c.String("token")
	passOld := c.String("pass-old")
	passNew := c.String("pass-new")

	params := &cognitoidentityprovider.ChangePasswordInput{
		AccessToken:      aws.String(accessToken),
		PreviousPassword: aws.String(passOld),
		ProposedPassword: aws.String(passNew),
	}

	newPassResponse, newPassErr := cip.ChangePassword(params)

	if newPassErr != nil {
		return newPassErr
	}

	fmt.Println(newPassResponse)

	return nil
}

// cmdAuthenticateUser invokes auth method with given params
// to get auth tokens.
//
func cmdAuthenticateUser(c *cli.Context) error {

	username := aws.String(c.String("username"))
	password := aws.String(c.String("password"))
	clientID := aws.String(c.String("clientID"))

	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": username,
			"PASSWORD": password,
		},
		ClientId: clientID,
	}

	authResponse, authError := cip.InitiateAuth(params)
	if authError != nil {
		return authError
	}

	fmt.Println(authResponse)

	return nil
}
