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
    <title>Reset Password Code</title>
</head>

<body>
	<div style="text-align: center; font-family: Arial, sans-serif;">
		<p>Reset Password Code</p>
		<p>Here is your reset code:</p>
		<p style="font-weight: bold; font-size: 20px; color: cornflowerblue">%s</p>
	</div>
</body>

</html>
`

const OtpTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>One-Time Passcode</title>
</head>

<body>
	<div style="text-align: center; font-family: Arial, sans-serif;">
		<p>One-Time Passcode</p>
		<p>Here is your one-time passcode:</p>
		<p style="font-weight: bold; font-size: 20px; color: cornflowerblue">%s</p>
	</div>
</body>

</html>
`
