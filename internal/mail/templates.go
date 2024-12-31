package mail

var EmailVerificationTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome!</title>
</head>

<body>
	<div style="text-align:center;">
		<p>Welcome!</p>
		<p>Please click this link to verify your email address:</p>
		<a href="%s">%s</a>  
	</div>
</body>

</html>
`
