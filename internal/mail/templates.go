package mail

const EmailVerificationTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome!</title>
</head>

<body>
	<div style="text-align: center; font-family: Arial, sans-serif;">
		<p>Welcome!</p>
		<p>Here is your verification code:</p>
		<p style="font-weight: bold; font-size: 20px; color: cornflowerblue">%s</p>
	</div>
</body>

</html>
`

const ResetPasswordTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password</title>
</head>

<body>
	<div style="text-align:center;">
		<p>Reset Password</p>
		<p>Please click this link to go to the password reset page:</p>
		<a href="%s">%s</a>  
	</div>
</body>

</html>
`
